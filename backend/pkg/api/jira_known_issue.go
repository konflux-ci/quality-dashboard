package api

import (
	"fmt"
	"net/http"
)

const (
	KNOWN_E2E_ISSUE_LABEL = "appstudio-e2e-tests-known-issues"
)

// Jira godoc
// @Summary Jira
// @Description returns a list of jira issues which contain the label appstudio-e2e-tests-known-issues
// @Tags Version API
// @Produce json
// @Router /api/jira/e2e-known/get [get]
// @Success 200 {object} api.MapResponse
func (s *Server) getE2eKnownIssues(w http.ResponseWriter, r *http.Request) {
	issuess := s.JiraApi.GetIssueByJQLQuery(fmt.Sprintf("labels in (%s) AND status not in (resolved, closed)", KNOWN_E2E_ISSUE_LABEL))

	s.JSONResponse(w, r, issuess)
}
