package client

import (
	"context"
	"fmt"

	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/repository"
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

func (d *Database) DeleteTeam(teamName string) (bool, error) {
	team, err := d.client.Teams.Query().Where(teams.TeamName(teamName)).First(context.Background())

	if err != nil {
		return false, convertDBError("failed to return teams status: %w", err)
	}

	reposAssigned, err := d.client.Teams.QueryRepositories(team).All(context.Background())
	if err != nil {
		return false, fmt.Errorf("failed to get repositories assigned to team: %s, status: %v", teamName, err)
	}

	for _, repo := range reposAssigned {
		if _, err := d.client.Repository.Delete().Where(repository.RepositoryName(repo.RepositoryName)).Exec(context.TODO()); err != nil {
			return false, fmt.Errorf("failing to delete repository from database: %v", err)
		}
	}

	if _, err := d.client.Teams.Delete().Where(teams.TeamName(teamName)).Exec(context.TODO()); err != nil {
		return false, fmt.Errorf("failing to delete team from database: %v", err)
	}

	return true, nil
}

func (d *Database) UpdateTeam(t *db.Teams, target string) error {
	teamFromDb, err := d.GetTeamByName(target)
	if err != nil {
		return fmt.Errorf("failing to get team from database, team: %s, error %v", t.TeamName, err)
	}
	if _, err := d.client.Teams.UpdateOneID(teamFromDb.ID).SetDescription(t.Description).SetTeamName(t.TeamName).Save(context.TODO()); err != nil {
		return fmt.Errorf("failing to update team: %v", err)
	}

	return nil
}
