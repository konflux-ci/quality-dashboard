package github

import (
	"context"

	"github.com/google/go-github/v44/github"
)

func (g *Github) GetGithubRepositoryInformation(organization string, repository string) (*github.Repository, error) {
	repo, _, err := g.client.Repositories.Get(context.Background(), organization, repository)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func (g *Github) GetRepositoryWorkflows(organization string, repository string) (*github.Workflows, error) {
	wk, _, err := g.client.Actions.ListWorkflows(context.Background(), organization, repository, &github.ListOptions{})

	if err != nil {
		return nil, err
	}
	return wk, nil
}
