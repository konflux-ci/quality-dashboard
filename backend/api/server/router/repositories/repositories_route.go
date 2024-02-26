package repositories

import (
	"github.com/redhat-appstudio/quality-studio/api/server/router"
	"github.com/redhat-appstudio/quality-studio/pkg/connectors/codecov"
	"github.com/redhat-appstudio/quality-studio/pkg/connectors/github"
	"github.com/redhat-appstudio/quality-studio/pkg/logger"
	"github.com/redhat-appstudio/quality-studio/pkg/storage"
	util "github.com/redhat-appstudio/quality-studio/pkg/utils"
	"go.uber.org/zap"
)

// systemRouter provides information about the server version.
type repositoryRouter struct {
	Route   []router.Route
	Storage storage.Storage
	Github  *github.Github
	CodeCov *codecov.API
	Logger  *zap.Logger
}

// NewRouter initializes a new system router
func NewRouter(s storage.Storage) router.Router {
	githubAPI := github.NewGithubClient(util.GetEnv("GITHUB_TOKEN", ""))

	logger, _ := logger.InitZap("info")

	codecovApi := codecov.NewCodeCoverageClient()
	r := &repositoryRouter{}

	r.Storage = s
	r.Github = githubAPI
	r.CodeCov = codecovApi
	r.Logger = logger

	r.Route = []router.Route{
		router.NewGetRoute("/repositories/list", r.listAllRepositoriesQuality),
		router.NewGetRoute("/workflows/get", r.getWorkflowByRepositoryName),
		router.NewGetRoute("/repositories/getJobTypesFromRepo", r.getJobTypesFromRepo),
		router.NewPostRoute("/repositories/create", r.createRepositoryHandler),
		router.NewDeleteRoute("/repositories/delete", r.deleteRepositoryHandler),
		router.NewGetRoute("/prs/get", r.getPullRequestsFromRepo),
		router.NewGetRoute("/repositories/verify", r.checkGithubRepository),
	}

	return r
}

// Routes returns all the API routes dedicated to the server system
func (s *repositoryRouter) Routes() []router.Route {
	return s.Route
}
