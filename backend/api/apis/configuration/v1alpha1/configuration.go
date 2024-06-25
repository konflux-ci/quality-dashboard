package v1alpha1

import (
	"github.com/google/uuid"
)

type Configuration struct {
	ID            uuid.UUID `json:"id"`
	TeamName      string    `json:"team_name"`
	JiraConfig    string    `json:"jira_config"`
	BugSLOsConfig string    `json:"bug_slos_config"`
}

type JiraConfig struct {
	BugsCollectQuery string   `json:"bugs_collect_query"`
	CiImpactQuery    string   `json:"ci_impact_query"`
	CiImpactBugs     []string `json:"ci_impact_bugs"`
}
