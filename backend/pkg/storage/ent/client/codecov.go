package client

import (
	"context"

	"github.com/flacatus/qe-dashboard-backend/pkg/storage"
	"github.com/flacatus/qe-dashboard-backend/pkg/storage/ent/db/codecov"
	"github.com/flacatus/qe-dashboard-backend/pkg/storage/ent/db/repository"
	"github.com/google/uuid"
)

// CreateRepository save provided repository information in database.
func (d *Database) CreateCoverage(repository storage.Coverage, repo_id uuid.UUID) error {
	c, err := d.client.CodeCov.Create().
		SetRepositoryName(repository.RepositoryName).
		SetGitOrganization(repository.GitOrganization).
		SetCoveragePercentage(repository.CoveragePercentage).
		Save(context.TODO())
	if err != nil {
		return convertDBError("create coverage: %w", err)
	}
	_, err = d.client.Repository.UpdateOneID(repo_id).AddCodecov(c).Save(context.TODO())
	if err != nil {
		return convertDBError("create workflows: %w", err)
	}
	return nil
}

func (d *Database) UpdateCoverage(codecoverage storage.Coverage, repoName string) error {
	_, err := d.client.CodeCov.Delete().Where(codecov.RepositoryName(repoName)).Exec(context.TODO())
	if err != nil {
		return convertDBError("delete workflow from database: %w", err)
	}

	c, err := d.client.CodeCov.Create().
		SetRepositoryName(codecoverage.RepositoryName).
		SetGitOrganization(codecoverage.GitOrganization).
		SetCoveragePercentage(codecoverage.CoveragePercentage).
		Save(context.TODO())
	if err != nil {
		return convertDBError("create workflows: %w", err)
	}
	_, err = d.client.Repository.Update().Where(repository.RepositoryName(repoName)).AddCodecov(c).Save(context.TODO())
	if err != nil {
		return convertDBError("create workflows: %w", err)
	}
	return nil
}
