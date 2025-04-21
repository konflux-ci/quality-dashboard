package prow

import (
	"context"
	"fmt"
	"net/http"

	prowv1Alpha1 "github.com/konflux-ci/quality-dashboard/api/apis/prow/v1alpha1"
	"github.com/konflux-ci/quality-dashboard/api/types"
	"github.com/konflux-ci/quality-dashboard/pkg/utils/httputils"
	"go.uber.org/zap"
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
	gitOrganization := r.URL.Query()["git_organization"]
	startDate := r.URL.Query()["start_date"]
	endDate := r.URL.Query()["end_date"]

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

	repoInfo, err := s.Storage.GetRepository(repositoryName[0], gitOrganization[0])
	if err != nil {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "failed to get repository from database; check if repository exist in quality studio",
			StatusCode: 400,
		})
	}

	prows, err := s.Storage.GetProwJobsResults(repoInfo, startDate[0], endDate[0])
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
	gitOrganization := r.URL.Query()["git_organization"]
	jobType := r.URL.Query()["job_type"]

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
	} else if len(jobType) == 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "job_type value not present in query",
			StatusCode: 400,
		})
	}

	repoInfo, err := s.Storage.GetRepository(repositoryName[0], gitOrganization[0])
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

// version godoc
// @Summary Github repositories info
// @Description returns all repository information stored in database
// @Tags Github Repositories API
// @Accept json
// @Produce json
// @Router /prow/repositories/list [get]
// @Success 200 {array} storage.RepositoryQualityInfo
// @Failure 400 {object} types.ErrorResponse
func (jb *jobRouter) listProwRepos(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "team_name value not present in query",
			StatusCode: http.StatusBadRequest,
		})
	}

	team, err := jb.Storage.GetTeamByName(teamName)
	if err != nil {
		jb.Logger.Error("Failed to fetch team", zap.String("team", teamName), zap.Error(err))
		return httputils.WriteJSON(w, http.StatusInternalServerError, types.ErrorResponse{
			Message:    "Unable to fetch team info",
			StatusCode: http.StatusInternalServerError,
		})
	}

	repos, err := jb.Storage.ListRepositories(team)
	if err != nil {
		jb.Logger.Error("Failed to fetch repositories", zap.String("team", teamName), zap.Error(err))
		return httputils.WriteJSON(w, http.StatusInternalServerError, types.ErrorResponse{
			Message:    "Unable to list repositories",
			StatusCode: http.StatusInternalServerError,
		})
	}

	var prowRepos []prowv1Alpha1.ProwRepository
	for _, repo := range repos {
		repoInfo, err := jb.Storage.GetRepository(repo.Name, repo.Owner.Login)
		if err != nil {
			jb.Logger.Warn("Skipping repository - details not found", zap.String("repo", repo.Name), zap.Error(err))
			continue
		}

		jobs, err := jb.Storage.GetProwJobsByRepoOrg(repoInfo)
		if err != nil {
			jb.Logger.Error("Failed to fetch Prow jobs", zap.String("repository", repo.Name), zap.Error(err))
			return httputils.WriteJSON(w, http.StatusInternalServerError, types.ErrorResponse{
				Message:    "Error fetching Prow jobs",
				StatusCode: http.StatusInternalServerError,
			})
		}

		if len(jobs) == 0 {
			continue
		}

		prowRepos = append(prowRepos, prowv1Alpha1.ProwRepository{
			Repository: repo,
			JobsList:   jobs,
		})
	}

	return httputils.WriteJSON(w, http.StatusOK, prowRepos)
}

func (jb *jobRouter) getProwJobsByType(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	repositoryName := r.URL.Query()["repository_name"]
	gitOrganization := r.URL.Query()["git_organization"]

	if len(gitOrganization) == 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "git_organization value not present in query",
			StatusCode: 400,
		})
	} else if len(gitOrganization) == 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "repository_name value not present in query",
			StatusCode: 400,
		})
	}

	repoInfo, err := jb.Storage.GetRepository(repositoryName[0], gitOrganization[0])
	if err != nil {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    fmt.Sprintf("Repository '%s' doesn't exists in quality studio database", repositoryName[0]),
			StatusCode: 400,
		})
	}

	jobsAndType, err := jb.Storage.GetJobsNameAndType(repoInfo)
	if err != nil {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    fmt.Sprintf("Unknown error requesting repository '%s'", repositoryName[0]),
			StatusCode: 500,
		})
	}

	return httputils.WriteJSON(w, http.StatusOK, jobsAndType)
}
