package client

import (
	"context"

	"github.com/google/uuid"
	"github.com/redhat-appstudio/quality-studio/pkg/storage"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/prowsuites"
)

// CreateProwJobResults save provided repository information in database.
func (d *Database) CreateProwJobSuites(suites storage.ProwJobSuites, repo_id uuid.UUID) error {
	c, err := d.client.ProwSuites.Create().
		SetJobID(suites.JobID).
		SetName(suites.TestCaseName).
		SetStatus(suites.TestCaseStatus).
		SetTime(suites.TestTiming).
		Save(context.TODO())
	if err != nil {
		return convertDBError("create prow: %w", err)
	}
	_, err = d.client.Repository.UpdateOneID(repo_id).AddProwSuites(c).Save(context.TODO())
	if err != nil {
		return convertDBError("create prow: %w", err)
	}
	return nil
}

func (d *Database) GetProwJobsResults(db *db.Repository) ([]*db.ProwSuites, error) {
	prowJobs, err := d.client.Repository.QueryProwSuites(db).Where().All(context.TODO())

	if err != nil {
		return nil, convertDBError("get repository: %w", err)
	}

	return prowJobs, nil
}

func (d *Database) GetSuitesByJobID(jobID string) ([]*db.ProwSuites, error) {
	prowSuites, err := d.client.ProwSuites.Query().Where(prowsuites.JobID(jobID)).All(context.Background())

	if err != nil {
		return nil, convertDBError("get prow job: %w", err)
	}

	return prowSuites, nil
}
