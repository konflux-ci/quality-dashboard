package client

import (
	"context"
	"fmt"
	"time"

	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db/oci"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db/repository"
)

// NotFoundError is returned when a requested entity is not found in the database.
type NotFoundError struct {
	label string
}

// Error returns the error message string.
func (e *NotFoundError) Error() string {
	return e.label
}

// CreateOCIArtifact creates a new OCI artifact record.
// Based on the provided schema, it can optionally link the artifact to a repository.
// To link to a repository, provide a non-nil repoID.
func (d *Database) CreateOCIArtifact(ociArtifact *db.OCI, repoID *string) (*db.OCI, error) {
	// Start building the creation query for the OCI artifact.
	// We only set fields that exist in the new schema.
	creator := d.client.OCI.Create().
		SetArtifactURL(ociArtifact.ArtifactURL)

	// Conditionally link to a repository if a repoID is provided.
	// The schema indicates the relationship is optional.
	// The method name is SetOciID because the edge is named "oci".
	if repoID != nil {
		creator.SetOciID(*repoID)
	}

	// Execute the create operation. This is atomic.
	ociNode, err := creator.Save(context.Background())
	if err != nil {
		return nil, convertDBError("create oci artifact: %w", err)
	}

	return ociNode, nil
}

// GetAllOCIArtifacts retrieves all OCI artifacts from the database.
// For performance, it eagerly loads the associated Repository to prevent N+1 query problems
// when accessing the parent repository of an artifact.
func (d *Database) GetAllOCIArtifacts() ([]*db.OCI, error) {
	// The method is WithOci() because the edge in the OCI schema is named "oci".
	ociArtifacts, err := d.client.OCI.Query().
		WithOci(). // Eagerly load the parent repository via the "oci" edge.
		All(context.Background())
	if err != nil {
		return nil, convertDBError("failed to get all oci artifacts: %w", err)
	}
	return ociArtifacts, nil
}

// GetOCIArtifactsByRepository retrieves all OCI artifacts associated with a given repository ID.
// This is a more practical and performant query than fetching all artifacts in the database.
func (d *Database) GetOCIArtifactsByRepository(repoID string) ([]*db.OCI, error) {
	exists, err := d.client.Repository.Query().
		Where(repository.ID(repoID)).
		Exist(context.Background())
	if err != nil {
		return nil, convertDBError(fmt.Sprintf("failed to check for repository '%s'", repoID), err)
	}
	if !exists {
		return nil, &NotFoundError{label: fmt.Sprintf("repository with id '%s' not found", repoID)}
	}

	// Query for the OCI artifacts by filtering on the "oci" edge.
	// This is the correct way to query for nodes that have a relationship with another node.
	ociArtifacts, err := d.client.OCI.Query().
		Where(oci.HasOciWith(repository.ID(repoID))).
		WithOci().
		All(context.Background())
	if err != nil {
		return nil, convertDBError(fmt.Sprintf("failed to get oci artifacts for repository '%s'", repoID), err)
	}

	return ociArtifacts, nil
}

// UpdateOCIArtifact updates an existing OCI artifact record identified by its ID.
func (d *Database) UpdateOCIArtifact(ctx context.Context, ociArtifact *db.OCI) (*db.OCI, error) {
	updater := d.client.OCI.UpdateOneID(ociArtifact.ID).
		SetArtifactURL(ociArtifact.ArtifactURL).
		SetUpdatedAt(time.Now())

	updatedNode, err := updater.Save(ctx)
	if err != nil {
		return nil, convertDBError("failed to update artifact", err)
	}
	return updatedNode, nil
}
