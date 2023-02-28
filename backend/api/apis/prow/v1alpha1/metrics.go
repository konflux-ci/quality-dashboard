package v1alpha1

// JobsMetrics return a set of ci metric for every repository in quality studio database
type JobsMetrics struct {
	// Indicate a GitHub repository name. "e2e-tests", "kubernetes"
	RepositoryName string `json:"repository_name"`

	// PostSubmits, periodics or presubmits
	JobType string `json:"type"`

	// Indicate a GitHub organization name
	GitOrganization string `json:"git_organization"`

	// A set of Jobs
	Jobs []Jobs `json:"jobs"`
}

// Metrics for specific job name
type Jobs struct {
	// Name of a prow job
	Name string `json:"name"`

	// Metadata about jobs
	Summary Summary `json:"summary"`

	// Set of metrics about job
	Metrics []Metrics `json:"metrics"`
}

// Set of Metrics for a specific job
type Metrics struct {
	// Return the percentage of success
	SuccessRate float64 `json:"success_rate"`

	// Return a percentage about how much fail a job
	FailureRate float64 `json:"failure_rate"`

	// Return percentage of prow ci failures
	CiFailedRate float64 `json:"ci_failed_rate"`

	Date string `json:"date"`
}

type Summary struct {
	DateFrom string `json:"date_from"`

	DateTo string `json:"date_to"`

	SuccessRateAvg float64 `json:"success_rate_avg"`

	JobFailedAvg float64 `json:"failure_rate_avg"`

	CIFailedAvg float64 `json:"ci_failed_rate_avg"`

	TotalJobs int `json:"total_jobs"`
}
