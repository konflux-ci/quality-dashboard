package v1alpha1

import (
	"time"
)

const (
	// PullRequestClosedAt is used when filtering when a Pull Request was closed
	PullRequestClosedAt PullRequestTimeField = iota
	// PullRequestCreatedAt is used when filtering when a Pull Request was opened
	PullRequestCreatedAt
	// PullRequestMergedAt is used when filtering when a Pull Request was merged
	PullRequestMergedAt
	// PullRequestNone is used when the results are not filtered by time. Without any other filters, using this could easily cause an access token to be rate limited
	PullRequestNone
)

// PullRequest is a GitHub pull request
type PullRequest struct {
	Number     int
	Title      string
	URL        string
	State      string
	Author     PullRequestAuthor
	Closed     bool
	IsDraft    bool
	Locked     bool
	Merged     bool
	ClosedAt   time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	MergedAt   time.Time
	Mergeable  string
	MergedBy   *PullRequestAuthor
	Repository Repository
}

// PullRequests is a list of GitHub Pull Requests
type PullRequests []PullRequest

// PullRequestAuthor is the structure of the Author object in a Pull Request (which requires a graphQL object expansion on `User`)
type PullRequestAuthor struct {
	User User `graphql:"... on User"`
}

// ListPullRequestsOptions are the available options when listing pull requests in a time range
// PullRequestTimeField defines what time field to filter pull requests by (closed, opened, merged...)
type PullRequestTimeField uint32

// ListPullRequestsOptions are the available options when listing pull requests in a time range
type ListPullRequestsOptions struct {
	// Repository is the name of the repository being queried (ex: grafana)
	Repository string `json:"repository"`

	// Owner is the owner of the repository (ex: grafana)
	Owner string `json:"owner"`

	// TimeField defines what time field to filter by
	TimeField PullRequestTimeField `json:"timeField"`

	Query *string `json:"query,omitempty"`
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

// PullRequestOptionsWithRepo adds the Owner and Repository options to a ListPullRequestsOptions type
func PullRequestOptionsWithRepo(opt ListPullRequestsOptions, owner string, repo string) ListPullRequestsOptions {
	return ListPullRequestsOptions{
		Owner:      owner,
		Repository: repo,
		Query:      opt.Query,
		TimeField:  opt.TimeField,
	}
}

func (d PullRequestTimeField) String() string {
	return [...]string{"closed", "created", "merged", "opened"}[d]
}
