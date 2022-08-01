package prow

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redhat-appstudio/quality-studio/api/types"
	"github.com/redhat-appstudio/quality-studio/pkg/storage"
	"github.com/redhat-appstudio/quality-studio/pkg/utils/httputils"
)

type GitRepositoryRequest struct {
	GitOrganization string `json:"git_organization"`
	GitRepository   string `json:"repository_name"`
}

var (
	suitesXml TestSuites
)

// version godoc
// @Summary Prow Jobs info
// @Description returns all prow jobs information stored in database
// @Tags Prow Jobs info
// @Accept json
// @Produce json
// @Param repository_name body GitRepositoryRequest true "repository name"
// @Param git_organization body GitRepositoryRequest true "repository name"
// @Param job_id body string true "repository name"
// @Router /prow/results/post [post]
// @Success 200 {Object} types.SuccessResponse
// @Failure 400 {object} types.ErrorResponse
func (s *jobRouter) createProwCIResults(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	repositoryName := r.URL.Query()["repository_name"]
	gitOrgazanitation := r.URL.Query()["git_organization"]
	jobID := r.URL.Query()["job_id"]
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
	} else if len(jobID) == 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "job_id value not present in query",
			StatusCode: 400,
		})
	} else if len(jobType) == 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "job_type value not present in query",
			StatusCode: 400,
		})
	}

	prowJobsInDatabase, _ := s.Storage.GetProwJobsResultsByJobID(jobID[0])
	if len(prowJobsInDatabase) > 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "ci jobid already exist in database",
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

	testXml, err := parseFileFromRequest(r, &s.Logger)

	if err != nil {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    fmt.Sprintf("Error parsing junit file %s", err),
			StatusCode: 400,
		})
	}
	for _, suites := range testXml.Suites {
		s.Storage.CreateProwJobResults(storage.ProwJobStatus{
			JobID:        jobID[0],
			CreatedAt:    time.Now(),
			Duration:     suites.Duration,
			TestsCount:   int64(suites.NumTests),
			FailedCount:  int64(suites.NumFailed),
			SkippedCount: int64(suites.NumSkipped),
			JobType:      jobType[0],
		}, repoInfo.ID)
		for _, testCase := range suites.TestCases {
			err := s.Storage.CreateProwJobSuites(storage.ProwJobSuites{
				JobID:          jobID[0],
				TestCaseName:   testCase.Name,
				TestCaseStatus: testCase.Status,
				TestTiming:     testCase.Duration,
				JobType:        jobType[0],
			}, repoInfo.ID)

			if err != nil {
				s.Logger.Sugar().Errorf("Failed to save test case %s: %s", testCase.Name, err)
			}
		}
	}

	return httputils.WriteJSON(w, http.StatusOK, types.SuccessResponse{
		Message:    "Successfully stored Prow Job",
		StatusCode: http.StatusCreated,
	})
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
