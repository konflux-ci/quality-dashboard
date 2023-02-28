package codecov

// Spec to get Total Coverage for a specific repo.
type Coverage struct {
	// RepositoryName identify an GitHub repository
	RepositoryName string `json:"repository_name"`

	// RepositoryName identify an github repository
	GitOrganization string `json:"git_organization"`

	// CoveragePercentage identify the total percentage of a repo coverage
	CoveragePercentage float64 `json:"coverage_percentage"`
}
