package client

import (
	"context"

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
