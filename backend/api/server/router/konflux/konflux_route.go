package konflux

import (
	"github.com/konflux-ci/quality-dashboard/api/server/router"
	"github.com/konflux-ci/quality-dashboard/pkg/logger"
	"github.com/konflux-ci/quality-dashboard/pkg/storage"
	"go.uber.org/zap"
)

type konfluxRouter struct {
	routes  []router.Route
	storage storage.Storage
	logger  *zap.Logger
}

func (k *konfluxRouter) Routes() []router.Route {
	return k.routes
}

func NewRouter(storage storage.Storage) router.Router {
	kloggger, _ := logger.InitZap("info")
	k := &konfluxRouter{}
	k.routes = []router.Route{
		router.NewPostRoute("/konflux/metadata/post", k.receiveMetrics),
	}
	k.storage = storage
	k.logger = kloggger
	return k
}
