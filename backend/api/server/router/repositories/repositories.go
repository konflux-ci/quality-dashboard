package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	coverageV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/codecov/v1alpha1"
	repoV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/github/v1alpha1"
	"github.com/redhat-appstudio/quality-studio/api/types"
	"github.com/redhat-appstudio/quality-studio/pkg/storage"
	"github.com/redhat-appstudio/quality-studio/pkg/utils/httputils"
	"go.uber.org/zap"
)

var repository GitRepositoryRequest

// version godoc
// @Summary Github repositories info
// @Description returns all repository information stored in database
// @Tags Github Repositories API
// @Accept json
// @Produce json
// @Router /repositories/list [get]
// @Success 200 {array} storage.RepositoryQualityInfo
// @Failure 400 {object} types.ErrorResponse
func (rp *repositoryRouter) listAllRepositoriesQuality(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	teamName := r.URL.Query()["team_name"]
	startDate := r.URL.Query()["start_date"]
	endDate := r.URL.Query()["end_date"]

	if len(teamName) == 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "team_name value not present in query",
			StatusCode: 400,
		})
	} else if len(startDate) == 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "start_date value not present in query",
			StatusCode: 400,
		})
	} else if len(endDate) == 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "end_date value not present in query",
			StatusCode: 400,
		})
	}

	team, err := rp.Storage.GetTeamByName(teamName[0])
	if err != nil {
		rp.Logger.Error("Failed to fetch team. Make sure the team exists", zap.String("team", repository.Team), zap.Error(err))

		return httputils.WriteJSON(w, http.StatusBadRequest, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
	}
	repos, err := rp.Storage.ListRepositoriesQualityInfo(team, startDate[0], endDate[0])
	if err != nil {
		return httputils.WriteJSON(w, http.StatusBadRequest, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
	}
	// in the OpenShift CI tab, do not include repos that do not exist in OpenShift CI
	isOpenShiftCI := r.URL.Query()["openshift_ci"]
	if len(isOpenShiftCI) != 0 && isOpenShiftCI[0] == "true" {
		clean := make([]storage.RepositoryQualityInfo, 0)
		for _, repo := range repos {
			if rp.Github.CheckIfRepoExistsInOpenshiftCI(repo.GitOrganization, repo.RepositoryName) &&
				len(rp.Github.GetJobTypes(repo.GitOrganization, repo.RepositoryName)) > 0 {

				clean = append(clean, repo)
			}
		}
		repos = clean
	}

	return httputils.WriteJSON(w, http.StatusOK, repos)
}

// version godoc
// @Summary Github repositories info
// @Description returns the Server information as a JSON
// @Tags Github Repositories API
// @Accept json
// @Produce json
// @Param repository body GitRepositoryRequest true "repository name"
// @Router /repositories/create [post]
// @Success 200 {object} types.SuccessResponse
// @Failure 400 {object} types.ErrorResponse
func (rp *repositoryRouter) createRepositoryHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := json.NewDecoder(r.Body).Decode(&repository); err != nil {
		return httputils.WriteJSON(w, http.StatusBadRequest, &types.ErrorResponse{
			Message:    "Error reading repository/git_organization/team_name value from body",
			StatusCode: http.StatusBadRequest,
		})
	}

	githubRepo, err := rp.Github.GetGithubRepositoryInformation(repository.GitOrganization, repository.GitRepository)
	if err != nil {
		rp.Logger.Error("Failed to fetch repository info from github", zap.String("repository", repository.GitRepository), zap.String("git_organization", repository.GitOrganization), zap.Error(err))

		return httputils.WriteJSON(w, http.StatusBadRequest, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
	}

	team, err := rp.Storage.GetTeamByName(repository.Team)
	if err != nil {
		rp.Logger.Error("Failed to fetch team. Make sure the team exists", zap.String("team", repository.Team), zap.Error(err))

		return httputils.WriteJSON(w, http.StatusBadRequest, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
	}

	description := githubRepo.GetDescription()
	if description == "" {
		description = "Repository don't contain a description"
	}
	createdRepo, err := rp.Storage.CreateRepository(repoV1Alpha1.Repository{
		ID:   fmt.Sprint(githubRepo.GetID()),
		Name: githubRepo.GetName(),
		Owner: repoV1Alpha1.Owner{
			Login: githubRepo.Owner.GetLogin(),
		},
		Description: description,
		URL:         githubRepo.GetHTMLURL(),
	}, team.ID)
	if err != nil {
		return httputils.WriteJSON(w, http.StatusBadRequest, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
	}

	coverage, covTrend, err := rp.CodeCov.GetCodeCovInfo(githubRepo.Owner.GetLogin(), githubRepo.GetName())
	if err != nil {
		rp.Logger.Error("Failed to fetch repository info from codecov", zap.String("repository", repository.GitRepository), zap.String("git_organization", repository.GitOrganization), zap.Error(err))

		return httputils.WriteJSON(w, http.StatusBadRequest, &types.ErrorResponse{
			Message:    "Failed to obtain repositories. There are no repository cached",
			StatusCode: http.StatusBadRequest,
		})
	}

	err = rp.Storage.CreateCoverage(coverageV1Alpha1.Coverage{
		GitOrganization:    githubRepo.Owner.GetLogin(),
		RepositoryName:     githubRepo.GetName(),
		CoveragePercentage: coverage,
		CoverageTrend:      covTrend,
	}, createdRepo.ID)

	if err != nil {
		return httputils.WriteJSON(w, http.StatusBadRequest, &types.ErrorResponse{
			Message:    "Failed to save coverage data in database. There are no repository cached",
			StatusCode: http.StatusBadRequest,
		})
	}

	workflows, err := rp.Github.GetRepositoryWorkflows(githubRepo.Owner.GetLogin(), githubRepo.GetName())
	for _, w := range workflows.Workflows {
		if err := rp.Storage.CreateWorkflows(repoV1Alpha1.Workflow{
			Name:     w.GetName(),
			BadgeURL: w.GetBadgeURL(),
			HTMLURL:  w.GetHTMLURL(),
			JobURL:   w.GetURL(),
			State:    w.GetState(),
		}, createdRepo.ID); err != nil {
			rp.Logger.Error("failed to save repository in database", zap.String("repository", repository.GitRepository), zap.String("git_organization", repository.GitOrganization), zap.Error(err))
		}
	}

	if err != nil {
		rp.Logger.Error("Failed to fetch github actions info from github", zap.String("repository", repository.GitRepository), zap.String("git_organization", repository.GitOrganization), zap.Error(err))

		return httputils.WriteJSON(w, http.StatusBadRequest, &types.ErrorResponse{
			Message:    "Failed to save workflows data in database. There are no repository cached",
			StatusCode: http.StatusBadRequest,
		})
	}

	prs, err := rp.Github.GetPullRequestsInRange(context.TODO(), repoV1Alpha1.ListPullRequestsOptions{
		Repository: githubRepo.GetName(),
		Owner:      githubRepo.Owner.GetLogin(),
		TimeField:  repoV1Alpha1.PullRequestNone,
	}, time.Now().AddDate(0, -3, 0), time.Now())

	if err != nil {
		rp.Logger.Error("Failed to fetch pull requests from github", zap.String("repository", repository.GitRepository), zap.String("git_organization", repository.GitOrganization), zap.Error(err))

		return httputils.WriteJSON(w, http.StatusBadRequest, &types.ErrorResponse{
			Message:    "Failed to get pull requests",
			StatusCode: http.StatusBadRequest,
		})
	}

	if err := rp.Storage.CreatePullRequests(prs, createdRepo.ID); err != nil {
		rp.Logger.Error("Failed to fetch pull requests info from github", zap.String("repository", repository.GitRepository), zap.String("git_organization", repository.GitOrganization), zap.Error(err))

		return httputils.WriteJSON(w, http.StatusBadRequest, &types.ErrorResponse{
			Message:    "Failed to save pull requests data in database. There are no repository cached",
			StatusCode: http.StatusBadRequest,
		})
	}

	return httputils.WriteJSON(w, http.StatusOK, types.SuccessResponse{
		Message:    "Successfully created repository in quality-studio",
		StatusCode: http.StatusCreated,
	})
}

// Version godoc
// @Summary Github repositories info
// @Description delete a given repository from a organization
// @Tags Github Repositories API
// @Produce json
// @Param repository body GitRepositoryRequest true "repository name"
// @Router /repositories/delete [delete]
// @Success 200 {object} types.SuccessResponse
// @Failure 400 {object} types.ErrorResponse
func (rp *repositoryRouter) deleteRepositoryHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := json.NewDecoder(r.Body).Decode(&repository); err != nil {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "incorrect data received to server",
			StatusCode: 400,
		})
	}

	if repository.GitRepository == "" {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "Failed to remove repository. Field 'repository_name' missing",
			StatusCode: 400,
		})
	}
	if repository.GitOrganization == "" {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "Failed to remove repository. Field 'git_organization' missing",
			StatusCode: 400,
		})
	}
	err := rp.Storage.DeleteRepository(repository.GitRepository, repository.GitOrganization)
	if err != nil {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "Failed to remove repository",
			StatusCode: 400,
		})
	}
	return httputils.WriteJSON(w, http.StatusOK, types.SuccessResponse{
		Message: "Repository deleted",
	})
}

