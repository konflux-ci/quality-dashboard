package client

import (
	"fmt"

	"github.com/flacatus/qe-dashboard-backend/pkg/storage"
	"github.com/flacatus/qe-dashboard-backend/pkg/storage/ent/db"
)

func convertDBError(t string, err error) error {
	if db.IsNotFound(err) {
		return storage.ErrNotFound
	}

	if db.IsConstraintError(err) {
		return storage.ErrAlreadyExists
	}

	return fmt.Errorf(t, err)
}
