package konflux

import "time"

type KonfluxMetadata struct {
	State           string    `json:"state"`
	JobId           string    `json:"job_id"`
	CreatedAt       time.Time `json:"created_at"`
	JobType         string    `json:"job_type"`
	JobName         string    `json:"job_name"`
	JobUrl          string    `json:"job_url"`
	ExternalImpact  bool      `json:"external_impact"`
	RepositoryName  string    `json:"repository_name"`
	GitOrganization string    `json:"git_organization"`
}
