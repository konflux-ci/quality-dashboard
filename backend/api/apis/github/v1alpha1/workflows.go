package v1alpha1

// GitHub Actions specs.
type Workflow struct {
	// GitHub workflow name
	Name string `json:"workflow_name"`

	// Return a url to some GItHub action workflow badge
	BadgeURL string `json:"badge_url"`

	// Url of a workflow
	HTMLURL string `json:"html_url"`

	// Job Workflow URL
	JobURL string `json:"job_url"`

	State string `json:"state"`
}
