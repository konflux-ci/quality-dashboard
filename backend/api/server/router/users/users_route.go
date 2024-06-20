package users

import (
	"github.com/konflux-ci/quality-studio/api/server/router"
	jiraAPI "github.com/konflux-ci/quality-studio/pkg/connectors/jira"
	"github.com/konflux-ci/quality-studio/pkg/logger"
	"github.com/konflux-ci/quality-studio/pkg/storage"
	"go.uber.org/zap"
)

// systemRouter provides information about the server version.
type usersRouter struct {
	Route   []router.Route
	Storage storage.Storage
	Logger  *zap.Logger
	Jira    jiraAPI.Jira
}

// NewRouter initializes a new system router
func NewRouter(s storage.Storage) router.Router {
	jiraAPI := jiraAPI.NewJiraConfig()
	logger, _ := logger.InitZap("info")
	r := &usersRouter{}

	r.Storage = s
	r.Logger = logger
	r.Jira = jiraAPI

	r.Route = []router.Route{
		router.NewPostRoute("/users/create", r.createUser),
		router.NewGetRoute("/users/get/all", r.listAllUsers),
		router.NewGetRoute("/users/get/user", r.getUser),
	}

	return r
}

// Routes returns all the API routes dedicated to the server system
func (s *usersRouter) Routes() []router.Route {
	return s.Route
}
