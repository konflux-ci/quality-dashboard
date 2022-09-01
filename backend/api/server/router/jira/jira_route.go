package jira

import (
	jiraAPI "github.com/redhat-appstudio/quality-studio/api/apis/jira"
	"github.com/redhat-appstudio/quality-studio/api/server/router"
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
	}

	return r
}

// Routes returns all the API routes dedicated to the server system
func (s *jiraRouter) Routes() []router.Route {
	return s.Route
}
