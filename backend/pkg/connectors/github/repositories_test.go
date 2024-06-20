package github

import (
	"testing"

	util "github.com/konflux-ci/quality-studio/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetGithubRepositoryInformation(t *testing.T) {
	gh := NewGithubClient(util.GetEnv("GITHUB_TOKEN", ""))

	cases := []struct {
		Name            string
		ExpectedError   string
		RepositoryName  string
		GitOrganization string
	}{
		{
			Name:            "get github repository information successfully",
			ExpectedError:   "",
			RepositoryName:  "quality-dashboard",
			GitOrganization: "konflux-ci",
		},
		{
			Name:            "get github repository information unsuccessfully",
			ExpectedError:   "unable to parse repository. Please verify your token",
			RepositoryName:  "",
			GitOrganization: "",
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			got, err := gh.GetGithubRepositoryInformation(c.GitOrganization, c.RepositoryName)
			if err != nil || c.ExpectedError != "" {
				return
			}

			assert.Equal(t, c.RepositoryName, *got.Name)
		})
	}
}

func TestGetRepositoryWorkflows(t *testing.T) {
	gh := NewGithubClient(util.GetEnv("GITHUB_TOKEN", ""))

	cases := []struct {
		Name               string
		ExpectedError      string
		RepositoryName     string
		GitOrganization    string
		WorkflowTotalCount int
	}{
		{
			Name:               "get repository workflows successfully",
			ExpectedError:      "",
			RepositoryName:     "quality-dashboard",
			GitOrganization:    "konflux-ci",
			WorkflowTotalCount: 0,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			got, err := gh.GetRepositoryWorkflows(c.GitOrganization, c.RepositoryName)

			if err != nil || c.ExpectedError != "" {
				return
			}

			assert.GreaterOrEqual(t, c.WorkflowTotalCount, (*got).GetTotalCount())
		})
	}
}
