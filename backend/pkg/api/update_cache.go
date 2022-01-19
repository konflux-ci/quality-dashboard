package api

import (
	"context"
	"time"

	"github.com/flacatus/qe-dashboard-backend/pkg/api/apis/codecov"
	"github.com/flacatus/qe-dashboard-backend/pkg/api/apis/github"
	"github.com/flacatus/qe-dashboard-backend/pkg/storage"
	"go.uber.org/zap"
)

// rotationStrategy describes a strategy for generating server configuration from a file.
type rotationStrategy struct {
	// Time between rotations.
	rotationFrequency time.Duration
}

// startUpdateCache begins repo information rotation in a new goroutine, closing once the context is canceled.
func (s *Server) startUpdateStorage(ctx context.Context, strategy rotationStrategy, now func() time.Time) {
	// Try to rotate immediately so properly configured repositories.
	if err := s.rotate(); err != nil {
		s.logger.Sugar().Infof("Update failed: %v", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(strategy.rotationFrequency):
				if err := s.rotate(); err != nil {
					s.logger.Sugar().Infof("failed to update cache: %v", err)
				}
			}
		}
	}()
}

func (s *Server) rotate() error {
	err := s.CacheRepositoriesInformation()
	if err != nil {
		s.logger.Sugar().Errorf("Failed to update cache", zap.Error(err))
		return err
	}

	return nil
}

func staticRotationStrategy() rotationStrategy {
	return rotationStrategy{
		// Setting these values to 4 hours is easier than having a flag indicating no rotation.
		rotationFrequency: time.Minute * 1,
	}
}

func (s *Server) CacheRepositoriesInformation() error {
	storageRepos, err := s.config.Storage.ListRepositories()
	if err != nil {
		return err
	}

	for _, repo := range storageRepos {
		workf, err := s.getGithubWorkflows(repo.GitOrganization, repo.RepositoryName)
		if err != nil {
			return err
		}

		for _, w := range workf.Workflows {
			err = s.config.Storage.ReCreateWorkflow(storage.GithubWorkflows{
				WorkflowName: w.Name,
				BadgeURL:     w.BadgeURL,
				HTMLURL:      w.HTML_URL,
				JobURL:       w.JobURL,
				State:        w.State,
			}, repo.RepositoryName)
			if err != nil {
				return err
			}
		}

		coverage, err := s.getCodeCoverage(repo.GitOrganization, repo.RepositoryName)
		if err != nil {
			return err
		}
		totalCoverageConverted, _ := coverage.Commit.Totals.TotalCoverage.Float64()
		err = s.config.Storage.UpdateCoverage(storage.Coverage{
			GitOrganization:    repo.GitOrganization,
			RepositoryName:     repo.RepositoryName,
			CoveragePercentage: totalCoverageConverted,
		}, repo.RepositoryName)
		if err != nil {
			return err
		}
	}
	s.logger.Info("Successfully updated the storage data")

	return nil
}

func (s *Server) getGithubWorkflows(gitOrganization string, repoName string) (github.GitHubActionsResponse, error) {
	return s.githubAPI.GetRepositoryWorkflows(gitOrganization, repoName)
}

func (s *Server) getCodeCoverage(gitOrganization string, repoName string) (codecov.GitHubTagResponse, error) {
	return s.codecovAPI.GetCodeCovInfo(gitOrganization, repoName)
}
