package client

import (
	"testing"

	"github.com/devfile/library/pkg/util"
	"github.com/google/uuid"
	s "github.com/redhat-appstudio/quality-studio/api/apis/github/v1alpha1"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
	"github.com/stretchr/testify/assert"
)

var toCreate = s.Repository{
	Name:         "managed-gitops",
	Organization: "redhat-appstudio",
	Description:  "GitOps Service: Backend/cluster-agent/utility components aiming to provided GitOps services via Kubernetes-controller-managed Argo CD",
	HTMLURL:      "https://github.com/redhat-appstudio/managed-gitops",
}

func TestCreateRepository(t *testing.T) {
	// get database client
	cfg := GetPostgresConnectionDetails()
	storage, err := cfg.Open()
	assert.NoError(t, err)

	// be sure that there is no test repo in the db
	err = storage.DeleteRepository(toCreate.Name, toCreate.Organization)
	assert.NoError(t, err)

	teamName := "team-" + util.GenerateRandomString(6)
	teamDescription := teamName

	// create a team
	team, err := storage.CreateQualityStudioTeam(teamName, teamDescription)
	assert.NoError(t, err)
	assert.Equal(t, teamName, team.TeamName)

	cases := []struct {
		Name          string
		Input         *s.Repository
		Expected      *db.Repository
		ExpectedError string
		Team          *uuid.UUID
	}{
		{
			Name:  "create a repository successfully",
			Input: &toCreate,
			Expected: &db.Repository{
				RepositoryName:  toCreate.Name,
				GitOrganization: toCreate.Organization,
				Description:     toCreate.Description,
				GitURL:          toCreate.HTMLURL,
			},
			ExpectedError: "",
			Team:          &team.ID,
		},
		{
			Name:          "create a repository unsuccessfully (with the same repo and team)",
			Input:         &toCreate,
			Expected:      &db.Repository{},
			ExpectedError: "Already exists",
			Team:          &team.ID,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			got, err := storage.CreateRepository(*c.Input, *c.Team)

			if err != nil || c.ExpectedError != "" {
				assert.EqualError(t, err, c.ExpectedError)
				return
			}

			assert.Equal(t, c.Input.Name, got.RepositoryName)
		})
	}
}

func TestListRepositories(t *testing.T) {
	// get database client
	cfg := GetPostgresConnectionDetails()
	storage, err := cfg.Open()
	assert.NoError(t, err)

	// be sure that there is no test repo in the db
	err = storage.DeleteRepository(toCreate.Name, toCreate.Organization)
	assert.NoError(t, err)

	teamName := "team-" + util.GenerateRandomString(6)
	teamDescription := teamName

	// create team1 without any repo associated
	team1, err := storage.CreateQualityStudioTeam(teamName, teamDescription)
	assert.NoError(t, err)
	assert.Equal(t, teamName, team1.TeamName)

	// create team2 with one repo associated
	team2, err := storage.CreateQualityStudioTeam(teamName+"-", teamDescription+"-")
	assert.NoError(t, err)
	assert.Equal(t, teamName+"-", team2.TeamName)

	_, err = storage.CreateRepository(toCreate, team2.ID)
	assert.NoError(t, err)

	cases := []struct {
		Name                string
		ExpectedReposNumber int
		Team                *db.Teams
	}{
		{
			Name:                "empty repository list",
			ExpectedReposNumber: 0,
			Team:                team1,
		},
		{
			Name:                "filled repository list",
			ExpectedReposNumber: 1,
			Team:                team2,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			repos, err := storage.ListRepositories(c.Team)
			assert.NoError(t, err)
			assert.Equal(t, c.ExpectedReposNumber, len(repos))
		})
	}
}

func TestGetRepository(t *testing.T) {
	// get database client
	cfg := GetPostgresConnectionDetails()
	storage, err := cfg.Open()
	assert.NoError(t, err)

	// be sure that there is no test repo in the db
	err = storage.DeleteRepository(toCreate.Name, toCreate.Organization)
	assert.NoError(t, err)

	teamName := "team-" + util.GenerateRandomString(6)
	teamDescription := teamName

	// create team with one repo associated
	team, err := storage.CreateQualityStudioTeam(teamName, teamDescription)
	assert.NoError(t, err)
	assert.Equal(t, teamName, team.TeamName)

	repo, err := storage.CreateRepository(toCreate, team.ID)
	assert.NoError(t, err)

	type Input struct {
		RepositoryName      string
		GitOrganizationName string
	}

	cases := []struct {
		Name          string
		Input         Input
		Expected      *db.Repository
		ExpectedError string
	}{
		{
			Name: "get repository successfully",
			Input: Input{
				RepositoryName:      repo.RepositoryName,
				GitOrganizationName: repo.GitOrganization,
			},
			Expected:      repo,
			ExpectedError: "",
		},
		{
			Name: "get repository unsuccessfully",
			Input: Input{
				RepositoryName:      "repository-not-exist",
				GitOrganizationName: "git-org-not-exist",
			},
			Expected:      &db.Repository{},
			ExpectedError: "not found",
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			got, err := storage.GetRepository(c.Input.RepositoryName, c.Input.GitOrganizationName)

			if err != nil || c.ExpectedError != "" {
				assert.EqualError(t, err, c.ExpectedError)
				return
			}

			assert.Equal(t, c.Expected.RepositoryName, got.RepositoryName)
			assert.Equal(t, c.Expected.GitOrganization, got.GitOrganization)
		})
	}
}
