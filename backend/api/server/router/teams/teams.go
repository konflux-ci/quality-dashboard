package teams

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/konflux-ci/quality-studio/api/types"
	"github.com/konflux-ci/quality-studio/pkg/storage/ent/db"
	"github.com/konflux-ci/quality-studio/pkg/utils/httputils"
)

type TeamsRequest struct {
	TeamName    string `json:"team_name"`
	JiraKeys    string `json:"jira_keys,omitempty"`
	Description string `json:"description"`
}

type UpdateTeamsRequest struct {
	TargetTeam  string `json:"target"`
	JiraKeys    string `json:"jira_keys"`
	TeamName    string `json:"team_name"`
	Description string `json:"description"`
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

	teams, err := s.Storage.CreateQualityStudioTeam(team.TeamName, team.Description, team.JiraKeys)
	if err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		})
	}
	if teams.JiraKeys != "" {
		bugs := s.Jira.GetBugsByJQLQuery(fmt.Sprintf("project in (%s) AND type = Bug", teams.JiraKeys))
		if err := s.Storage.CreateJiraBug(bugs, teams); err != nil {
			return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
				Message:    err.Error(),
				StatusCode: http.StatusInternalServerError,
			})
		}
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

	if team.JiraKeys != "" {
		old := strings.Split(t.JiraKeys, ",")
		update := strings.Split(team.JiraKeys, ",")

		// delete team jira keys that are not present in the update ones
		for _, oldKey := range old {
			if !contains(oldKey, update) {
				if err := rp.Storage.DeleteJiraBugsByProject(oldKey, t); err != nil {
					return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
						Message:    err.Error(),
						StatusCode: http.StatusInternalServerError,
					})
				}
			}
		}

		// create jira keys
		projectsToAdd := getProjectsToAdd(old, update)
		if projectsToAdd != "" {
			bugs := rp.Jira.GetBugsByJQLQuery(fmt.Sprintf("project in (%s) AND type = Bug", projectsToAdd))
			if err := rp.Storage.CreateJiraBug(bugs, t); err != nil {
				return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
					Message:    err.Error(),
					StatusCode: http.StatusInternalServerError,
				})
			}
		}
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
