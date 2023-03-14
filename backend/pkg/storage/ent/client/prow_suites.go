package client

import (
	"context"

	prowV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/prow/v1alpha1"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/prowsuites"
)

func (d *Database) CreateProwJobSuites(suites prowV1Alpha1.JobSuites, repo_id string) error {
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

func (d *Database) GetProwJobsResults(db *db.Repository) ([]*db.ProwJobs, error) {
	prowJobs, err := d.client.Repository.QueryProwJobs(db).Where().All(context.TODO())

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
