package repositories

import (
	"github.com/redhat-appstudio/quality-studio/api/apis/codecov"
	"github.com/redhat-appstudio/quality-studio/api/apis/github"
	"github.com/redhat-appstudio/quality-studio/api/server/router"
	"github.com/redhat-appstudio/quality-studio/pkg/storage"
)

// systemRouter provides information about the server version.
type repositoryRouter struct {
	Route   []router.Route
	Storage storage.Storage
	Github  *github.Github
	CodeCov *codecov.API
}

// NewRouter initializes a new system router
func NewRouter(s storage.Storage) router.Router {
	githubAPI := github.NewGithubClient("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")
	codecovApi := codecov.NewCodeCoverageClient()
	r := &repositoryRouter{}

	r.Storage = s
	r.Github = githubAPI
	r.CodeCov = codecovApi

	r.Route = []router.Route{
		router.NewGetRoute("/repositories/list", r.listAllRepositoriesQuality),
		router.NewGetRoute("/workflows/get", r.getWorkflowByRepositoryName),
		router.NewPostRoute("/repositories/create", r.createRepositoryHandler),
		router.NewDeleteRoute("/repositories/delete", r.deleteRepositoryHandler),
	}

	return r
}

// Routes returns all the API routes dedicated to the server system
func (s *repositoryRouter) Routes() []router.Route {
	return s.Route
}
