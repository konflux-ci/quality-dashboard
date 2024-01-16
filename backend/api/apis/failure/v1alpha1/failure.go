package v1alpha1

import (
	"time"

	"github.com/google/uuid"
)

type Failure struct {
	TeamName      string    `json:"team"`            // Team Name as defined in QD
	TeamID        uuid.UUID `json:"team_id"`         // Team ID as generated for QD
	JiraID        uuid.UUID `json:"jira_id"`         // Jira ID in QD
	JiraKey       string    `json:"jira_key"`        // Jira Identfier fetched from Jira
	JiraStatus    string    `json:"jira_status"`     // Jira status fetched from Jira
	ErrorMessage  string    `json:"error_message"`   // Error messag from logs
	Frequency     float64   `json:"frequency"`       // Precentage of number times this issue occuers
	TitleFromJira string    `json:"title_from_jira"` // Fetched from Jira Summary
	CreatedDate   time.Time `json:"created_date"`    // Date of creating this Jira, fetched from Jira
	ClosedDate    time.Time `json:"closed_date"`     // Date of closing Jira fetched from Jira
	Labels        string    `json:"labels"`          // Labels fetched from Jira
}

type Failures []Failure
