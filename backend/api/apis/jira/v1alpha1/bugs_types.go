package v1alpha1

import (
	"time"

	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
)

// Bugs is used to represent bugs in some serialized content.  It also tracks some additional metadata.
type JiraBug struct {
	// Unique key for a specific Jira: E.G. "STONE-<number>" or "AD-<number>"
	JiraKey string `json:"jira_key"`

	// Indicate the date when a bug was created.
	CreatedAt time.Time `json:"created_at"`

	// Indicate the date when a bug was updated.
	UpdatedAt time.Time `json:"updated_at"`

	// Indicate the date when a bug was resolved.
	ResolvedAt time.Time `json:"resolved_at"`

	// Specific bug is resolved.
	IsResolved bool `json:"resolved"`

	// Time to resolve a bug in jira
	ResolutionTime float64 `json:"resolution_time"`

	// Return information if a Jira is In Progress, Done etc.
	Status string `json:"status"`

	// Return information if a Jira is a blocker etc.
	Priority string `json:"priority"`

	// Jira Summary.
	Summary string `json:"summary"`

	// A complete link to the jira url.
	Url string `json:"url"`
}

type BugsMetrics struct {
	ResolutionTimeTotal ResolutionTime `json:"resolution_time"`
}

type ResolutionTime struct {
	Total float64 `json:"total"`

	Priority string `json:"priority"`

	NumberOfTotalBugs int `json:"resolved_bugs"`

	Months []MonthsResolution `json:"months"`
}

type MonthsResolution struct {
	Name string `json:"name"`

	Total float64 `json:"total"`

	NumberOfResolvedBugs int `json:"resolved_bugs"`

	Bugs []*db.Bugs `json:"bugs"`
}
