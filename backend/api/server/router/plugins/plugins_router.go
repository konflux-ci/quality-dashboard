package plugins

import (
	"github.com/redhat-appstudio/quality-studio/api/server/router"
	"github.com/redhat-appstudio/quality-studio/pkg/storage"
)

type pluginsRouter struct {
	Storage storage.Storage
	routes  []router.Route
}

func NewRouter(s storage.Storage) router.Router {
	r := &pluginsRouter{}
	r.Storage = s
	r.routes = []router.Route{
		router.NewGetRoute("/plugins/hub/list", r.getAllPlugins),
		router.NewPostRoute("/plugins/hub/create", r.createPlugin),
		router.NewPostRoute("/plugins/hub/install", r.installTeamPlugin),
		router.NewGetRoute("/plugins/hub/get/team", r.getPluginsByTeam),
	}

	return r
}

func (s *pluginsRouter) Routes() []router.Route {
	return s.routes
}
