package v1alpha1

import "github.com/google/uuid"

type Failure struct {
	TeamName     string    `json:"team"`
	TeamID       uuid.UUID `json:"team_id"`
	JiraID       uuid.UUID `json:"jira_id"`
	JiraKey      string    `json:"jira_key"`
	JiraStatus   string    `json:"jira_status"`
	ErrorMessage string    `json:"error_message"`
	Frequency    float64   `json:"frequency"`
}

type Failures []Failure
