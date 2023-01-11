package teams

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/redhat-appstudio/quality-studio/api/types"
	"github.com/redhat-appstudio/quality-studio/pkg/utils/httputils"
)

type TeamsRequest struct {
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
