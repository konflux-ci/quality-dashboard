package jira

import (
	"context"
	"encoding/json"
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

type ResolutionRequest struct {
	Priority string `json:"priority"`
}

// Jira godoc
// @Summary Jira API Info
// @Description returns all bugs stored in database
// @Tags Jira API Info
// @Produce json
// @Router /jira/bugs/resolution [post]
// @Success 200 {object} v1alpha1.BugsMetrics
func (s *jiraRouter) calculateRates(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	var resolution ResolutionRequest
	if err := json.NewDecoder(r.Body).Decode(&resolution); err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    "Error reading team_name/description value from body",
			StatusCode: http.StatusBadRequest,
		})
	}

	if resolution.Priority == "" {
		resolution.Priority = "Global"
	}

	totalAvg, err := s.Storage.TotalBugsResolutionTime(resolution.Priority)

	if err != nil {
		s.Logger.Error("Failed to fetch bugs")

		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
	}

	return httputils.WriteJSON(w, http.StatusOK, totalAvg)
}
