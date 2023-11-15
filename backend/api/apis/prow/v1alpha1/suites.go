package v1alpha1

import "time"

// Indicate the test suites specs for a job. Need to be in XUNIT format
type JobSuites struct {
	// Job ID identification
	JobID string `json:"job_id"`

	// Job URL of a specific test suite
	JobURL string `json:"job_url"`

	// JobName the name of the job which a suite belongs
	JobName string `json:"job_name"`

	// the name of suite
	SuiteName string `json:"suite_name"`

	// Test name for a specific case
	TestCaseName string `json:"test_name"`

	// Indicate the status for a specific tests
	TestCaseStatus string `json:"test_status"`

	// Return the total time for a test
	TestTiming float64 `json:"test_timing"`

	// Return if is postsubmit, presubmit or periodic
	JobType string `json:"job_type"`

	CreatedAt time.Time `json:"created_at"`

	ErrorMessage string `json:"error_message"`
}

type FlakyFrequency struct {
	// Indicate percentage of the flaky tests impact in ci jobs
	GlobalImpact float64 `json:"global_impact"`

	// Indicate a GitHub organization name
	GitOrganization string `json:"git_organization"`

	// Indicate a GitHub repository name. "e2e-tests", "kubernetes"
	RepositoryName string `json:"repository_name"`

	// JobName the name of the job which a suite belongs
	JobName string `json:"job_name"`

	SuitesFailureFrequency []SuitesFailureFrequency `json:"suites"`
}

type SuitesFailureFrequency struct {
	SuiteName string `json:"suite_name"`

	Status string `json:"status"`

	TestCases []TestCases `json:"test_cases"`

	AverageImpact float64 `json:"average_impact"`
}

type TestCases struct {
	Name string `json:"name"`

	TestCaseImpact float64 `json:"test_case_impact"`

	Count int `json:"count"`

	Messages []Messages `json:"messages"`
}

type Messages struct {
	JobId string `json:"job_id"`

	JobURL string `json:"job_url"`

	Message string `json:"error_message"`

	FailureDate *time.Time `json:"failure_date"`
}

type FlakyMetrics struct {
	Date string `jaon:"date"`

	GlobalImpact float64 `json:"global_impact"`
}
