package plugins

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	v1alphaPlugins "github.com/redhat-appstudio/quality-studio/api/apis/plugins/v1alpha1"
	"github.com/redhat-appstudio/quality-studio/api/types"
	"github.com/redhat-appstudio/quality-studio/pkg/utils/httputils"
	"k8s.io/utils/strings/slices"
)

var categories = []string{"Openshift CI", "Jira", "GitHub"}

type TeamPluginRequest struct {
	TeamName string `json:"team_name"`

	PluginName string `json:"plugin_name"`
}

// Teams godoc
// @Summary Teams API Info
// @Description returns a list of teams created in quality studio
// @Tags Teams API Info
// @Produce json
// @Router /teams/list/all [get]
// @Success 200 {object} []db.Teams
func (s *pluginsRouter) getAllPlugins(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	plugins, err := s.Storage.ListPlugins()
	if err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusConflict,
		})
	}

	return httputils.WriteJSON(w, http.StatusOK, plugins)
}

// Teams godoc
// @Summary Teams API Info
// @Description create a team in quality studio
// @Tags Teams API Info
// @Produce json
// @Router /teams/create [post]
// @Param request body TeamsRequest true "Body json params"
// @Success 200 {object} db.Teams
func (s *pluginsRouter) createPlugin(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	var plugin *v1alphaPlugins.Plugin
	if err := json.NewDecoder(r.Body).Decode(&plugin); err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    "Error reading plugin object from body",
			StatusCode: http.StatusBadRequest,
		})
	}

	if !slices.Contains(categories, plugin.Category) {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    "plugin have wronng category; should be one of " + strings.Join(categories, ","),
			StatusCode: http.StatusBadRequest,
		})
	}

	plugins, err := s.Storage.CreatePlugin(plugin)
	if err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		})
	}

	return httputils.WriteJSON(w, http.StatusOK, plugins)
}

// Teams godoc
// @Summary Teams API Info
// @Description create a team in quality studio
// @Tags Teams API Info
// @Produce json
// @Router /teams/create [post]
// @Param request body TeamsRequest true "Body json params"
// @Success 200 {object} db.Teams
func (s *pluginsRouter) installTeamPlugin(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	var teamPluginRequest *TeamPluginRequest
	if err := json.NewDecoder(r.Body).Decode(&teamPluginRequest); err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    "error reading plugin object from body.",
			StatusCode: http.StatusBadRequest,
		})
	}

	team, err := s.Storage.GetTeamByName(teamPluginRequest.TeamName)
	if err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
	}

	plugin, err := s.Storage.GetPluginByName(teamPluginRequest.PluginName)
	if err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		})
	}

	if _, err := s.Storage.InstallPlugin(team, plugin); err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    "error installing plugin",
			StatusCode: http.StatusBadRequest,
		})
	}

	return httputils.WriteJSON(w, http.StatusOK, plugin)
}

// Teams godoc
// @Summary Teams API Info
// @Description create a team in quality studio
// @Tags Teams API Info
// @Produce json
// @Router /teams/create [post]
// @Param request body TeamsRequest true "Body json params"
// @Success 200 {object} db.Teams
func (s *pluginsRouter) getPluginsByTeam(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
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
			StatusCode: http.StatusBadRequest,
		})
	}

	plugins, err := s.Storage.GetPluginsByTeam(team)
	if err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
	}
	return httputils.WriteJSON(w, http.StatusOK, plugins)
}
