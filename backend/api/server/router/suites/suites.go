package suites

import (
	"context"
	"net/http"

	"github.com/redhat-appstudio/quality-studio/pkg/utils/httputils"
)

func (s *suitesRouter) getOcurrencies(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	freq, _ := s.Storage.GetSuitesFailureFrequency("redhat-appstudio", "infra-deployments", "pull-ci-redhat-appstudio-infra-deployments-main-appstudio-e2e-tests", "2023-10-10 00:00:00", "2023-11-15 23:59:59")
	s.Storage.GetFlakyTest("redhat-appstudio", "infra-deployments")
	return httputils.WriteJSON(w, http.StatusOK, freq)
}

func (s *suitesRouter) getTrends(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	freq := s.Storage.GetProwFlakyTrendsMetrics("redhat-appstudio", "infra-deployments", "pull-ci-redhat-appstudio-infra-deployments-main-appstudio-e2e-tests", "2023-11-10 00:00:00", "2023-11-15 23:59:59")

	return httputils.WriteJSON(w, http.StatusOK, freq)
}

func (s *suitesRouter) getFlakyJobs(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	freq := s.Storage.GetProwFlakyTrendsMetrics("redhat-appstudio", "infra-deployments", "pull-ci-redhat-appstudio-infra-deployments-main-appstudio-e2e-tests", "2023-11-10 00:00:00", "2023-11-15 23:59:59")

	return httputils.WriteJSON(w, http.StatusOK, freq)
}
