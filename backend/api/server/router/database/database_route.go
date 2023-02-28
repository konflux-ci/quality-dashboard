package database

import (
	"database/sql"

	"github.com/redhat-appstudio/quality-studio/api/server/router"
)

type databaseRouter struct {
	routes []router.Route
	db     *sql.DB
}

func NewRouter(db *sql.DB) router.Router {
	r := &databaseRouter{}
	r.db = db
	r.routes = []router.Route{
		router.NewGetRoute("/database/ok", r.getDbConnection),
	}

	return r
}

func (s *databaseRouter) Routes() []router.Route {
	return s.routes
}
