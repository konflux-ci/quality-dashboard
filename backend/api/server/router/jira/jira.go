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
	issues := s.Jira.GetIssueByJQLQuery(`project in (DEVHAS, SRVKP, GITOPSRVCE, HACBS, RHTAP, RHTAPBUGS) AND status not in (Closed) AND labels = ci-fail`)

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
	TeamName  string `json:"team_name"`
	Priority  string `json:"priority"`
	StartDate string `json:"start"`
	EndDate   string `json:"end"`
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

	team, err := s.Storage.GetTeamByName(resolution.TeamName)
	if err != nil {
		s.Logger.Error("Failed to fetch bugs")

		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
	}
	totalAvg, err := s.Storage.TotalBugsResolutionTime(resolution.Priority, resolution.StartDate, resolution.EndDate, team)

	if err != nil {
		s.Logger.Error("Failed to fetch bugs")

		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
	}

	return httputils.WriteJSON(w, http.StatusOK, totalAvg)
}

// Jira godoc
// @Summary Jira API Info
// @Description returns all bugs stored in database
// @Tags Jira API Info
// @Produce json
// @Router /jira/bugs/resolution [post]
// @Success 200 {object} v1alpha1.BugsMetrics
func (s *jiraRouter) openBugsMetrics(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	var resolution ResolutionRequest
	if err := json.NewDecoder(r.Body).Decode(&resolution); err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    "Error reading priority value from body",
			StatusCode: http.StatusBadRequest,
		})
	}

	if resolution.Priority == "" {
		resolution.Priority = "Global"
	}

	team, err := s.Storage.GetTeamByName(resolution.TeamName)
	if err != nil {
		s.Logger.Error("Failed to fetch bugs")

		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
	}

	openBugMetrics, err := s.Storage.GetOpenBugsMetricsByStatusAndPriority(resolution.Priority, resolution.StartDate, resolution.EndDate, team)

	if err != nil {
		s.Logger.Error("Failed to fetch bugs")

		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
	}

	return httputils.WriteJSON(w, http.StatusOK, openBugMetrics)
}

type BugCategoriesMetrics struct {
	Priority     string `json:"priority"`
	OpenBugs     int    `json:"open_bugs"`
	ResolvedBugs int    `json:"resolved_bugs"`
}

// Jira godoc
// @Summary Jira API Info
// @Description returns all bugs stored in database
// @Tags Jira API Info
// @Produce json
// @Router /jira/bugs/metrics/priorities [get]
// @Success 200 {object} v1alpha1.BugsMetrics
func (s *jiraRouter) getCountBugsForAllCategories(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	teamName := r.URL.Query()["team_name"]
	startDate := r.URL.Query()["start_date"]
	endDate := r.URL.Query()["end_date"]
	bugByCategory := []BugCategoriesMetrics{}
	priorities := []string{"Global", "Blocker", "Critical", "Major", "Normal", "Minor", "Undefined"}

	if len(teamName) == 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "team_name value not present in query",
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

	team, err := s.Storage.GetTeamByName(teamName[0])
	if err != nil {
		s.Logger.Error("Failed to fetch bugs")

		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
	}

	for _, priority := range priorities {
		totalAvg, err := s.Storage.TotalBugsResolutionTime(priority, startDate[0], endDate[0], team)
		if err != nil {
			s.Logger.Error("Failed to fetch bugs")

			return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
			})
		}

		openBugMetrics, err := s.Storage.GetOpenBugsMetricsByStatusAndPriority(priority, startDate[0], endDate[0], team)

		if err != nil {
			s.Logger.Error("Failed to fetch bugs")

			return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
			})
		}

		bugByCategory = append(bugByCategory, BugCategoriesMetrics{
			Priority:     priority,
			OpenBugs:     openBugMetrics.TotalOpenBugs.NumberOfOpenBugs,
			ResolvedBugs: totalAvg.ResolutionTimeTotal.NumberOfTotalBugs,
		})
	}

	return httputils.WriteJSON(w, http.StatusOK, bugByCategory)
}

type ProjectListSimplified struct {
	ProjectKey string `json:"project_key"`

	ProjectName string `json:"project_name"`
}

// Jira godoc
// @Summary Jira API Info
// @Description returns all bugs stored in database
// @Tags Jira API Info
// @Produce json
// @Router /jira/project/list [get]
// @Success 200 {object} []ProjectListSimplified
func (s *jiraRouter) getJiraProjects(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	projects := []ProjectListSimplified{}

	list, err := s.Jira.GetJiraProjects()
	if err != nil {
		s.Logger.Error("Failed to fetch jira projects")

		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
	}

	for _, p := range *list {
		projects = append(projects, ProjectListSimplified{
			ProjectKey:  p.Key,
			ProjectName: p.Name,
		})
	}

	return httputils.WriteJSON(w, http.StatusOK, projects)
}

func (s *jiraRouter) bugExists(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	teamName := r.URL.Query()["team_name"]
	key := r.URL.Query()["jira_key"]

	if len(teamName) == 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "team_name value not present in query",
			StatusCode: 400,
		})
	} else if len(key) == 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "jira_key value not present in query",
			StatusCode: 400,
		})
	}

	team, err := s.Storage.GetTeamByName(teamName[0])
	if err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
	}

	exists, err := s.Storage.BugExists(key[0], team)
	if err != nil {
		s.Logger.Error("Failed to check if jira key exists")

		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
	}

	return httputils.WriteJSON(w, http.StatusOK, exists)
}
