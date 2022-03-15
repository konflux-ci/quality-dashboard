package api

import (
	"net/http"

	"github.com/flacatus/qe-dashboard-backend/pkg/version"
)

// Version godoc
// @Summary Version
// @Description returns quality backend version
// @Tags Version API
// @Produce json
// @Router /api/version [get]
// @Success 200 {object} api.MapResponse
func (s *Server) versionHandler(w http.ResponseWriter, r *http.Request) {
	result := map[string]string{
		"version": version.VERSION,
	}
	s.JSONResponse(w, r, result)
}
