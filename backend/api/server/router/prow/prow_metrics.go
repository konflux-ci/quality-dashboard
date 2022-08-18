package prow

import (
	"context"
	"fmt"
	"net/http"

	"github.com/redhat-appstudio/quality-studio/pkg/utils/httputils"
)

// version godoc
// @Summary Prow Jobs info
// @Description returns all prow jobs related to git_organization and repository_name
// @Tags Prow Jobs info
// @Accept json
// @Produce json
// @Param repository body GitRepositoryRequest true "repository name"
// @Param organization body GitRepositoryRequest true "git_organization"
// @Router /prow/results/get [get]
// @Success 200 {Object} []db.Prow
// @Failure 400 {object} types.ErrorResponse
func (s *jobRouter) getProwMetrics(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	metrics := s.Storage.GetMetrics()
	fmt.Println(metrics)
	return httputils.WriteJSON(w, http.StatusOK, metrics)
}
