package client

import (
	"context"

	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/teams"
)

func (d *Database) CreateQualityStudioTeam(teamName string, description string) (*db.Teams, error) {
	team, err := d.client.Teams.Create().
		SetTeamName(teamName).
		SetDescription(description).
		Save(context.TODO())
	if err != nil {
		return nil, convertDBError("create team status: %w", err)
	}

	return team, nil
}

func (d *Database) GetAllTeamsFromDB() ([]*db.Teams, error) {
	teams, err := d.client.Teams.Query().All(context.Background())

	if err != nil {
		return nil, convertDBError("failed to return teams status: %w", err)
	}

	return teams, nil
}

func (d *Database) GetTeamByName(teamName string) (*db.Teams, error) {
	teams, err := d.client.Teams.Query().Where(teams.TeamName(teamName)).First(context.Background())

	if err != nil {
		return nil, convertDBError("failed to return teams status: %w", err)
	}

	return teams, nil
}
