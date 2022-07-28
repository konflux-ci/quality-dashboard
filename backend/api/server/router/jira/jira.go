package jira

import (
	"context"
	"fmt"
	"net/http"

	"github.com/redhat-appstudio/quality-studio/pkg/utils/httputils"
)

const (
	KNOWN_E2E_ISSUE_LABEL = "appstudio-e2e-tests-known-issues"
)

// Jira godoc
// @Summary Jira API Info
// @Description returns a list of jira issues which contain the label appstudio-e2e-tests-known-issues
// @Tags Jira API Info
// @Produce json
// @Router /jira/bugs/e2e [get]
// @Success 200 {object} []jira.Issue
func (s *jiraRouter) listE2EBugsKnown(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	issuess := s.Jira.GetIssueByJQLQuery(fmt.Sprintf("labels in (appstudio-e2e-tests-known-issues) AND status not in (resolved, closed)", KNOWN_E2E_ISSUE_LABEL))

	return httputils.WriteJSON(w, http.StatusOK, issuess)
}
