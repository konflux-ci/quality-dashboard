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
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
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
	teamArr, err := s.cfg.Storage.GetAllTeamsFromDB()
	if err != nil {
		s.cfg.Logger.Sugar().Errorf("Failed to update cache", zap.Error(err))
		return err
	}
	for _, team := range teamArr {
		if team.JiraKeys == "" {
			continue
		}
		if err := s.rotateJiraBugs(team.JiraKeys, team); err != nil {
			return err
		}
	}

	s.UpdateProwStatusByTeam()
	err = s.UpdateDataBaseRepoByTeam()
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

func (s *Server) rotateJiraBugs(jiraKeys string, team *db.Teams) error {
	bugs := s.cfg.Jira.GetBugsByJQLQuery(fmt.Sprintf("project in (%s) AND type = Bug", team.JiraKeys))
	wg := new(sync.WaitGroup)
	for keyString, bugValue := range bugs {
		wg.Add(1)
		go func(keyString int, bug jira.Issue, wg *sync.WaitGroup) {
			defer wg.Done()
			var bugIsResolved bool
			var diff float64
			if bug.Fields.Status.Name == "Closed" || bug.Fields.Status.Name == "Resolved" || bug.Fields.Status.Name == "Done" {
				resolvedTime := time.Time(bug.Fields.Resolutiondate).UTC()
				createdTime := time.Time(bug.Fields.Created).UTC()

				diff = resolvedTime.Sub(createdTime).Hours()
				bugIsResolved = true
			}

			if err := s.cfg.Storage.CreateJiraBug(v1alpha1.JiraBug{
				JiraKey:        bug.Key,
				CreatedAt:      time.Time(bug.Fields.Created),
				UpdatedAt:      time.Time(bug.Fields.Updated),
				ResolvedAt:     time.Time(bug.Fields.Resolutiondate),
				IsResolved:     bugIsResolved,
				ResolutionTime: diff,
				Priority:       bug.Fields.Priority.Name,
				Status:         bug.Fields.Status.Name,
				Summary:        bug.Fields.Summary,
				Url:            fmt.Sprintf("https://issues.redhat.com/browse/%s", bug.Key),
			}, team); err != nil {
				s.cfg.Logger.Sugar().Errorf("failed to update jiras %s, %v", bug.Key, err)
			}
		}(keyString, bugValue, wg)
	}

	wg.Wait()

	return nil
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
	currentTime := time.Now()
	for _, repo := range storageRepos {
		prs, err := s.cfg.Github.GetPullRequestsInRange(context.TODO(), repoV1Alpha1.ListPullRequestsOptions{
			Repository: repo.Name,
			Owner:      repo.Owner.Login,
			TimeField:  repoV1Alpha1.PullRequestNone,
		}, currentTime.AddDate(0, -3, 0), currentTime)
		if err != nil {
			return err
		}

		if err := s.cfg.Storage.CreatePullRequests(prs, repo.ID); err != nil {
			return err
		}

		totalRetestRepoAvg, err := s.cfg.Github.RetestsToMerge(fmt.Sprintf("%s/%s", repo.Owner.Login, repo.Name))
		if err != nil {
			return err
		}
		//update coverage info
		coverage, err := s.getCodeCoverage(repo.Owner.Login, repo.Name)

		if err != nil {
			return err
		}
		totalCoverageConverted, _ := coverage.Totals.Coverage.Float64()
		err = s.cfg.Storage.UpdateCoverage(coverageV1Alpha1.Coverage{
			GitOrganization:            repo.Owner.Login,
			RepositoryName:             repo.Name,
			CoveragePercentage:         totalCoverageConverted,
			AverageToRetestPullRequest: totalRetestRepoAvg,
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
