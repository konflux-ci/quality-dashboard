package teams

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/andygrunwald/go-jira"
	configurationV1Alpha1 "github.com/konflux-ci/quality-dashboard/api/apis/configuration/v1alpha1"
	"github.com/konflux-ci/quality-dashboard/api/types"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db"
	"github.com/konflux-ci/quality-dashboard/pkg/utils/httputils"
)

type TeamsRequest struct {
	TeamName    string `json:"team_name"`
	JiraKeys    string `json:"jira_keys,omitempty"`
	Description string `json:"description"`
	JiraConfig  string `json:"jira_config"`
}

type UpdateTeamsRequest struct {
	TargetTeam  string `json:"target"`
	JiraKeys    string `json:"jira_keys"`
	TeamName    string `json:"team_name"`
	Description string `json:"description"`
	JiraConfig  string `json:"jira_config"`
}

// Teams godoc
// @Summary Teams API Info
// @Description returns a list of teams created in quality studio
// @Tags Teams API Info
// @Produce json
// @Router /teams/list/all [get]
// @Success 200 {object} []db.Teams
func (s *teamsRouter) listAllQualityStudioTeams(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	teams, err := s.Storage.GetAllTeamsFromDB()
	if err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		})
	}

	return httputils.WriteJSON(w, http.StatusOK, teams)
}

// Teams godoc
// @Summary Teams API Info
// @Description create a team in quality studio
// @Tags Teams API Info
// @Produce json
// @Router /teams/create [post]
// @Param request body TeamsRequest true "Body json params"
// @Success 200 {object} db.Teams
func (s *teamsRouter) createQualityStudioTeam(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	var team TeamsRequest
	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    "Error reading team_name/description/jira_keys value from body",
			StatusCode: http.StatusBadRequest,
		})
	}

	jiraCfg, err := parseJiraConfig(team.JiraConfig)
	if err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		})
	}

	teams, err := s.Storage.CreateQualityStudioTeam(team.TeamName, team.Description, team.JiraKeys)
	if err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		})
	}

	if teams.JiraKeys != "" {
		bugs := s.Jira.GetBugsByJQLQuery(jiraCfg.BugsCollectQuery)
		if err := s.Storage.CreateJiraBug(bugs, teams); err != nil {
			return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
				Message:    err.Error(),
				StatusCode: http.StatusInternalServerError,
			})
		}

		jiraCfg.CiImpactBugs = extractKeys(s.Jira.GetBugsByJQLQuery(jiraCfg.CiImpactQuery))
	}

	if err := s.saveConfiguration(team.TeamName, jiraCfg); err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		})
	}

	return httputils.WriteJSON(w, http.StatusOK, teams)
}

// Teams godoc
// @Summary Teams API Info
// @Description delete a team in quality studio
// @Tags Teams API Info
// @Produce json
// @Param repository body TeamsRequest true "team_name"
// @Router /teams/delete [delete]
// @Success 200 {object} types.SuccessResponse
// @Failure 400 {object} types.ErrorResponse
func (rp *teamsRouter) deleteTeamHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	var team TeamsRequest
	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "incorrect data received to server",
			StatusCode: 400,
		})
	}

	if team.TeamName == "" {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "Failed to remove team. Field 'team_name' missing",
			StatusCode: 400,
		})
	}

	_, err := rp.Storage.DeleteTeam(team.TeamName)
	if err != nil {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "Failed to remove team",
			StatusCode: 400,
		})
	}
	return httputils.WriteJSON(w, http.StatusOK, types.SuccessResponse{
		Message:    "Team deleted",
		StatusCode: 200,
	})
}

