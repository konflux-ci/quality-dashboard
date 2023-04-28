package login

import (
	"github.com/redhat-appstudio/quality-studio/api/server/router"
)

// systemRouter provides information about the server version.
type loginRouter struct {
	routes []router.Route
}

// NewRouter initializes a new system router
func NewRouter() router.Router {
	r := &loginRouter{}

	r.routes = []router.Route{
		router.NewGetRoute("/login", r.login),
		router.NewGetRoute("/login/callback", r.handleOAuth2Callback),
	}

	return r
}

// Routes returns all the API routes dedicated to the server system
func (s *loginRouter) Routes() []router.Route {
	return s.routes
}
