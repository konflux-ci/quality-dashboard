package codecov

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
)

type API struct {
	httpClient   *http.Client
	githubAPIURL string
}

type TotalsSpec struct {
	Coverage json.Number `json:"coverage"`
}

type CoverageSpec struct {
	Totals TotalsSpec `json:"totals"`
}

func NewCodeCoverageClient() *API {
	api := API{
		githubAPIURL: "https://codecov.io/api/v2/github/",
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
	req, err := http.NewRequestWithContext(ctx, "GET", c.githubAPIURL+organization+"/repos/"+repository+"/report/", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)

	return c.Do(req)
}

func (c *API) GetCodeCovInfo(organization string, repo string) (repository CoverageSpec, err error) {
	gh := CoverageSpec{}
	response, err := c.Get(context.Background(), "aplication/json", organization, repo)

	if err != nil {
		return gh, err
	}

	if err = json.NewDecoder(response.Body).Decode(&gh); err != nil {
		return gh, err
	}

	return gh, nil
}
