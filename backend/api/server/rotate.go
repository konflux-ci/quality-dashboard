package server

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/andygrunwald/go-jira"
	coverageV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/codecov/v1alpha1"
	configurationV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/configuration/v1alpha1"
	repoV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/github/v1alpha1"
	"github.com/redhat-appstudio/quality-studio/pkg/storage"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
	"go.uber.org/zap"
)

// nolint:all
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

	s.UpdateProwStatusByTeam()

	s.cfg.Logger.Info("starting to rotate jira bugs")
	for _, team := range teamArr {
		if team.JiraKeys == "" {
			continue
		}

		// temporary until previous teams has config defined
		_, err := s.cfg.Storage.GetConfiguration(team.TeamName)
		if err != nil && err == storage.ErrNotFound {
			query := fmt.Sprintf("project in (%s) AND type = Bug", team.JiraKeys)
			jiraCfg := configurationV1Alpha1.JiraConfig{
				BugsCollectQuery: query,
				CiImpactQuery:    query + " AND labels=ci-fail",
			}
			b, err := json.Marshal(jiraCfg)
			if err != nil {
				return err
			}

			config := configurationV1Alpha1.Configuration{
				TeamName:      team.TeamName,
				JiraConfig:    string(b),
				BugSLOsConfig: "",
			}

			err = s.cfg.Storage.CreateConfiguration(config)
			if err != nil {
				return err
			}
		}

		s.cfg.Logger.Info(fmt.Sprintf("starting to rotate jira bugs for team %s", team.TeamName))
		if err := s.rotateJiraBugs(team.JiraKeys, team); err != nil {
			return err
		}
	}
	s.cfg.Logger.Info("finishing to rotate jira bugs")

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
		rotationFrequency: time.Minute * 90,
	}
}

func shouldBeDeleted(jiraKey string, bugs []jira.Issue) bool {
	for _, bug := range bugs {
		if bug.Key == jiraKey {
			return false
		}
	}
	return true
}

func remove(l []string, item string) []string {
	for i, other := range l {
		if other == item {
			return append(l[:i], l[i+1:]...)
		}
	}

	return l
}

func (s *Server) rotateJiraBugs(jiraKeys string, team *db.Teams) error {
	cfg, err := s.cfg.Storage.GetConfiguration(team.TeamName)
	if err != nil {
		return err
	}

	jiraCfg := configurationV1Alpha1.JiraConfig{}
	err = json.Unmarshal([]byte(cfg.JiraConfig), &jiraCfg)
	if err != nil {
		return err
	}

	bugs := s.cfg.Jira.GetBugsByJQLQuery(jiraCfg.BugsCollectQuery)
	if err := s.cfg.Storage.CreateJiraBug(bugs, team); err != nil {
		return err
	}

	projects := strings.Split(team.JiraKeys, ",")
	bugsInDb := make([]*db.Bugs, 0)
	for _, project := range projects {
		bgs, err := s.cfg.Storage.GetAllJiraBugsByProject(project)
		if err != nil {
			return err
		}
		bugsInDb = append(bugsInDb, bgs...)
	}

	// clean bugs that changed project or jira type
	for _, bugInDb := range bugsInDb {
		deleted := shouldBeDeleted(bugInDb.JiraKey, bugs)
		if deleted {
			if err := s.cfg.Storage.DeleteJiraBugByJiraKey(bugInDb.JiraKey); err != nil {
				return err
			}

			// also clean in CI Impact bugs
			jiraCfg.CiImpactBugs = remove(jiraCfg.CiImpactBugs, bugInDb.JiraKey)
		}
	}

	// update config
	jiraCfgUpdated, err := json.Marshal(jiraCfg)
	if err != nil {
		return err
	}

	config := configurationV1Alpha1.Configuration{
		TeamName:      cfg.TeamName,
		JiraConfig:    string(jiraCfgUpdated),
		BugSLOsConfig: cfg.JiraConfig,
	}

	return s.cfg.Storage.CreateConfiguration(config)
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
