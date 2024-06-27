package repositories

import (
	"context"
	"net/http"

	"github.com/konflux-ci/quality-dashboard/api/types"
	"github.com/konflux-ci/quality-dashboard/pkg/utils/httputils"
)

func (rp *repositoryRouter) getPullRequestsFromRepo(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
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

	prs, err := rp.Storage.GetPullRequestsByRepository(repositoryName[0], gitOrganization[0], startDate[0], endDate[0])

	if err != nil {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "Failed to get pull requests by repository.",
			StatusCode: 400,
		})

	}
	return httputils.WriteJSON(w, http.StatusOK, prs)
}
