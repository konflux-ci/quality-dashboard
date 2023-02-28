package v1alpha1

import (
	"time"
)

// Bugs is used to represent bugs in some serialized content.  It also tracks some additional metadata.
type JiraBug struct {
	// Unique key for a specific Jira: E.G. "STONE-<number>" or "AD-<number>"
	JiraKey string `json:"jira_key"`

	// Indicate the date when a bug was created.
	CreatedAt time.Time `json:"created_at"`

	// Indicate the date when a bug was updated.
	UpdatedAt time.Time `json:"updated_at"`

	// Return information if a Jira is In Progress, Done etc.
	Status string `json:"status"`

	// Return information if a Jira is a blocker etc.
	Priority string `json:"priority"`

	// Jira Summary.
	Summary string `json:"summary"`

	// A complete link to the jira url.
	Url string `json:"url"`
}
