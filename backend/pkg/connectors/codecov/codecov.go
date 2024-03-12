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

type ResultSpec struct {
	CommitID string     `json:"commitid"`
	Totals   TotalsSpec `json:"totals"`
}

type ResultsSpec struct {
	Results []ResultSpec `json:"results"`
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

func (c *API) Get(ctx context.Context, contentType string, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)

	return c.Do(req)
}

func (c *API) GetCodeCovInfo(organization, repo string) (currentCov float64, covTrend string, err error) {
	res := ResultsSpec{}
	url := c.githubAPIURL + organization + "/repos/" + repo + "/commits?branch=main"

	response, err := c.Get(context.Background(), "aplication/json", url)
	if err != nil {
		return 0, "n/a", err
	}

	if err = json.NewDecoder(response.Body).Decode(&res); err != nil {
		return 0, "n/a", err
	}

	if len(res.Results) == 0 {
		return 0, "n/a", err
	}

	// codecov api returns results in desc order
	lastCov, err := res.Results[0].Totals.Coverage.Float64()
	if err != nil {
		return 0, "n/a", err
	}
	trending := "stable"

	if len(res.Results) >= 2 {
		penultimateCov, err := res.Results[1].Totals.Coverage.Float64()
		if err == nil {
			if lastCov > penultimateCov {
				trending = "ascending"
			} else if lastCov < penultimateCov {
				trending = "descending"
			}
		}
	}

	return lastCov, trending, nil
}
