package client

import (
	"context"

	repoV1Alpha1 "github.com/konflux-ci/quality-dashboard/api/apis/github/v1alpha1"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db/repository"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db/workflows"
)

// CreateWorkflows saves provided workflow information in database.
func (d *Database) CreateWorkflows(workflow repoV1Alpha1.Workflow, repo_id string) error {
	w, err := d.client.Workflows.Create().
		SetWorkflowName(workflow.Name).
		SetBadgeURL(workflow.BadgeURL).
		SetHTMLURL(workflow.HTMLURL).
		SetJobURL(workflow.JobURL).
		SetState(workflow.State).
		Save(context.TODO())
	if err != nil {
		return convertDBError("create workflows: %w", err)
	}
	_, err = d.client.Repository.UpdateOneID(repo_id).AddWorkflows(w).Save(context.TODO())
	if err != nil {
		return convertDBError("create workflows: %w", err)
	}
	return nil
}

func (d *Database) ReCreateWorkflow(workflow repoV1Alpha1.Workflow, repoName string) error {
	_, err := d.client.Workflows.Delete().Where(workflows.WorkflowName(workflow.Name)).Exec(context.TODO())
	if err != nil {
		return convertDBError("delete workflow from database: %w", err)
	}

	w, err := d.client.Workflows.Create().
		SetWorkflowName(workflow.Name).
		SetBadgeURL(workflow.BadgeURL).
		SetHTMLURL(workflow.HTMLURL).
		SetJobURL(workflow.JobURL).
		SetState(workflow.State).
		Save(context.TODO())
	if err != nil {
		return convertDBError("create workflows: %w", err)
	}
	_, err = d.client.Repository.Update().Where(repository.RepositoryName(repoName)).AddWorkflows(w).Save(context.TODO())
	if err != nil {
		return convertDBError("create workflows: %w", err)
	}
	return nil
}

func (d *Database) ListWorkflowsByRepository(repositoryName string) (w []repoV1Alpha1.Workflow, err error) {
	repositories, err := d.client.Repository.Query().All(context.TODO())
	if err != nil {
		return nil, convertDBError("list repositories: %w", err)
	}

	storageWorkflows := make([]repoV1Alpha1.Workflow, 0, len(repositories))
	for _, p := range repositories {
		if p.RepositoryName == repositoryName {
			w, _ := d.client.Repository.QueryWorkflows(p).All(context.TODO())
			for _, workflow := range w {
				storageWorkflows = append(storageWorkflows, toStorageWorkflows(workflow))
			}
		}
	}
	return storageWorkflows, err
}
