package database

import (
	"context"
	"net/http"

	"github.com/konflux-ci/quality-studio/api/types"
	"github.com/konflux-ci/quality-studio/pkg/utils/httputils"
)

func (d *databaseRouter) getDbConnection(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := d.db.Ping(); err != nil {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "database connection is down",
			StatusCode: 400,
		})
	}

	return httputils.WriteJSON(w, http.StatusOK, "database connection is up")
}
