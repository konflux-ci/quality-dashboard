package client

import (
	"context"

	"github.com/google/uuid"
	"github.com/redhat-appstudio/quality-studio/pkg/storage"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/prowjobs"
)

func (d *Database) CreateProwJobResults(prowJobStatus storage.ProwJobStatus, repo_id uuid.UUID) error {
	c, err := d.client.ProwJobs.Create().
		SetJobID(prowJobStatus.JobID).
		SetState(prowJobStatus.State).
		SetCreatedAt(prowJobStatus.CreatedAt).
		SetDuration(prowJobStatus.Duration).
		SetTestsCount(prowJobStatus.TestsCount).
		SetFailedCount(prowJobStatus.FailedCount).
		SetSkippedCount(prowJobStatus.SkippedCount).
		SetJobType(prowJobStatus.JobType).
		SetJobName(prowJobStatus.JobName).
		SetJobURL(prowJobStatus.JobURL).
		SetCiFailed(prowJobStatus.CIFailed).
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

// GetRepository returns a git repo given its url
func (d *Database) GetProwJobsResultsByJobID(jobID string) ([]*db.ProwJobs, error) {
	prowSuites, err := d.client.ProwJobs.Query().Where(prowjobs.JobID(jobID)).All(context.TODO())

	if err != nil {
		return nil, convertDBError("get prow suites: %w", err)
	}

	return prowSuites, nil
}
