package suites

import (
	"context"
	"net/http"

	"github.com/redhat-appstudio/quality-studio/pkg/utils/httputils"
)

func (s *suitesRouter) getOcurrencies(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	freq, _ := s.Storage.GetSuitesFailureFrequency("redhat-appstudio", "infra-deployments", "", "")

	return httputils.WriteJSON(w, http.StatusOK, freq)
}
