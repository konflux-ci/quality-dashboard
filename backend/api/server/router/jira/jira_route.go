package jira

import (
	"github.com/redhat-appstudio/quality-studio/api/server/router"
	jiraAPI "github.com/redhat-appstudio/quality-studio/pkg/connectors/jira"
	"github.com/redhat-appstudio/quality-studio/pkg/logger"
	"github.com/redhat-appstudio/quality-studio/pkg/storage"
	"go.uber.org/zap"
)

// systemRouter provides information about the server version.
type jiraRouter struct {
	Route   []router.Route
	Storage storage.Storage
	Jira    jiraAPI.Jira
	Logger  *zap.Logger
}

// NewRouter initializes a new system router
func NewRouter(s storage.Storage) router.Router {
	jiraAPI := jiraAPI.NewJiraConfig()
	logger, _ := logger.InitZap("info")
	r := &jiraRouter{}

	r.Storage = s
	r.Jira = jiraAPI
	r.Logger = logger

	r.Route = []router.Route{
		router.NewGetRoute("/jira/bugs/e2e", r.listE2EBugsKnown),
		router.NewGetRoute("/jira/bugs/all", r.listAllBugs), ///jira/project/list
		router.NewGetRoute("/jira/bugs/metrics/priorities", r.getCountBugsForAllCategories),
		router.NewGetRoute("/jira/project/list", r.getJiraProjects),
		router.NewPostRoute("/jira/bugs/metrics/resolution", r.calculateRates),
		router.NewPostRoute("/jira/bugs/metrics/open", r.openBugsMetrics),
		router.NewGetRoute("/jira/bugs/exist", r.bugExists),
		router.NewGetRoute("/jira/slos/list", r.getBugSLOs),
	}

	return r
}

// Routes returns all the API routes dedicated to the server system
func (s *jiraRouter) Routes() []router.Route {
	return s.Route
}
