package client

import (
	"context"
	"strings"

	"github.com/google/uuid"
	repoV1Alpha1 "github.com/konflux-ci/quality-dashboard/api/apis/github/v1alpha1"
	"github.com/konflux-ci/quality-dashboard/pkg/storage"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db/repository"
)

// CreateRepository saves provided repository information in database.
func (d *Database) CreateRepository(repository repoV1Alpha1.Repository, team_id uuid.UUID) (*db.Repository, error) {
	repo, err := d.client.Repository.Create().
		SetRepositoryName(repository.Name).
		SetID(repository.ID).
		SetGitOrganization(repository.Owner.Login).
		SetDescription(repository.Description).
		SetGitURL(repository.URL).
		Save(context.TODO())
	if err != nil {
		return nil, convertDBError("create repository: %w", err)
	}

	_, err = d.client.Teams.UpdateOneID(team_id).AddRepositories(repo).Save(context.TODO())
	if err != nil {
		return nil, convertDBError("create repository: %w", err)
	}
	return repo, nil
}

// ListPasswords extracts an array of repositories from the database.
func (d *Database) ListRepositories(team *db.Teams) ([]repoV1Alpha1.Repository, error) {
	repositories, err := d.client.Teams.QueryRepositories(team).All(context.TODO())
	if err != nil {
		return nil, convertDBError("list repositories: %w", err)
	}

	storageRepositories := make([]repoV1Alpha1.Repository, 0, len(repositories))
	for _, p := range repositories {
		storageRepositories = append(storageRepositories, toStorageRepository(p))
	}
	return storageRepositories, nil
}

// GetRepository returns a git repo given its repository name and git organization name.
func (d *Database) GetRepository(repositoryName string, gitOrganizationName string) (*db.Repository, error) {
	repository, err := d.client.Repository.Query().
		Where(repository.RepositoryName(repositoryName)).Where(repository.GitOrganization(gitOrganizationName)).Only(context.TODO())

	if err != nil {
		return nil, convertDBError("get repository: %w", err)
	}

	return repository, nil
}

// ListRepositoriesQualityInfo extracts an array of repositories from the database.
func (d *Database) ListRepositoriesQualityInfo(team *db.Teams, startDate, endDate string) ([]storage.RepositoryQualityInfo, error) {
	repositories, err := d.client.Teams.QueryRepositories(team).All(context.TODO())
	if err != nil {
		return nil, convertDBError("list repositories: %w", err)
	}

	storageRepositories := make([]storage.RepositoryQualityInfo, 0, len(repositories))
	for _, p := range repositories {
		c, err := d.client.Repository.QueryCodecov(p).Only(context.TODO())
		if err != nil {
			return nil, convertDBError("list coverage failed: %w", err)
		}

		var prs repoV1Alpha1.PullRequestsInfo
		var workflows []repoV1Alpha1.Workflow

		prs, err = d.GetPullRequestsByRepository(p.RepositoryName, p.GitOrganization, startDate, endDate)
		if err != nil {
			return nil, convertDBError("list pull requests failed: %w", err)
		}

		workflows, err = d.ListWorkflowsByRepository(p.RepositoryName)
		if err != nil {
			return nil, convertDBError("list workflows failed: %w", err)
		}

		storageRepositories = append(storageRepositories, toStorageRepositoryAllInfo(p, c, prs, workflows))
	}
	return storageRepositories, nil
}

// DeleteRepository deletes a repository from the database by repository name and git organization name.
func (d *Database) DeleteRepository(repositoryName string, gitOrganizationName string) error {
	repositoryName = strings.ToLower(repositoryName)
	_, err := d.client.Repository.Delete().
		Where(repository.RepositoryName(repositoryName)).Where(repository.GitOrganization(gitOrganizationName)).
		Exec(context.TODO())

	if err != nil {
		return convertDBError("delete repository: %w", err)
	}
	return nil
}

// ListAllRepositories extracts an array of repositories from the database.
func (d *Database) ListAllRepositories() ([]*db.Repository, error) {
	repositories, err := d.client.Repository.Query().All(context.Background())
	if err != nil {
		return nil, convertDBError("list repositories: %w", err)
	}

	return repositories, nil
}
