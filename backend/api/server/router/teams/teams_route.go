package teams

import (
	"github.com/redhat-appstudio/quality-studio/api/server/router"
	"github.com/redhat-appstudio/quality-studio/pkg/logger"
	"github.com/redhat-appstudio/quality-studio/pkg/storage"
	"go.uber.org/zap"
)

// systemRouter provides information about the server version.
type teamsRouter struct {
	Route   []router.Route
	Storage storage.Storage
	Logger  *zap.Logger
}

// NewRouter initializes a new system router
func NewRouter(s storage.Storage) router.Router {
	logger, _ := logger.InitZap("info")
	r := &teamsRouter{}

	r.Storage = s
	r.Logger = logger

	r.Route = []router.Route{
		router.NewGetRoute("/teams/list/all", r.listAllQualityStudioTeams),
		router.NewPostRoute("/teams/create", r.createQualityStudioTeam),
	}

	return r
}

// Routes returns all the API routes dedicated to the server system
func (s *teamsRouter) Routes() []router.Route {
	return s.Route
}
