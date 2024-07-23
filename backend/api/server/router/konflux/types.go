package konflux

import (
	"time"
)

// KonfluxMetadata Represents metadata of tests used in Konflux metrics
type KonfluxMetadata struct {
	// State indicate if a job is running, failed, aborted etc
	State string `json:"state"`

	// Unique ID generated for every prow job.
	JobId string `json:"job_id"`

	// Date when a job was created in Openshift ci
	CreatedAt time.Time `json:"created_at"`

	// PostSubmits, periodics or presubmits
	JobType string `json:"job_type"`

	// Indicate the name of the job
	JobName string `json:"job_name"`

	// Url to some prow cluster
	JobUrl string `json:"job_url"`

	// Indicate if job was impacted by an external service
	ExternalImpact bool `json:"external_impact"`

	// Name of the repository where job was ran
	RepositoryName string `json:"repository_name"`

	// Name of the git organization where job was ran
	GitOrganization string `json:"git_organization"`
}
