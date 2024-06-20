package repositories

import (
	"context"
	"net/http"

	"github.com/konflux-ci/quality-studio/api/types"
	"github.com/konflux-ci/quality-studio/pkg/utils/httputils"
)

// Version godoc
// @Summary Quality Repositories Workflow
// @Description return github workflows from a given repository
// @Tags Github Workflows API
// @Produce json
// @Router /workflows/get [get]
// @Param   repository_name     query     string     false  "string example"   example(string)
// @Success 200 {array} storage.GithubWorkflows
// @Failure 400 {object} types.ErrorResponse
func (rp *repositoryRouter) getWorkflowByRepositoryName(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	repositoryName := r.URL.Query()["repository_name"]
	workflows, err := rp.Storage.ListWorkflowsByRepository(repositoryName[0])
	if err != nil {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "Failed to remove repository. Field 'repository_name' missing",
			StatusCode: 400,
		})

	}
	return httputils.WriteJSON(w, http.StatusOK, workflows)
}
