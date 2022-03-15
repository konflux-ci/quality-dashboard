package codecov

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"os"
)

type API struct {
	httpClient   *http.Client
	githubAPIURL string
}

func NewCodeCoverageClient() *API {
	api := API{
		githubAPIURL: "https://codecov.io/api/gh/",
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
	req.Header.Set("Authorization", os.Getenv("CODECOV_TOKEN"))
	return c.Do(req)
}

type TotalsSpec struct {
	TotalCoverage json.Number `json:"c"`
}

type CommitSpec struct {
	Totals TotalsSpec `json:"totals"`
}

type GitHubTagResponse struct {
	Commit CommitSpec `json:"commit"`
}

func (c *API) GetCodeCovInfo(organization string, repo string) (repository GitHubTagResponse, err error) {
	gh := GitHubTagResponse{}

	response, err := c.Get(context.Background(), "aplication/json", organization, repo)
	if err != nil {
		return gh, err
	}
	err = json.NewDecoder(response.Body).Decode(&gh)
	if err != nil {
		return gh, err
	}

	return gh, nil
}
