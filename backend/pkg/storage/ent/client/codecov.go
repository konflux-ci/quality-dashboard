package client

import (
	"context"

	coverageV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/codecov/v1alpha1"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/codecov"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/repository"
)

// CreateCoverage saves provided coverage information in database.
func (d *Database) CreateCoverage(codecoverage coverageV1Alpha1.Coverage, repo_id string) error {
	c, err := d.client.CodeCov.Create().
		SetRepositoryName(codecoverage.RepositoryName).
		SetGitOrganization(codecoverage.GitOrganization).
		SetCoveragePercentage(codecoverage.CoveragePercentage).
		SetAverageRetestsToMerge(codecoverage.AverageToRetestPullRequest).
		Save(context.TODO())
	if err != nil {
		return convertDBError("create coverage: %w", err)
	}
	_, err = d.client.Repository.UpdateOneID(repo_id).AddCodecov(c).Save(context.TODO())
	if err != nil {
		return convertDBError("create coverage: %w", err)
	}
	return nil
}

func (d *Database) UpdateCoverage(codecoverage coverageV1Alpha1.Coverage, repoName string) error {
	_, err := d.client.CodeCov.Delete().Where(codecov.RepositoryName(repoName)).Exec(context.TODO())
	if err != nil {
		return convertDBError("delete coverage from database: %w", err)
	}

	c, err := d.client.CodeCov.Create().
		SetRepositoryName(codecoverage.RepositoryName).
		SetGitOrganization(codecoverage.GitOrganization).
		SetCoveragePercentage(codecoverage.CoveragePercentage).
		SetAverageRetestsToMerge(codecoverage.AverageToRetestPullRequest).
		Save(context.TODO())
	if err != nil {
		return convertDBError("create coverage: %w", err)
	}
	_, err = d.client.Repository.Update().Where(repository.RepositoryName(repoName)).AddCodecov(c).Save(context.TODO())
	if err != nil {
		return convertDBError("create coverage: %w", err)
	}
	return nil
}
