package repositories

import (
	"context"
	"encoding/json"
	"net/http"

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

	if len(teamName) == 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "team_name value not present in query",
			StatusCode: 400,
		})
	}

	team, err := rp.Storage.GetTeamByName(teamName[0])
	if err != nil {
		rp.Logger.Error("Failed to fetch team. Make sure the team exists", zap.String("team", repository.Team), zap.Error(err))

		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
	}
	repos, err := rp.Storage.ListRepositoriesQualityInfo(team)

	if err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		})
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
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    "Error reading repository/git_organization/team_name value from body",
			StatusCode: http.StatusBadRequest,
		})
	}

	githubRepo, err := rp.Github.GetGithubRepositoryInformation(repository.GitOrganization, repository.GitRepository)
	if err != nil {
		rp.Logger.Error("Failed to fetch repository info from github", zap.String("repository", repository.GitRepository), zap.String("git_organization", repository.GitOrganization), zap.Error(err))

		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
	}

	team, err := rp.Storage.GetTeamByName(repository.Team)
	if err != nil {
		rp.Logger.Error("Failed to fetch team. Make sure the team exists", zap.String("team", repository.Team), zap.Error(err))

		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
	}

	description := githubRepo.GetDescription()
	if description == "" {
		description = "Repository don't contain a description"
	}
	createdRepo, err := rp.Storage.CreateRepository(storage.Repository{
		RepositoryName:  githubRepo.GetName(),
		GitOrganization: githubRepo.Owner.GetLogin(),
		Description:     description,
		GitURL:          githubRepo.GetHTMLURL(),
	}, team.ID)
	if err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
	}

	coverage, err := rp.CodeCov.GetCodeCovInfo(githubRepo.Owner.GetLogin(), githubRepo.GetName())
	if err != nil {
		rp.Logger.Error("Failed to fetch repository info from codecov", zap.String("repository", repository.GitRepository), zap.String("git_organization", repository.GitOrganization), zap.Error(err))

		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    "Failed to obtain repositories. There are no repository cached",
			StatusCode: http.StatusBadRequest,
		})
	}

	totalCoverageConverted, _ := coverage.Totals.Coverage.Float64()
	err = rp.Storage.CreateCoverage(storage.Coverage{
		GitOrganization:    githubRepo.Owner.GetLogin(),
		RepositoryName:     githubRepo.GetName(),
		CoveragePercentage: totalCoverageConverted,
	}, createdRepo.ID)

	if err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    "Failed to save coverage data in database. There are no repository cached",
			StatusCode: http.StatusBadRequest,
		})
	}

	workflows, err := rp.Github.GetRepositoryWorkflows(githubRepo.Owner.GetLogin(), githubRepo.GetName())
	for _, w := range workflows.Workflows {
		rp.Storage.CreateWorkflows(storage.GithubWorkflows{
			WorkflowName: w.GetName(),
			BadgeURL:     w.GetBadgeURL(),
			HTMLURL:      w.GetHTMLURL(),
			JobURL:       w.GetURL(),
			State:        w.GetState(),
		}, createdRepo.ID)
	}

	if err != nil {
		rp.Logger.Error("Failed to fetch github actions info from github", zap.String("repository", repository.GitRepository), zap.String("git_organization", repository.GitOrganization), zap.Error(err))

		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    "Failed to save workflows data in database. There are no repository cached",
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
	json.NewDecoder(r.Body).Decode(&repository)

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
