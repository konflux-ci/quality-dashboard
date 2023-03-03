package server

import (
	"context"
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/andygrunwald/go-jira"
	coverageV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/codecov/v1alpha1"
	repoV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/github/v1alpha1"
	"github.com/redhat-appstudio/quality-studio/api/apis/jira/v1alpha1"
	"github.com/redhat-appstudio/quality-studio/pkg/connectors/codecov"
	"go.uber.org/zap"
)

var RegexpCompiler = regexp.MustCompile("(-main-|-master-)(.*?)(\\/)")

// rotationStrategy describes a strategy for generating server configuration from a file.
type rotationStrategy struct {
	// Time between rotations.
	rotationFrequency time.Duration
}

// startUpdateCache begins repo information rotation in a new goroutine, closing once the context is canceled.
func (s *Server) startUpdateStorage(ctx context.Context, strategy rotationStrategy, now func() time.Time) {
	go func() {
		// Try to rotate immediately so properly configured repositories.
		if err := s.rotate(); err != nil {
			s.cfg.Logger.Sugar().Infof("Update failed: %v", err)
		}
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(strategy.rotationFrequency):
				if err := s.rotate(); err != nil {
					s.cfg.Logger.Sugar().Infof("failed to update cache: %v", err)
				}
			}
		}
	}()
}

func (s *Server) rotate() error {
	s.rotateJiraBugs()

	s.UpdateProwStatusByTeam()
	err := s.UpdateDataBaseRepoByTeam()
	if err != nil {
		s.cfg.Logger.Sugar().Errorf("Failed to update cache", zap.Error(err))
		return err
	}

	return nil
}
func staticRotationStrategy() rotationStrategy {
	return rotationStrategy{
		rotationFrequency: time.Second * 30,
	}
}

func (s *Server) rotateJiraBugs() {
	bugs := s.cfg.Jira.GetBugsByJQLQuery(`project in ("Hybrid Application Service", "Stonesoup", "CodeReady Toolchain", "GitOps Service", "Pipeline Service", "SVPI", "Stonesoup Build") AND type = Bug`)
	wg := new(sync.WaitGroup)
	for keyString, bugValue := range bugs {
		wg.Add(1)
		go func(keyString int, bug jira.Issue, wg *sync.WaitGroup) {
			defer wg.Done()
			if bug.Fields.Priority.Name == "" {
				bug.Fields.Priority.Name = "No Data"
			}
			if err := s.cfg.Storage.CreateJiraBug(v1alpha1.JiraBug{
				JiraKey:   bug.Key,
				CreatedAt: time.Time(bug.Fields.Created),
				UpdatedAt: time.Time(bug.Fields.Updated),
				Priority:  bug.Fields.Priority.Name,
				Status:    bug.Fields.Status.Name,
				Summary:   bug.Fields.Summary,
				Url:       fmt.Sprintf("https://issues.redhat.com/browse/%s", bug.Key),
			}); err != nil {
				s.cfg.Logger.Sugar().Errorf("failed to update jiras %s, %v", bug.Key, err)
			}
		}(keyString, bugValue, wg)
	}

	wg.Wait()

	s.cfg.Logger.Sugar().Info("successfully update jira bugs in database")
}

func (s *Server) UpdateDataBaseRepoByTeam() error {
	teamArr, err := s.cfg.Storage.GetAllTeamsFromDB()
	if err != nil {
		return err
	}

	for _, team := range teamArr {
		repo, err := s.cfg.Storage.ListRepositories(team)
		if err != nil {
			return err
		}

		if err := s.CacheRepositoriesInformation(repo); err != nil {
			s.cfg.Logger.Sugar().Errorf("failed to update repository %s, %v", repo, err)
		}
	}
	return nil
}

func (s *Server) CacheRepositoriesInformation(storageRepos []repoV1Alpha1.Repository) error {
	for _, repo := range storageRepos {

		// update prs info
		prs, err := s.cfg.Github.GetRepositoryPullRequests(repo.Organization, repo.Name)
		if err != nil {
			return err
		}

		for _, pr := range prs {
			s.cfg.Storage.CreatePullRequest(repoV1Alpha1.PullRequest{
				RepositoryName:         repo.Name,
				RepositoryOrganization: repo.Organization,
				Number:                 pr.GetNumber(),
				Title:                  pr.GetTitle(),
				CreatedAt:              pr.GetCreatedAt(),
				MergedAt:               pr.GetMergedAt(),
				ClosedAt:               pr.GetClosedAt(),
				State:                  pr.GetState(),
				Author:                 *pr.GetUser().Login,
			}, repo.ID)
		}

		//update coverage info
		coverage, err := s.getCodeCoverage(repo.Organization, repo.Name)

		if err != nil {
			return err
		}
		totalCoverageConverted, _ := coverage.Totals.Coverage.Float64()
		err = s.cfg.Storage.UpdateCoverage(coverageV1Alpha1.Coverage{
			GitOrganization:    repo.Organization,
			RepositoryName:     repo.Name,
			CoveragePercentage: totalCoverageConverted,
		}, repo.Name)
		if err != nil {
			return err
		}
	}
	s.cfg.Logger.Info("Successfully updated the storage data")

	return nil
}

func (s *Server) getCodeCoverage(gitOrganization string, repoName string) (codecov.CoverageSpec, error) {
	return s.cfg.CodeCov.GetCodeCovInfo(gitOrganization, repoName)
}
