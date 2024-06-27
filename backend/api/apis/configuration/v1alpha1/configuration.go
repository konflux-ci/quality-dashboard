package v1alpha1

import (
	"github.com/google/uuid"
)

// Configuration holds all the configuration related to team specifications.
type Configuration struct {
	ID            uuid.UUID `json:"id"`
	TeamName      string    `json:"team_name"`
	JiraConfig    string    `json:"jira_config"`
	BugSLOsConfig string    `json:"bug_slos_config"`
}

// JiraConfig holds the configuration for the Jira plugin.
type JiraConfig struct {
	// JQL query to save all the bugs for a specific team.
	BugsCollectQuery string `json:"bugs_collect_query"`
	// JQL query to grab all the bugs that are affecting CI.
	CiImpactQuery string `json:"ci_impact_query"`
	// List of Jira keys related to the bugs impacting CI based on the CiImpactQuery.
	CiImpactBugs []string `json:"ci_impact_bugs"`
}
