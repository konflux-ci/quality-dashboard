package client

import (
	"context"

	"github.com/google/uuid"
	"github.com/redhat-appstudio/quality-studio/pkg/storage"
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
