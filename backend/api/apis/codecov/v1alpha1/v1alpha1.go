package codecov

// Spec to get Total Coverage for a specific repo.
type Coverage struct {
	// RepositoryName identifies a GitHub repository
	RepositoryName string `json:"repository_name"`

	// RepositoryName identifies a GitHub repository
	GitOrganization string `json:"git_organization"`

	// Metric to determine the average of retest in a pull request
	AverageToRetestPullRequest float64 `json:"average_to_retest_before_merge"`

	// CoveragePercentage identifies the total percentage of a repo coverage
	CoveragePercentage float64 `json:"coverage_percentage"`

	// CoverageTrend identifies the coverage trend between the two last commits
	CoverageTrend string `json:"coverage_trend"`
}
