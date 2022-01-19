package github

import (
	"context"
	"encoding/json"
)

type WorkflowSpec struct {
	Name     string `json:"name"`
	BadgeURL string `json:"badge_url"`
	HTML_URL string `json:"html_url"`
	JobURL   string `json:"job_url"`
	State    string `json:"state"`
}

type GitHubActionsResponse struct {
	Workflows []WorkflowSpec `json:"workflows"`
}

func (c *API) GetRepositoryWorkflows(organization string, repo string) (repository GitHubActionsResponse, err error) {
	gh := GitHubActionsResponse{}

	response, err := c.GetWorkflows(context.Background(), "aplication/json", organization, repo)
	if err != nil {
		return gh, err
	}
	err = json.NewDecoder(response.Body).Decode(&gh)
	if err != nil {
		panic(err)
	}

	return gh, nil
}
