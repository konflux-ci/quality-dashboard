package github

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"

	"github.com/flacatus/qe-dashboard-backend/pkg/utils"
)

type API struct {
	httpClient   *http.Client
	githubAPIURL string
}

func NewGitubClient() *API {
	api := API{
		githubAPIURL: "https://api.github.com/repos/",
	}
	api.httpClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	return &api
}

func (c *API) Do(req *http.Request) (*http.Response, error) {
	res, err := c.httpClient.Do(req)
	return res, err
}

func (c *API) Get(ctx context.Context, contentType string, organization string, repository string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.githubAPIURL+organization+"/"+repository, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)

	if utils.CheckIfEnvironmentExists("GITHUB_TOKEN") {
		req.Header.Set("Authorization", "token "+os.Getenv("GITHUB_TOKEN"))

	}
	return c.Do(req)
}

func (c *API) GetWorkflows(ctx context.Context, contentType string, organization string, repository string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.githubAPIURL+organization+"/"+repository+"/actions/workflows", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)

	if utils.CheckIfEnvironmentExists("GITHUB_TOKEN") {
		req.Header.Set("Authorization", "token "+os.Getenv("GITHUB_TOKEN"))

	}
	return c.Do(req)
}
