package client

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	prowV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/prow/v1alpha1"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/prowsuites"
)

func (d *Database) CreateProwJobSuites(suites prowV1Alpha1.JobSuites, repo_id string) error {
	alreadyExists := d.client.ProwSuites.Query().
		Where(prowsuites.JobID(suites.JobID)).
		Where(prowsuites.SuiteName(suites.SuiteName)).
		Where(prowsuites.Name(suites.TestCaseName)).
		ExistX(context.TODO())

	if !alreadyExists {
		c, err := d.client.ProwSuites.Create().
			SetJobID(suites.JobID).
			SetJobURL(suites.JobURL).
			SetJobName(suites.JobName).
			SetSuiteName(suites.SuiteName).
			SetName(suites.TestCaseName).
			SetStatus(suites.TestCaseStatus).
			SetTime(suites.TestTiming).
			SetErrorMessage(suites.ErrorMessage).
			SetCreatedAt(suites.CreatedAt).
			Save(context.TODO())
		if err != nil {
			return convertDBError("create prow: %w", err)
		}
		_, err = d.client.Repository.UpdateOneID(repo_id).AddProwSuites(c).Save(context.TODO())
		if err != nil {
			return convertDBError("create prow: %w", err)
		}
	}

	return nil
}

func (d *Database) GetProwJobsResults(repo *db.Repository, startDate, endDate string) ([]*db.ProwJobs, error) {
	prowJobs, err := d.client.Repository.QueryProwJobs(repo).
		Where(func(s *sql.Selector) { // "created_at BETWEEN ? AND 2022-08-17", "2022-08-16"
			s.Where(sql.ExprP(fmt.Sprintf("created_at BETWEEN '%s' AND '%s'", startDate, endDate)))
		}).All(context.TODO())

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
