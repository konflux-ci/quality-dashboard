package codecov

import (
	"testing"

	repoV1Alpha1 "github.com/konflux-ci/quality-studio/api/apis/github/v1alpha1"
	"github.com/konflux-ci/quality-studio/pkg/storage/ent/client"
	util "github.com/konflux-ci/quality-studio/pkg/utils"
	"github.com/stretchr/testify/assert"
)

var repository = repoV1Alpha1.Repository{
	ID:   "12345678",
	Name: "managed-gitops",
	Owner: repoV1Alpha1.Owner{
		Login: "konflux-ci",
	},
	Description: "GitOps Service: Backend/cluster-agent/utility components aiming to provided GitOps services via Kubernetes-controller-managed Argo CD",
	URL:         "https://github.com/redhat-appstudio/managed-gitops",
}

func TestGetCodeCovInfo(t *testing.T) {
	cfg := client.GetPostgresConnectionDetails()
	storage, _, err := cfg.Open()
	assert.NoError(t, err)

	// be sure that there is no test repo in the db
	err = storage.DeleteRepository(repository.Name, repository.Owner.Login)
	assert.NoError(t, err)

	teamName := "team-" + util.GenerateRandomString(6)
	teamDescription := teamName

	// create team with one repo associated
	team, err := storage.CreateQualityStudioTeam(teamName, teamDescription, "team_jira")
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
		ExpectedCov     float64
	}{
		{
			Name:            "get codecov info successfully",
			WantError:       false,
			RepositoryName:  repo.RepositoryName,
			GitOrganization: repo.GitOrganization,
			ExpectedCov:     10,
		},
		{
			Name:            "get codecov info unsuccessfully",
			WantError:       false,
			RepositoryName:  "",
			GitOrganization: "",
			ExpectedCov:     0,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			currentCov, _, err := api.GetCodeCovInfo(c.GitOrganization, c.RepositoryName)
			if c.WantError != (err != nil) {
				t.Errorf("GetCodeCovInfo(%s, %s) got error = %v, want error %v", c.GitOrganization, c.RepositoryName, err, c.WantError)
				return
			}
			assert.GreaterOrEqual(t, currentCov, c.ExpectedCov)
		})
	}
}
