package server

import (
	"context"
	"regexp"
	"time"

	"github.com/redhat-appstudio/quality-studio/api/apis/codecov"
	"github.com/redhat-appstudio/quality-studio/pkg/storage"
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
		rotationFrequency: time.Minute * 30,
	}
}

func (s *Server) UpdateDataBaseRepoByTeam() error {
	teamArr, _ := s.cfg.Storage.GetAllTeamsFromDB()

	for _, team := range teamArr {
		repo, _ := s.cfg.Storage.ListRepositories(team)
		if err := s.CacheRepositoriesInformation(repo); err != nil {
			s.cfg.Logger.Sugar().Errorf("failed to update repository %s", repo)
		}
	}
	return nil
}

func (s *Server) CacheRepositoriesInformation(storageRepos []storage.Repository) error {
	for _, repo := range storageRepos {
		coverage, err := s.getCodeCoverage(repo.GitOrganization, repo.RepositoryName)

		if err != nil {
			return err
		}
		totalCoverageConverted, _ := coverage.Totals.Coverage.Float64()
		err = s.cfg.Storage.UpdateCoverage(storage.Coverage{
			GitOrganization:    repo.GitOrganization,
			RepositoryName:     repo.RepositoryName,
			CoveragePercentage: totalCoverageConverted,
		}, repo.RepositoryName)
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
