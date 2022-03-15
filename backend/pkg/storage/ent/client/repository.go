package client

import (
	"context"
	"strings"

	"github.com/flacatus/qe-dashboard-backend/pkg/storage"
	"github.com/flacatus/qe-dashboard-backend/pkg/storage/ent/db"
	"github.com/flacatus/qe-dashboard-backend/pkg/storage/ent/db/repository"
)

// CreateRepository save provided repository information in database.
func (d *Database) CreateRepository(repository storage.Repository) (*db.Repository, error) {
	repo, err := d.client.Repository.Create().
		SetRepositoryName(repository.RepositoryName).
		SetGitOrganization(repository.GitOrganization).
		SetDescription(repository.Description).
		SetGitURL(repository.GitURL).
		Save(context.TODO())
	if err != nil {
		return nil, convertDBError("create repository: %w", err)
	}
	return repo, nil
}

// ListPasswords extracts an array of repositories from the database.
func (d *Database) ListRepositories() ([]storage.Repository, error) {
	repositories, err := d.client.Repository.Query().All(context.TODO())
	if err != nil {
		return nil, convertDBError("list repositories: %w", err)
	}

	storageRepositories := make([]storage.Repository, 0, len(repositories))
	for _, p := range repositories {
		storageRepositories = append(storageRepositories, toStorageRepository(p))
	}
	return storageRepositories, nil
}

// ListRepositoriesQualityInfo extracts an array of repositories from the database.
func (d *Database) ListRepositoriesQualityInfo() ([]storage.RepositoryQualityInfo, error) {
	repositories, err := d.client.Repository.Query().All(context.TODO())
	if err != nil {
		return nil, convertDBError("list repositories: %w", err)
	}

	storageRepositories := make([]storage.RepositoryQualityInfo, 0, len(repositories))
	for _, p := range repositories {
		w, _ := d.client.Repository.QueryWorkflows(p).All(context.TODO())
		c, _ := d.client.Repository.QueryCodecov(p).Only(context.TODO())
		storageRepositories = append(storageRepositories, toStorageRepositoryAllInfo(p, w, c))
	}
	return storageRepositories, nil
}

// DeletePassword deletes a password from the database by email.
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
