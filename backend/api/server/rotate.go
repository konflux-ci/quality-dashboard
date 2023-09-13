package server

import (
	"context"
	"fmt"
	"regexp"
	"time"

	coverageV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/codecov/v1alpha1"
	repoV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/github/v1alpha1"
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
	s.UpdateFailuresByTeam()
	err = s.UpdateDataBaseRepoByTeam()
	if err != nil {
		s.cfg.Logger.Sugar().Errorf("Failed to update cache", zap.Error(err))
		return err
	}

	return nil
}
func staticRotationStrategy() rotationStrategy {
	return rotationStrategy{
		rotationFrequency: time.Minute * 2,
	}
}

/**
bugs := s.Jira.GetBugsByJQLQuery(fmt.Sprintf("project in (%s) AND type = Bug", teams.JiraKeys))
if err := s.Storage.CreateJiraBug(bugs, teams); err != nil {
	return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
		Message:    err.Error(),
		StatusCode: http.StatusInternalServerError,
	})
}

*/

func (s *Server) rotateJiraBugs(jiraKeys string, team *db.Teams) error {
	bugs := s.cfg.Jira.GetBugsByJQLQuery(fmt.Sprintf("project in (%s) AND type = Bug", team.JiraKeys))
	if err := s.cfg.Storage.CreateJiraBug(bugs, team); err != nil {
		return err
	}

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

		//update coverage info
		coverage, covTrend, err := s.getCodeCoverage(repo.Owner.Login, repo.Name)
		if err != nil {
			return err
		}

		err = s.cfg.Storage.UpdateCoverage(coverageV1Alpha1.Coverage{
			GitOrganization:    repo.Owner.Login,
			RepositoryName:     repo.Name,
			CoveragePercentage: coverage,
			CoverageTrend:      covTrend,
		}, repo.Name)
		if err != nil {
			return err
		}
	}
	s.cfg.Logger.Info("Successfully updated the storage data")

	return nil
}

func (s *Server) getCodeCoverage(gitOrganization string, repoName string) (float64, string, error) {
	return s.cfg.CodeCov.GetCodeCovInfo(gitOrganization, repoName)
}
