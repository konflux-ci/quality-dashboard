package prow

import (
	"context"
	"fmt"
	"net/http"

	prowv1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/prow/v1alpha1"
	"github.com/redhat-appstudio/quality-studio/api/types"
	"github.com/redhat-appstudio/quality-studio/pkg/utils/httputils"
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

	repoInfo, err := s.Storage.GetRepository(repositoryName[0], gitOrganization[0])
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
	teamName := r.URL.Query()["team_name"]

	if len(teamName) == 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "team_name value not present in query",
			StatusCode: 400,
		})
	}

	team, err := jb.Storage.GetTeamByName(teamName[0])
	if err != nil {
		jb.Logger.Error("Failed to fetch team. Make sure the team exists", zap.String("team", teamName[0]), zap.Error(err))

		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
	}

	repos, err := jb.Storage.ListRepositories(team)
	if err != nil {
		jb.Logger.Error("Failed to fetch repositories", zap.String("team", teamName[0]), zap.Error(err))

		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
	}

	prowRepos := make([]prowv1Alpha1.ProwRepository, 0)

	for _, repo := range repos {
		repoInfo, err := jb.Storage.GetRepository(repo.Name, repo.Owner.Login)
		if err != nil {
			continue
		}

		if jb.Github.CheckIfRepoExistsInOpenshiftCI(repo.Owner.Login, repo.Name) {
			list, err := jb.Storage.GetProwJobsByRepoOrg(repoInfo)
			if err != nil {
				jb.Logger.Error("Failed to fetch team. Make sure the team exists", zap.String("team", teamName[0]), zap.Error(err))

				return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
					Message:    err.Error(),
					StatusCode: http.StatusBadRequest,
				})
			}

			prowRepos = append(prowRepos, prowv1Alpha1.ProwRepository{
				Repository: repo,
				JobsList:   list,
			})

		}
	}

	return httputils.WriteJSON(w, http.StatusOK, prowRepos)
}
