package teams

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/redhat-appstudio/quality-studio/api/types"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
	"github.com/redhat-appstudio/quality-studio/pkg/utils/httputils"
)

type TeamsRequest struct {
	TeamName    string `json:"team_name"`
	Description string `json:"description"`
}

type UpdateTeamsRequest struct {
	TargetTeam  string `json:"target"`
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
			Message:    "Error reading team_name/description value from body",
			StatusCode: http.StatusBadRequest,
		})
	}

	teams, err := s.Storage.CreateQualityStudioTeam(team.TeamName, team.Description)
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
// @Description delete a team in quality studio
// @Tags Teams API Info
// @Produce json
// @Param repository body TeamsRequest true "team_name"
// @Router /teams/delete [delete]
// @Success 200 {object} types.SuccessResponse
// @Failure 400 {object} types.ErrorResponse
func (rp *teamsRouter) deleteTeamHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	var team TeamsRequest
	json.NewDecoder(r.Body).Decode(&team)

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
	json.NewDecoder(r.Body).Decode(&team)

	if team.TeamName == "" {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "Failed to remove team. Field 'team_name' missing",
			StatusCode: 400,
		})
	}

	err := rp.Storage.UpdateTeam(
		&db.Teams{
			TeamName:    team.TeamName,
			Description: team.Description,
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
