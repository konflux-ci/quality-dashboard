package jira

import (
	"context"
	"net/http"

	"github.com/redhat-appstudio/quality-studio/api/types"
	"github.com/redhat-appstudio/quality-studio/pkg/utils/httputils"
)

// Jira godoc
// @Summary Jira API Info
// @Description returns a list of jira issues which contain the label appstudio-e2e-tests-known-issues
// @Tags Jira API Info
// @Produce json
// @Router /jira/bugs/e2e [get]
// @Success 200 {object} []jira.Issue
func (s *jiraRouter) listE2EBugsKnown(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	issues := s.Jira.GetIssueByJQLQuery(`project in (STONE, DEVHAS, SRVKP, GITOPSRVCE, HACBS) AND status not in (Closed) AND labels = ci-fail`)

	return httputils.WriteJSON(w, http.StatusOK, issues)
}

// Jira godoc
// @Summary Jira API Info
// @Description returns all bugs stored in database
// @Tags Jira API Info
// @Produce json
// @Router /jira/bugs/all [get]
// @Success 200 {object} []db.Bugs
func (s *jiraRouter) listAllBugs(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	bugs, err := s.Storage.GetAllJiraBugs()

	if err != nil {
		s.Logger.Error("Failed to fetch bugs")

		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
	}
	return httputils.WriteJSON(w, http.StatusOK, bugs)
}
