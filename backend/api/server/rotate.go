package server

import (
	"context"
	"regexp"
	"time"

	coverageV1Alpha1 "github.com/konflux-ci/quality-dashboard/api/apis/codecov/v1alpha1"
	repoV1Alpha1 "github.com/konflux-ci/quality-dashboard/api/apis/github/v1alpha1"
	ociloader "github.com/konflux-ci/quality-dashboard/api/server/oci-loader"
	"go.uber.org/zap"
)

// nolint:all
var RegexpCompiler = regexp.MustCompile("(-main-|-master-)(.*?)(\\/)")

const (
	// directory to download the oci artifacts blobs
	OCICache = "/tmp/oci-artifacts"

	// The directory where all the oci artifacts will be extracted
	ArtifactsDir = "/tmp/oci-cache"
)

// rotationStrategy describes a strategy for generating server configuration from a file.
type rotationStrategy struct {
	// Time between rotations.
	rotationFrequency time.Duration
}

// startUpdateStorage begins repo information rotation in a new goroutine, closing once the context is canceled.
func (s *Server) startUpdateStorage(ctx context.Context, strategy rotationStrategy) {
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
	_, err := s.cfg.Storage.GetAllTeamsFromDB()
	if err != nil {
		s.cfg.Logger.Sugar().Errorf("Failed to update cache", zap.Error(err))
		return err
	}

	s.UpdateProwStatusByTeam()

	loader := ociloader.NewLoader(ociloader.LoaderConfig{
		NumWorkers:          5,
		ProcessingThreshold: time.Minute * 45,
		ArtifactsDir:        "/tmp/oci-artifacts",
		CacheDir:            "/tmp/oci-cache",
		CleanupOnCompletion: true,
	}, s.cfg.Storage)

	if loader.Run(context.Background()) != nil {
		s.cfg.Logger.Sugar().Errorf("failed to update the oci artifact", zap.Error(err))
	}

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
		rotationFrequency: time.Minute * 15,
	}
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

		// update coverage info
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

	return nil
}

func (s *Server) getCodeCoverage(gitOrganization string, repoName string) (float64, string, error) {
	return s.cfg.CodeCov.GetCodeCovInfo(gitOrganization, repoName)
}
