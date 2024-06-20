package server

import "github.com/konflux-ci/quality-dashboard/pkg/utils/httputils"

// handlerWithGlobalMiddlewares wraps the handler function for a request with
// the server's global middlewares. The order of the middlewares is backwards,
// meaning that the first in the list will be evaluated last.
func (s *Server) handlerWithGlobalMiddlewares(handler httputils.APIFunc) httputils.APIFunc {
	next := handler
	for _, m := range s.middlewares {
		next = m.WrapHandler(next)
	}

	return next
}
