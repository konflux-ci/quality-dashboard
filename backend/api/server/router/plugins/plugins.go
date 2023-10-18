package plugins

import (
	"context"
	"encoding/json"
	"fmt"
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

// Plugins godoc
// @Summary List Plugins
// @Description returns a list of plugins created in quality studio
// @Tags Plugins API Info
// @Produce json
// @Router /plugins/hub/list [get]
// @Success 200 {object} []db.Plugins
// @Failure 400 {object} types.ErrorResponse
func (s *pluginsRouter) getAllPlugins(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	plugins, err := s.Storage.ListPlugins()
	if err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
	}

	return httputils.WriteJSON(w, http.StatusOK, plugins)
}

// Plugins godoc
// @Summary Create plugins
// @Description create a plugin in quality studio
// @Tags Plugins API Info
// @Produce json
// @Router /plugins/hub/create [post]
// @Param request body TeamPluginRequest true "Body json params"
// @Success 200 {object} db.Plugins
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
func (s *pluginsRouter) createPlugin(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	var plugin *v1alphaPlugins.PluginSpec
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

// Plugins godoc
// @Summary Assign Plugin to a team
// @Description Assign a plugin to a team in Quality Studio
// @Tags Plugins API Info
// @Produce json
// @Router /plugins/hub/install [post]
// @Param request body TeamPluginRequest true "Body json params"
// @Success 200 {object} db.Plugins
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
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
		fmt.Println(err)
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    "error installing plugin",
			StatusCode: http.StatusBadRequest,
		})
	}

	return httputils.WriteJSON(w, http.StatusOK, plugin)
}

// Plugins godoc
// @Summary Get plugin by team
// @Description Get a plugin from a given team_name in params
// @Tags Plugins API Info
// @Produce json
// @Router /plugins/hub/get/team [post]
// @Param   team_name     query     string     true  "string example"   example(string)
// @Success 200 {object} v1alphaPlugins.Plugin
// @Failure 400 {object} types.ErrorResponse
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
