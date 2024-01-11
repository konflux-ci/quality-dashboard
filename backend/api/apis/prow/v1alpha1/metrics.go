package v1alpha1

// JobsMetrics return a set of ci metric for every repository in quality studio database
type JobsMetrics struct {
	// Indicate a GitHub organization name
	GitOrganization string `json:"git_organization"`

	// Indicate a GitHub repository name. "e2e-tests", "kubernetes"
	RepositoryName string `json:"repository_name"`

	// Name of a prow job
	JobName string `json:"name"`

	// From date where start to measure jobs impact
	StartDate string `json:"start_date"`

	// // From date where stop to measure jobs impact
	EndDate string `json:"end_date"`

	// A set of metrics about a specific job name
	JobsRuns JobsRuns `json:"jobs_runs"`

	// Metrics which show how CI jobs are impacted
	JobsImpacts JobsImpacts `json:"jobs_impacts"`
}

// Return basic metrics for a specific job in prow
type JobsRuns struct {
	// All the jobs executed for a job
	Total int `json:"total"`

	// Total number of jobs which finished in success
	Success int `json:"success"`

	// Total number of jobs which finished in fail
	Failures int `json:"failures"`

	// Percentage of all success jobs executed for a specific job
	SuccessPercentage float64 `json:"success_percentage"`

	// Percentage of all success jobs executed for a specific job
	FailedPercentage float64 `json:"failed_percentage"`
}

// Indicate why Openshift CI jobs are failing
type JobsImpacts struct {
	// Show trends of infrastructure impact in ci jobs
	InfrastructureImpact InfrastructureImpact `json:"infrastructure_impact"`

	// Show trends about flaky tests impact in ci jobs
	FlakyTestsImpact FlakyTestsImpact `json:"flaky_tests_impact"`

	// Show trends about flaky tests impact in ci jobs
	ExternalServicesImpact ExternalServicesImpact `json:"external_services_impact"`

	// Show trends about undetected failures by quality dashboard impact in ci jobs
	UnknowFailuresImpact UnknowFailuresImpact `json:"unknown_failures_impact"`
}

type InfrastructureImpact struct {
	// Number of total jobs executed in ci
	Total int `json:"total"`

	// The percentage of the total job impacted
	Percentage float64 `json:"percentage"`
}

type FlakyTestsImpact struct {
	// Number of total jobs executed in ci
	Total int `json:"total"`

	// The percentage of the total job impacted
	Percentage float64 `json:"percentage"`
}

type ExternalServicesImpact struct {
	// Number of total jobs executed in ci
	Total int `json:"total"`

	// The percentage of the total job impacted
	Percentage float64 `json:"percentage"`
}

type UnknowFailuresImpact struct {
	// Number of total jobs executed in ci
	Total int `json:"total"`

	// The percentage of the total job impacted
	Percentage float64 `json:"percentage"`
}
