package konflux

import (
	"github.com/redhat-appstudio/quality-studio/api/server/router"
	"github.com/redhat-appstudio/quality-studio/pkg/logger"
	"github.com/redhat-appstudio/quality-studio/pkg/storage"
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