func (rp *repositoryRouter) getJobTypesFromRepo(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	repositoryName := r.URL.Query()["repository_name"]
	gitOrganization := r.URL.Query()["git_organization"]

	if len(repositoryName) == 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "repository_name value not present in query",
			StatusCode: 400,
		})
	} else if len(gitOrganization) == 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "git_organization value not present in query",
			StatusCode: 400,
		})
	}

	jobTypes := rp.Github.GetJobTypes(gitOrganization[0], repositoryName[0])

	return httputils.WriteJSON(w, http.StatusOK, jobTypes)
}

func (rp *repositoryRouter) checkGithubRepositoryUrl(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	repoName := r.URL.Query()["repository_name"]
	repoOrg := r.URL.Query()["git_organization"]

	if len(repoName) == 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "repository_name value not present in query",
			StatusCode: 400,
		})
	} else if len(repoOrg) == 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "git_organization value not present in query",
			StatusCode: 400,
		})
	}

	_, err := rp.Github.GetGithubRepositoryInformation(repoOrg[0], repoName[0])
	if err != nil {
		rp.Logger.Error("Failed to fetch repository info from github", zap.String("repository", repoName[0]), zap.String("git_organization", repoOrg[0]), zap.Error(err))

		return httputils.WriteJSON(w, http.StatusBadRequest, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
	}

	return httputils.WriteJSON(w, http.StatusOK, types.SuccessResponse{
		Message:    "Successfully verified that repository exists",
		StatusCode: http.StatusCreated,
	})
}

func (rp *repositoryRouter) checkGithubRepositoryExists(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	repoName := r.URL.Query()["repository_name"]
	repoOrg := r.URL.Query()["git_organization"]

	if len(repoName) == 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "repository_name value not present in query",
			StatusCode: 400,
		})
	} else if len(repoOrg) == 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "git_organization value not present in query",
			StatusCode: 400,
		})
	}

	teams, err := rp.Storage.GetAllTeamsFromDB()
	if err != nil {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "Failed to get repository",
			StatusCode: 400,
		})
	}

	for _, team := range teams {
		repos, err := rp.Storage.ListRepositories(team)
		if err != nil {
			return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
				Message:    "Failed to get repository",
				StatusCode: 400,
			})
		}

		for _, repo := range repos {
			if repo.Name == repoName[0] && repo.Owner.Login == repoOrg[0] {
				return httputils.WriteJSON(w, http.StatusOK, team.TeamName)
			}
		}
	}

	return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
		Message:    "Failed to get repository",
		StatusCode: 400,
	})

}
