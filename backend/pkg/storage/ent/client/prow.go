package client

import (
	"context"

	"github.com/google/uuid"
	"github.com/redhat-appstudio/quality-studio/pkg/storage"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/prow"
)

// CreateProwJobResults save provided repository information in database.
func (d *Database) CreateProwJobResults(repository storage.ProwJob, repo_id uuid.UUID) error {
	c, err := d.client.Prow.Create().
		SetJobID(repository.JobID).
		SetName(repository.TestCaseName).
		SetStatus(repository.TestCaseStatus).
		SetTime(repository.TestTiming).
		Save(context.TODO())
	if err != nil {
		return convertDBError("create prow: %w", err)
	}
	_, err = d.client.Repository.UpdateOneID(repo_id).AddProw(c).Save(context.TODO())
	if err != nil {
		return convertDBError("create prow: %w", err)
	}
	return nil
}

// GetRepository returns a git repo given its url
func (d *Database) GetProwJobsResults(db *db.Repository) ([]*db.Prow, error) {
	prowJobs, err := d.client.Repository.QueryProw(db).Where().All(context.TODO())

	if err != nil {
		return nil, convertDBError("get repository: %w", err)
	}

	return prowJobs, nil
}

// GetRepository returns a git repo given its url
func (d *Database) GetProwJobsResultsByJobID(jobID string) ([]*db.Prow, error) {
	prowJobs, err := d.client.Prow.Query().Where(prow.JobID(jobID)).All(context.TODO())

	if err != nil {
		return nil, convertDBError("get repository: %w", err)
	}

	return prowJobs, nil
}
