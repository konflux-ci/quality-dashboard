package version

import (
	"github.com/konflux-ci/quality-dashboard/api/server/router"
)

// systemRouter provides information about the server version.
type versionRouter struct {
	routes []router.Route
}

// NewRouter initializes a new system router
func NewRouter() router.Router {
	r := &versionRouter{}

	r.routes = []router.Route{
		router.NewGetRoute("/server/info", r.getVersion),
	}

	return r
}

// Routes returns all the API routes dedicated to the server system
func (s *versionRouter) Routes() []router.Route {
	return s.routes
}
