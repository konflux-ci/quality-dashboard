package codecov

import (
	"testing"

	"github.com/devfile/library/pkg/util"
	repoV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/github/v1alpha1"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/client"
	"github.com/stretchr/testify/assert"
)

var repository = repoV1Alpha1.Repository{
	Name:         "managed-gitops",
	Organization: "redhat-appstudio",
	Description:  "GitOps Service: Backend/cluster-agent/utility components aiming to provided GitOps services via Kubernetes-controller-managed Argo CD",
	HTMLURL:      "https://github.com/redhat-appstudio/managed-gitops",
}

func TestGetCodeCovInfo(t *testing.T) {
	cfg := client.GetPostgresConnectionDetails()
	storage, _, err := cfg.Open()
	assert.NoError(t, err)

	// be sure that there is no test repo in the db
	err = storage.DeleteRepository(repository.Name, repository.Organization)
	assert.NoError(t, err)

	teamName := "team-" + util.GenerateRandomString(6)
	teamDescription := teamName

	// create team with one repo associated
	team, err := storage.CreateQualityStudioTeam(teamName, teamDescription)
	assert.NoError(t, err)
	assert.Equal(t, teamName, team.TeamName)

	repo, err := storage.CreateRepository(repository, team.ID)
	assert.NoError(t, err)

	// create codecov client
	api := NewCodeCoverageClient()

	cases := []struct {
		Name            string
		WantError       bool
		RepositoryName  string
		GitOrganization string
		ExpectedCov     string
	}{
		{
			Name:            "get codecov info successfully",
			WantError:       false,
			RepositoryName:  repo.RepositoryName,
			GitOrganization: repo.GitOrganization,
			ExpectedCov:     "10",
		},
		{
			Name:            "get codecov info unsuccessfully",
			WantError:       false, // GetCodeCovInfo is returning no error after the change codecov api to use v2
			RepositoryName:  "",
			GitOrganization: "",
			ExpectedCov:     "",
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			got, err := api.GetCodeCovInfo(c.GitOrganization, c.RepositoryName)
			if c.WantError != (err != nil) {
				t.Errorf("GetCodeCovInfo(%s, %s) got error = %v, want error %v", c.GitOrganization, c.RepositoryName, err, c.WantError)
				return
			}
			assert.GreaterOrEqual(t, got.Totals.Coverage.String(), c.ExpectedCov)
		})
	}
}
