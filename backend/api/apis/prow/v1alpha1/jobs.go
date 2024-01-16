package v1alpha1

import (
	"time"

	"github.com/redhat-appstudio/quality-studio/api/apis/github/v1alpha1"
)

// Prow Jobs desired specs.
type Job struct {
	// Unique ID generated for every prow job.
	JobID string `json:"job_id"`

	// Date when a job was created in Openshift ci
	CreatedAt time.Time `json:"created_at"`

	// State indicate if a job is running, failed, aborted etc
	State string `json:"state"`

	// Return how much a job takes to finalize
	Duration float64 `json:"duration"`

	// Number of the tests running in a job
	TestsCount int64 `json:"tests_count"`

	// Number of tests failed
	FailedCount int64 `json:"failed_count"`

	// Number of tests skipped
	SkippedCount int64 `json:"skipped_count"`

	// PostSubmits, periodics or presubmits
	JobType string `json:"job_type"`

	// Indicate the name of the job
	JobName string `json:"job_name"`

	// Url to some prow cluster
	JobURL string `json:"job_url"`

	// Indicate if the test infrastructure failed or not
	CIFailed int16 `json:"ci_failed"`

	// Indicate if job was impacted by an external service
	ExternalServiceImpact bool `json:"external_service_impact"`

	// Set of failed test cases from /artifacts/redhat-appstudio-e2e/redhat-appstudio-e2e/artifacts/e2e-report.xml
	E2EFailedTestMessages string `json:"e2e_failed_test_messages"`

	// URL to the e2e-report.xml of the job
	SuitesXmlUrl string `json:"suites_xml_url"`

	// Set of error logs from build-log.txt
	BuildErrorLogs string `json:"build_error_logs"`
}

type ProwRepository struct {
	Repository v1alpha1.Repository

	JobsList []string `json:"jobs"`
}
