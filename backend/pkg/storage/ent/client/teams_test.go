package client

import (
	"testing"

	"github.com/devfile/library/pkg/util"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateQualityStudioTeam(t *testing.T) {
	// get database client
	cfg := GetPostgresConnectionDetails()
	storage, err := cfg.Open()
	require.NoError(t, err, "GetPostgresConnectionDetails() should not return an error")

	teamName := "team-" + util.GenerateRandomString(6)
	teamDescription := teamName

	type Input struct {
		TeamName        string
		TeamDescription string
	}

	cases := []struct {
		Name          string
		Input         *Input
		Expected      *db.Teams
		ExpectedError string
	}{
		{
			Name: "create a quality studio team successfully",
			Input: &Input{
				TeamName:        teamName,
				TeamDescription: teamDescription,
			},
			Expected: &db.Teams{
				TeamName:    teamName,
				Description: teamDescription,
			},
			ExpectedError: "",
		},
		{
			Name: "create a quality studio team unsuccessfully (with a name that already exists)",
			Input: &Input{
				TeamName:        teamName,
				TeamDescription: teamDescription,
			},
			Expected:      &db.Teams{},
			ExpectedError: "Already exists",
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			got, err := storage.CreateQualityStudioTeam(c.Input.TeamName, c.Input.TeamDescription)

			if err != nil || c.ExpectedError != "" {
				assert.EqualError(t, err, c.ExpectedError)
				return
			}

			assert.Equal(t, c.Expected.Description, got.Description)
		})
	}
}

func TestGetAllTeamsFromDB(t *testing.T) {
	// get database client
	cfg := GetPostgresConnectionDetails()
	storage, err := cfg.Open()
	require.NoError(t, err, "GetPostgresConnectionDetails() should not return an error")

	teamName := "team-" + util.GenerateRandomString(6)
	teamDescription := teamName

	toFind, err := storage.CreateQualityStudioTeam(teamName, teamDescription)
	assert.NoError(t, err)
	assert.Equal(t, teamName, toFind.TeamName)

	teams, err := storage.GetAllTeamsFromDB()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(teams), 1)

	if !exists(teams, toFind) {
		t.Errorf("GetAllTeamsFromDB() did not get %s team", toFind.TeamName)
	}
}

func TestGetTeamByName(t *testing.T) {
	// get database client
	cfg := GetPostgresConnectionDetails()
	storage, err := cfg.Open()
	require.NoError(t, err, "GetPostgresConnectionDetails() should not return an error")

	teamName := "team-" + util.GenerateRandomString(6)
	teamDescription := teamName

	expected, err := storage.CreateQualityStudioTeam(teamName, teamDescription)
	assert.NoError(t, err)
	assert.Equal(t, teamName, expected.TeamName)

	type Input struct {
		TeamName        string
		TeamDescription string
	}

	cases := []struct {
		Name          string
		Input         Input
		Expected      *db.Teams
		ExpectedError string
	}{
		{
			Name: "get a quality studio team by name successfully",
			Input: Input{
				TeamName: teamName,
			},
			Expected:      expected,
			ExpectedError: "",
		},
		{
			Name: "get a quality studio team by name unsuccessfully",
			Input: Input{
				TeamName: "team-not-exist",
			},
			Expected:      &db.Teams{},
			ExpectedError: "not found",
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			got, err := storage.GetTeamByName(c.Input.TeamName)

			if err != nil || c.ExpectedError != "" {
				assert.EqualError(t, err, c.ExpectedError)
				return
			}

			assert.Equal(t, c.Expected.TeamName, got.TeamName)
			assert.Equal(t, c.Expected.Description, got.Description)
			assert.Equal(t, c.Expected.ID, got.ID)
		})
	}
}

func exists(teams []*db.Teams, find *db.Teams) bool {
	for _, team := range teams {
		if team.TeamName == find.TeamName &&
			team.Description == find.Description &&
			team.ID == find.ID {
			return true
		}
	}
	return false
}
