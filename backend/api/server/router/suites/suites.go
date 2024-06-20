package suites

import (
	"context"
	"net/http"

	"github.com/konflux-ci/quality-dashboard/api/types"
	"github.com/konflux-ci/quality-dashboard/pkg/utils/httputils"
)

func (s *suitesRouter) getOccurrences(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	gitOrganization := r.URL.Query()["git_org"]
	repositoryName := r.URL.Query()["repository_name"]
	jobName := r.URL.Query()["job_name"]
	startDate := r.URL.Query()["start_date"]
	endDate := r.URL.Query()["end_date"]

	if len(repositoryName) == 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "repository_name value not present in query",
			StatusCode: 400,
		})
	} else if len(gitOrganization) == 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "git_org value not present in query",
			StatusCode: 400,
		})
	} else if len(jobName) == 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "job_name value not present in query",
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

	freq, err := s.Storage.GetSuitesFailureFrequency(gitOrganization[0], repositoryName[0], jobName[0], startDate[0], endDate[0])
	if err != nil {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "failed to get flaky suites metrics.",
			StatusCode: 400,
		})
	}

	return httputils.WriteJSON(w, http.StatusOK, freq)
}

func (s *suitesRouter) getTrends(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	gitOrganization := r.URL.Query()["git_org"]
	repositoryName := r.URL.Query()["repository_name"]
	jobName := r.URL.Query()["job_name"]
	startDate := r.URL.Query()["start_date"]
	endDate := r.URL.Query()["end_date"]

	if len(repositoryName) == 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "repository_name value not present in query",
			StatusCode: 400,
		})
	} else if len(gitOrganization) == 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "git_org value not present in query",
			StatusCode: 400,
		})
	} else if len(jobName) == 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "job_name value not present in query",
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

	freq := s.Storage.GetProwFlakyTrendsMetrics(gitOrganization[0], repositoryName[0], jobName[0], startDate[0], endDate[0])

	return httputils.WriteJSON(w, http.StatusOK, freq)
}
