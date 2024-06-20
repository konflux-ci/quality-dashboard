package suites

import (
	"github.com/konflux-ci/quality-dashboard/api/server/router"
	"github.com/konflux-ci/quality-dashboard/pkg/logger"
	"github.com/konflux-ci/quality-dashboard/pkg/storage"
	"go.uber.org/zap"
)

// failureRouter provides information about the server version.
type suitesRouter struct {
	Route   []router.Route
	Storage storage.Storage
	Logger  *zap.Logger
}

// NewRouter initializes a new system router
func NewRouter(s storage.Storage) router.Router {
	logger, _ := logger.InitZap("info")
	r := &suitesRouter{}

	r.Storage = s
	r.Logger = logger

	r.Route = []router.Route{
		router.NewGetRoute("/suites/occurrences", r.getOccurrences),
		router.NewGetRoute("/suites/flaky/trends", r.getTrends),
	}

	return r
}

// Routes returns all the API routes dedicated to the server system
func (s *suitesRouter) Routes() []router.Route {
	return s.Route
}
