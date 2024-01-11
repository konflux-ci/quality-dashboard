package prow

import (
	"github.com/redhat-appstudio/quality-studio/api/server/router"
	"github.com/redhat-appstudio/quality-studio/pkg/connectors/github"
	"github.com/redhat-appstudio/quality-studio/pkg/logger"
	"github.com/redhat-appstudio/quality-studio/pkg/storage"
	util "github.com/redhat-appstudio/quality-studio/pkg/utils"
	"go.uber.org/zap"
)

// systemRouter provides information about the server version.
type jobRouter struct {
	routes  []router.Route
	Storage storage.Storage
	Logger  zap.Logger
	Github  *github.Github
}

// NewRouter initializes a new system router
func NewRouter(storage storage.Storage) router.Router {
	githubAPI := github.NewGithubClient(util.GetEnv("GITHUB_TOKEN", ""))
	r := &jobRouter{}
	r.Github = githubAPI
	logger, _ := logger.InitZap("info")
	r.Logger = *logger
	r.Storage = storage
	r.routes = []router.Route{ // getProwJobsByType
		router.NewGetRoute("/prow/results/get", r.getProwJobs),
		router.NewGetRoute("/prow/jobs/types", r.getProwJobsByType),
		router.NewGetRoute("/prow/results/latest/get", r.getLatestSuitesExecution),
		router.NewGetRoute("/prow/metrics/get", r.getProwMetrics),
		router.NewGetRoute("/prow/repositories/list", r.listProwRepos),
		router.NewGetRoute("/prow/metrics/daily", r.getProwMetricsByDay),
	}

	return r
}

// Routes returns all the API routes dedicated to the server system
func (s *jobRouter) Routes() []router.Route {
	return s.routes
}
