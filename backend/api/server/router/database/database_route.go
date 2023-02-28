package database

import (
	"github.com/redhat-appstudio/quality-studio/api/server/router"
)

type databaseRouter struct {
	routes []router.Route
}

func NewRouter() router.Router {
	r := &databaseRouter{}

	r.routes = []router.Route{
		router.NewGetRoute("/database/ok", r.getDbConnection),
	}

	return r
}

func (s *databaseRouter) Routes() []router.Route {
	return s.routes
}