// Teams godoc
// @Summary Teams API Info
// @Description update a team in quality studio
// @Tags Teams API Info
// @Produce json
// @Param repository body TeamsRequest true "team_name"
// @Router /teams/update [put]
// @Success 200 {object} types.SuccessResponse
// @Failure 400 {object} types.ErrorResponse
func (rp *teamsRouter) updateTeamHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	var team UpdateTeamsRequest

	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "incorrect data received to server",
			StatusCode: 400,
		})
	}

	if team.TeamName == "" {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "Failed to remove team. Field 'team_name' missing",
			StatusCode: 400,
		})
	}

	t, err := rp.Storage.GetTeamByName(team.TeamName)
	if err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		})
	}

	// parse jira config
	jiraCfg, err := parseJiraConfig(team.JiraConfig)
	if err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		})
	}
	if jiraCfg.BugsCollectQuery == "" {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "Failed to get the jira config. Field 'bugs_collect_query' missing",
			StatusCode: 400,
		})
	}

	// delete all
	if err := rp.Storage.DeleteAllJiraBugByTeam(t.TeamName); err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		})
	}

	if team.JiraKeys != "" {
		// create jira issues
		bugs := rp.Jira.GetBugsByJQLQuery(jiraCfg.BugsCollectQuery)
		if err := rp.Storage.CreateJiraBug(bugs, t); err != nil {
			return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
				Message:    err.Error(),
				StatusCode: http.StatusInternalServerError,
			})
		}

		// get CI Impact bugs
		jiraCfg.CiImpactBugs = extractKeys(rp.Jira.GetBugsByJQLQuery(jiraCfg.CiImpactQuery))
	}

	err = rp.Storage.UpdateTeam(
		&db.Teams{
			TeamName:    team.TeamName,
			Description: team.Description,
			JiraKeys:    team.JiraKeys,
		},
		team.TargetTeam,
	)
	if err != nil {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "Message: " + err.Error(),
			StatusCode: 400,
		})
	}

	err = rp.saveConfiguration(team.TeamName, jiraCfg)
	if err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		})
	}

	return httputils.WriteJSON(w, http.StatusOK, types.SuccessResponse{
		Message:    "Team update",
		StatusCode: 200,
	})
}

func contains(key string, keys []string) bool {
	for _, k := range keys {
		if k == key {
			return true
		}
	}

	return false
}

func getProjectsToAdd(old, update []string) string {
	projectsToAdd := make([]string, 0)
	for _, p := range update {
		if !contains(p, old) {
			projectsToAdd = append(projectsToAdd, p)
		}
	}

	return strings.Join(projectsToAdd, ",")
}

func (s *teamsRouter) getConfiguration(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	teamName := r.URL.Query()["team_name"]
	if len(teamName) == 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "team_name value not present in query",
			StatusCode: 400,
		})
	}

	teams, err := s.Storage.GetConfiguration(teamName[0])
	if err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		})
	}

	return httputils.WriteJSON(w, http.StatusOK, teams)
}

func (s *teamsRouter) getJiraKeys(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	teamName := r.URL.Query()["team_name"]
	if len(teamName) == 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "team_name value not present in query",
			StatusCode: 400,
		})
	}

	team, err := s.Storage.GetTeamByName(teamName[0])
	if err != nil {
		s.Logger.Error("Failed to fetch team")

		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
	}

	return httputils.WriteJSON(w, http.StatusOK, team.JiraKeys)
}

func (s *teamsRouter) getTeam(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	teamName := r.URL.Query()["team_name"]
	if len(teamName) == 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "team_name value not present in query",
			StatusCode: 400,
		})
	}

	team, err := s.Storage.GetTeamByName(teamName[0])
	if err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		})
	}

	return httputils.WriteJSON(w, http.StatusOK, team)
}

// parseJiraConfig parses the Jira configuration from a JSON string.
func parseJiraConfig(configStr string) (configurationV1Alpha1.JiraConfig, error) {
	var jiraCfg configurationV1Alpha1.JiraConfig
	err := json.Unmarshal([]byte(configStr), &jiraCfg)
	return jiraCfg, err
}

// marshalJiraConfig converts the Jira configuration to a JSON string.
func marshalJiraConfig(jiraCfg configurationV1Alpha1.JiraConfig) string {
	cfg, _ := json.Marshal(jiraCfg)
	return string(cfg)
}

// extractKeys extracts the keys from a list of issues.
func extractKeys(bugs []jira.Issue) []string {
	keys := make([]string, len(bugs))
	for i, bug := range bugs {
		keys[i] = bug.Key
	}
	return keys
}

// saveConfiguration saves the configuration for the team
func (s *teamsRouter) saveConfiguration(teamName string, jiraCfg configurationV1Alpha1.JiraConfig) error {
	config := configurationV1Alpha1.Configuration{
		TeamName:      teamName,
		JiraConfig:    marshalJiraConfig(jiraCfg),
		BugSLOsConfig: "",
	}
	return s.Storage.CreateConfiguration(config)
}
