package failure

import (
	"github.com/konflux-ci/quality-studio/api/server/router"
	"github.com/konflux-ci/quality-studio/pkg/logger"
	"github.com/konflux-ci/quality-studio/pkg/storage"
	"go.uber.org/zap"
)

// failureRouter provides information about the server version.
type failureRouter struct {
	Route   []router.Route
	Storage storage.Storage
	Logger  *zap.Logger
}

// NewRouter initializes a new system router
func NewRouter(s storage.Storage) router.Router {
	logger, _ := logger.InitZap("info")
	r := &failureRouter{}

	r.Storage = s
	r.Logger = logger

	r.Route = []router.Route{
		router.NewPostRoute("/failures/create", r.createFailure),
		router.NewGetRoute("/failures/get", r.getFailures),
		router.NewDeleteRoute("/failures/delete", r.deleteFailure),
	}

	return r
}

// Routes returns all the API routes dedicated to the server system
func (s *failureRouter) Routes() []router.Route {
	return s.Route
}
