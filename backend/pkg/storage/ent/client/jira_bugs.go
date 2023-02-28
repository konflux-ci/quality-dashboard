package client

import (
	"context"

	jiraV1Alpha "github.com/redhat-appstudio/quality-studio/api/apis/jira/v1alpha1"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/bugs"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/predicate"
)

// CreateRepository save provided repository information in database.
func (d *Database) CreateJiraBug(jiraBug jiraV1Alpha.JiraBug) error {
	bugAlreadyExists := d.client.Bugs.Query().Where(bugs.JiraKey(jiraBug.JiraKey)).ExistX(context.TODO())

	if bugAlreadyExists {
		_, err := d.client.Bugs.Update().Where(predicate.Bugs(bugs.JiraKey(jiraBug.JiraKey))).
			SetCreatedAt(jiraBug.CreatedAt).
			SetUpdatedAt(jiraBug.UpdatedAt).
			SetJiraKey(jiraBug.JiraKey).
			SetPriority(jiraBug.Priority).
			SetSummary(jiraBug.Summary).
			SetURL(jiraBug.Url).
			SetStatus(jiraBug.Status).
			Save(context.TODO())
		if err != nil {
			return convertDBError("failed to create bug: %w", err)
		}
	} else {
		_, err := d.client.Bugs.Create().
			SetCreatedAt(jiraBug.CreatedAt).
			SetUpdatedAt(jiraBug.UpdatedAt).
			SetJiraKey(jiraBug.JiraKey).
			SetPriority(jiraBug.Priority).
			SetSummary(jiraBug.Summary).
			SetURL(jiraBug.Url).
			SetStatus(jiraBug.Status).
			Save(context.TODO())
		if err != nil {
			return convertDBError("failed to update bug: %w", err)
		}
	}

	return nil
}

func (d *Database) GetAllJiraBugs() ([]*db.Bugs, error) {
	bugs, err := d.client.Bugs.Query().All(context.Background())

	if err != nil {
		return nil, convertDBError("failed to return bugs: %w", err)
	}

	return bugs, nil
}
