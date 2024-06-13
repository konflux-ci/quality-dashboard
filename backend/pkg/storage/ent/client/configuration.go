package client

import (
	"context"

	"github.com/redhat-appstudio/quality-studio/api/apis/configuration/v1alpha1"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/configuration"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/predicate"
)

func (d *Database) CreateConfiguration(c v1alpha1.Configuration) error {
	configAlreadyExists := d.client.Configuration.Query().
		Where(configuration.TeamName(c.TeamName)).
		ExistX(context.TODO())
	if configAlreadyExists {
		_, err := d.client.Configuration.Update().
			Where(predicate.Configuration(configuration.TeamName(c.TeamName))).
			SetJiraConfig(c.JiraConfig).
			SetBugSlosConfig(c.BugSLOsConfig).
			Save(context.TODO())
		if err != nil {
			return convertDBError("failed to update configuration: %w", err)
		}
	} else {
		_, err := d.client.Configuration.Create().
			SetTeamName(c.TeamName).
			SetJiraConfig(c.JiraConfig).
			SetBugSlosConfig(c.BugSLOsConfig).
			Save(context.TODO())
		if err != nil {
			return convertDBError("failed to create configuration: %w", err)
		}
	}

	return nil
}

func (d *Database) GetConfiguration(teamName string) (*db.Configuration, error) {
	configuration, err := d.client.Configuration.Query().
		Where(configuration.TeamName(teamName)).Only(context.TODO())
	if err != nil {
		return nil, convertDBError("failed to get configuration: %w", err)
	}

	return configuration, nil
}
