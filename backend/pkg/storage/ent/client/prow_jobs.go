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
		SetCreatedAt(prowJobStatus.CreatedAt).
		SetDuration(prowJobStatus.Duration).
		SetTestsCount(prowJobStatus.TestsCount).
		SetFailedCount(prowJobStatus.FailedCount).
		SetSkippedCount(prowJobStatus.SkippedCount).
		SetJobType(prowJobStatus.JobType).
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
	prowJob, err := d.client.Repository.QueryProwJobs(r).Where(prowjobs.JobType(jobType)).Order(db.Desc(prowjobs.FieldCreatedAt)).First(context.Background())
	if err != nil {
		return nil, convertDBError("get prow job: %w", err)
	}

	return prowJob, nil
}
