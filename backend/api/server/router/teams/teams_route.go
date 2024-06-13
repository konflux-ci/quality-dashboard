package teams

import (
	"github.com/redhat-appstudio/quality-studio/api/server/router"
	jiraAPI "github.com/redhat-appstudio/quality-studio/pkg/connectors/jira"
	"github.com/redhat-appstudio/quality-studio/pkg/logger"
	"github.com/redhat-appstudio/quality-studio/pkg/storage"
	"go.uber.org/zap"
)

// systemRouter provides information about the server version.
type teamsRouter struct {
	Route   []router.Route
	Storage storage.Storage
	Logger  *zap.Logger
	Jira    jiraAPI.Jira
}

// NewRouter initializes a new system router
func NewRouter(s storage.Storage) router.Router {
	jiraAPI := jiraAPI.NewJiraConfig()
	logger, _ := logger.InitZap("info")
	r := &teamsRouter{}

	r.Storage = s
	r.Logger = logger
	r.Jira = jiraAPI

	r.Route = []router.Route{
		router.NewGetRoute("/teams/list/all", r.listAllQualityStudioTeams),
		router.NewPostRoute("/teams/create", r.createQualityStudioTeam),
		router.NewDeleteRoute("/teams/delete", r.deleteTeamHandler),
		router.NewPutRoute("/teams/put", r.updateTeamHandler),
		router.NewGetRoute("/teams/get/configuration", r.getConfiguration),
		router.NewGetRoute("/teams/get/jira_keys", r.getJiraKeys),
		router.NewGetRoute("/teams/get/team", r.getTeam),
	}

	return r
}

// Routes returns all the API routes dedicated to the server system
func (s *teamsRouter) Routes() []router.Route {
	return s.Route
}
