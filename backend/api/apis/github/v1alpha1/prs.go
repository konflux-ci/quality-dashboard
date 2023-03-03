package v1alpha1

import (
	"time"
)

// PullRequest represents pull request of a repository.
type PullRequest struct {
	// Title of the pull request.
	Title string `json:"title"`

	// Time when the pull request was created.
	CreatedAt time.Time `json:"created_at"`

	// Time when the pull request was closed.
	ClosedAt time.Time `json:"closed_at"`

	// Time when the pull request was merged.
	MergedAt time.Time `json:"merged_at"`

	// State of the pull request (open, closed).
	State string `json:"state"`

	// User who created the pull request.
	Author string `json:"author"`
}

// Summary represents all the collected information regarding all the pull requests of a repository.
type Summary struct {
	// Number of merged pull requests.
	MergedPrsCount int `json:"merged_prs"`

	// Number of open pull requests.
	OpenPrsCount int `json:"open_prs"`

	// Average time to merge a pull request.
	MergeAvg float64 `json:"merge_avg"`
}

// PullRequestsInfo represents the metrics by day and the summary of the pull requests.
type PullRequestsInfo struct {
	// Metadata of the pull requests.
	Summary Summary `json:"summary"`

	// Set of metrics about pull requests.
	Metrics []Metrics `json:"metrics"`
}

// Metrics represents the metrics by day of the pull requests.
type Metrics struct {
	// Date target of the collected metrics.
	Date string `json:"date"`

	// Number of pull requests created on the target day.
	CreatedPullRequestsCount int `json:"created_prs_count"`

	// Number of pull requests merged on the target day.
	MergedPullRequestsCount int `json:"merged_prs_count"`
}
