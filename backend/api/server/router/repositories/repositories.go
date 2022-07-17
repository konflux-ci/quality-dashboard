package repositories

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/redhat-appstudio/quality-studio/api/types"
	"github.com/redhat-appstudio/quality-studio/pkg/storage"
	"github.com/redhat-appstudio/quality-studio/pkg/utils/httputils"
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
// @Failure 400 {object} ErrorResponse
func (rp *repositoryRouter) listAllRepositoriesQuality(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	repos, err := rp.Storage.ListRepositoriesQualityInfo()

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
// @Success 200 {object} db.Repository
// @Failure 400 {object} ErrorResponse
func (rp *repositoryRouter) createRepositoryHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := json.NewDecoder(r.Body).Decode(&repository); err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    "Error reading repository/git_organization value from body",
			StatusCode: http.StatusBadRequest,
		})
	}

	githubRepo, err := rp.Github.GetGithubRepositoryInformation(repository.GitOrganization, repository.GitRepository)
	if err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
	}

	createdRepo, err := rp.Storage.CreateRepository(storage.Repository{
		RepositoryName:  githubRepo.GetName(),
		GitOrganization: githubRepo.Owner.GetLogin(),
		Description:     githubRepo.GetDescription(),
		GitURL:          githubRepo.GetHTMLURL(),
	})
	if err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
	}
	coverage, err := rp.CodeCov.GetCodeCovInfo(githubRepo.Owner.GetLogin(), githubRepo.GetName())
	if err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    "Failed to obtain repositories. There are no repository cached",
			StatusCode: http.StatusBadRequest,
		})
	}

	totalCoverageConverted, _ := coverage.Commit.Totals.TotalCoverage.Float64()
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
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    "Failed to save workflows data in database. There are no repository cached",
			StatusCode: http.StatusBadRequest,
		})
	}

	return httputils.WriteJSON(w, http.StatusOK, createdRepo)
}

// Version godoc
// @Summary Github repositories info
// @Description delete a given repository from a organization
// @Tags Github Repositories API
// @Produce json
// @Param repository body GitRepositoryRequest true "repository name"
// @Router /repositories/delete [delete]
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
func (rp *repositoryRouter) deleteRepositoryHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	json.NewDecoder(r.Body).Decode(&repository)

	if repository.GitRepository == "" {
		return httputils.WriteJSON(w, http.StatusOK, ErrorResponse{
			Message:    "Failed to remove repository. Field 'repository_name' missing",
			StatusCode: 400,
		})
	}
	if repository.GitOrganization == "" {
		return httputils.WriteJSON(w, http.StatusOK, ErrorResponse{
			Message:    "Failed to remove repository. Field 'git_organization' missing",
			StatusCode: 400,
		})
	}
	err := rp.Storage.DeleteRepository(repository.GitRepository, repository.GitOrganization)
	if err != nil {
		return httputils.WriteJSON(w, http.StatusOK, ErrorResponse{
			Message:    "Failed to remove repository",
			StatusCode: 400,
		})
	}
	return httputils.WriteJSON(w, http.StatusOK, SuccessResponse{
		Message: "Repository deleted",
	})
}
