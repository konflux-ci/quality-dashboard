package github

import (
	"context"
	"errors"
	"sync"

	"github.com/google/go-github/v44/github"
)

var (
	GENERIC_ERROR_GIT_RESPONSE = errors.New("unable to parse repository. Please verify your token")
)

func (g *Github) GetGithubRepositoryInformation(organization string, repository string) (*github.Repository, error) {
	repo, resp, err := g.client.Repositories.Get(context.Background(), organization, repository)
	if resp.StatusCode != 200 {
		return nil, GENERIC_ERROR_GIT_RESPONSE
	}
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func (g *Github) GetRepositoryWorkflows(organization string, repository string) (*github.Workflows, error) {
	wk, resp, err := g.client.Actions.ListWorkflows(context.Background(), organization, repository, &github.ListOptions{})
	if resp.StatusCode != 200 {
		return nil, GENERIC_ERROR_GIT_RESPONSE
	}

	if err != nil {
		return nil, err
	}
	return wk, nil
}

func (g *Github) GetRepositoryPullRequests(organization string, repository string) ([]*github.PullRequest, error) {
	prs, resp, err := g.client.PullRequests.List(
		context.Background(),
		organization,
		repository,
		&github.PullRequestListOptions{
			State: "all",
			ListOptions: github.ListOptions{
				Page:    0,
				PerPage: 100,
			},
		})

	if resp.StatusCode != 200 {
		return nil, GENERIC_ERROR_GIT_RESPONSE
	}

	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	var m sync.Mutex

	for page := resp.NextPage; page <= resp.LastPage; page++ {
		wg.Add(1)
		page := page

		go func(page int) {
			defer wg.Done()

			prsNextPage, _, errNextPage := g.client.PullRequests.List(
				context.Background(),
				organization,
				repository,
				&github.PullRequestListOptions{
					State: "all",
					ListOptions: github.ListOptions{
						PerPage: 100,
						Page:    page,
					},
				})

			if errNextPage != nil {
				return
			}

			m.Lock()
			prs = append(prs, prsNextPage...)
			m.Unlock()
		}(page)
	}

	wg.Wait()

	return prs, nil
}
