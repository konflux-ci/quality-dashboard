package client

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	prowV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/prow/v1alpha1"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/prowjobs"
)

func (d *Database) CreateProwJobResults(job prowV1Alpha1.Job, repo_id string) error {
	c, err := d.client.ProwJobs.Create().
		SetJobID(job.JobID).
		SetState(job.State).
		SetCreatedAt(job.CreatedAt).
		SetDuration(job.Duration).
		SetTestsCount(job.TestsCount).
		SetFailedCount(job.FailedCount).
		SetSkippedCount(job.SkippedCount).
		SetJobType(job.JobType).
		SetJobName(job.JobName).
		SetJobURL(job.JobURL).
		SetCiFailed(job.CIFailed).
		SetE2eFailedTestMessages(job.E2EFailedTestMessages).
		SetSuitesXMLURL(job.SuitesXmlUrl).
		SetBuildErrorLogs(job.BuildErrorLogs).
		Save(context.TODO())
	if err != nil {
		return convertDBError("create prow status: %w", err)
	}
	_, err = d.client.Repository.UpdateOneID(repo_id).AddProwJobs(c).Save(context.TODO())
	if err != nil {
		return convertDBError("create prow status: %w", err)
	}
	return nil
}

func (d *Database) GetLatestProwTestExecution(r *db.Repository, jobType string) (*db.ProwJobs, error) {
	jobId, err := d.client.Repository.QueryProwJobs(r).Where(prowjobs.JobType(jobType)).Where(prowjobs.State("success")).Order(db.Desc(prowjobs.FieldCreatedAt)).FirstID(context.Background())
	if err != nil {
		return nil, convertDBError("get prow job: %w", err)
	}

	prowJob, err := d.client.Repository.QueryProwJobs(r).Where(prowjobs.ID(jobId)).Order(db.Desc(prowjobs.FieldCreatedAt)).First(context.Background())

	if err != nil {
		return nil, convertDBError("get prow job: %w", err)
	}

	return prowJob, nil
}

func (d *Database) GetProwJobsResultsByJobID(jobID string) ([]*db.ProwJobs, error) {
	prowSuites, err := d.client.ProwJobs.Query().Where(prowjobs.JobID(jobID)).All(context.TODO())

	if err != nil {
		return nil, convertDBError("get prow suites: %w", err)
	}

	return prowSuites, nil
}

func (d *Database) UpdateBuildLogErrors(jobID, buildErrorLogs string) error {
	alreadyExists := d.client.ProwJobs.Query().
		Where(prowjobs.JobID(jobID)).
		ExistX(context.TODO())

	if !alreadyExists {
		return fmt.Errorf("jobID '%s' not found", jobID)
	}
	_, err := d.client.ProwJobs.Update().
		Where(prowjobs.JobID(jobID)).
		SetBuildErrorLogs(buildErrorLogs).
		Save(context.TODO())
	if err != nil {
		return convertDBError("failed to update failure: %w", err)
	}

	return nil
}

func (d *Database) GetAllProwJobs(startDate, endDate string) ([]*db.ProwJobs, error) {
	prowJobs, err := d.client.ProwJobs.Query().
		Where(func(s *sql.Selector) { // "created_at BETWEEN ? AND 2022-08-17", "2022-08-16"
			s.Where(sql.ExprP(fmt.Sprintf("created_at BETWEEN '%s' AND '%s'", startDate, endDate)))
		}).All(context.TODO())

	if err != nil {
		return nil, convertDBError("failed to get all prow jobs: %w", err)
	}

	return prowJobs, nil
}
