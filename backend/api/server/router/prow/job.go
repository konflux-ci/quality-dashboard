package prow

import (
	"context"
	"fmt"
	"net/http"

	"github.com/redhat-appstudio/quality-studio/api/types"
	"github.com/redhat-appstudio/quality-studio/pkg/utils/httputils"
)

type GitRepositoryRequest struct {
	GitOrganization string `json:"git_organization"`
	GitRepository   string `json:"repository_name"`
}

// version godoc
// @Summary Prow Jobs info
// @Description returns all prow jobs related to git_organization and repository_name
// @Tags Prow Jobs info
// @Accept json
// @Produce json
// @Param repository body GitRepositoryRequest true "repository name"
// @Param organization body GitRepositoryRequest true "git_organization"
// @Router /prow/results/get [get]
// @Success 200 {Object} []db.Prow
// @Failure 400 {object} types.ErrorResponse
func (s *jobRouter) getProwJobs(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	repositoryName := r.URL.Query()["repository_name"]
	gitOrgazanitation := r.URL.Query()["git_organization"]

	if len(repositoryName) == 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "repository_name value not present in query",
			StatusCode: 400,
		})
	} else if len(gitOrgazanitation) == 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "git_organization value not present in query",
			StatusCode: 400,
		})
	}

	repoInfo, err := s.Storage.GetRepository(repositoryName[0], gitOrgazanitation[0])
	if err != nil {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "failed to get repository from database; check if repository exist in quality studio",
			StatusCode: 400,
		})
	}

	prows, err := s.Storage.GetProwJobsResults(repoInfo)
	if err != nil {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "failed to get repository from database; check if repository exist in quality studio",
			StatusCode: 400,
		})
	}

	return httputils.WriteJSON(w, http.StatusOK, prows)
}

// version godoc
// @Summary Prow Jobs info
// @Description returns all prow jobs related to git_organization and repository_name
// @Tags Prow Jobs info
// @Accept json
// @Produce json
// @Param repository body GitRepositoryRequest true "repository name"
// @Param organization body GitRepositoryRequest true "git_organization"
// @Router /prow/results/latest/get [get]
// @Success 200 {Object} []db.ProwJobSuites
// @Failure 400 {object} types.ErrorResponse
func (s *jobRouter) getLatestSuitesExecution(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	repositoryName := r.URL.Query()["repository_name"]
	gitOrgazanitation := r.URL.Query()["git_organization"]
	jobType := r.URL.Query()["job_type"]

	if len(repositoryName) == 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "repository_name value not present in query",
			StatusCode: 400,
		})
	} else if len(gitOrgazanitation) == 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "git_organization value not present in query",
			StatusCode: 400,
		})
	} else if len(jobType) == 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "job_type value not present in query",
			StatusCode: 400,
		})
	}

	repoInfo, err := s.Storage.GetRepository(repositoryName[0], gitOrgazanitation[0])
	if err != nil {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    fmt.Sprintf("Repository '%s' doesn't exists in quality studio database", repositoryName[0]),
			StatusCode: 400,
		})
	}
	latest, err := s.Storage.GetLatestProwTestExecution(repoInfo, jobType[0])
	if err != nil {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "failed to get latest prow execution.",
			StatusCode: 400,
		})
	}

	suiterlandia, err := s.Storage.GetSuitesByJobID(latest.JobID)

	if err != nil {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "failed to get latest prow execution",
			StatusCode: 400,
		})
	}

	return httputils.WriteJSON(w, http.StatusOK, suiterlandia)
}
