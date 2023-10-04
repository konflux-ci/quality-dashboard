package jira

type Alert struct {
	AlertMessage *string `json:"alert_message"`
	Signal       string  `json:"signal"`
}

type BugSLO struct {
	JiraKey               string   `json:"jira_key"`
	JiraURL               string   `json:"jira_url"`
	TriageSLI             *Alert   `json:"triage_sli"`
	ResponseSLI           *Alert   `json:"response_sli"`
	ResolutionSLI         *Alert   `json:"resolution_sli"`
	DaysWithoutAssignee   *float64 `json:"days_without_assignee"`
	DaysWithoutPriority   *float64 `json:"days_without_priority"`
	DaysWithoutResolution *float64 `json:"days_without_resolution"`
}

type BugSLOInfo struct {
	Average float64 `json:"average"`
	Sum     float64 `json:"sum"`
	Total   int     `json:"total"`
}

type Project struct {
	ProjectKey                     string     `json:"project_key"`
	BugSLOs                        []BugSLO   `json:"bug_slos"`
	RedTriageTimeBugSLOInfo        BugSLOInfo `json:"red_triage_time_bug_slo_info"`
	YellowTriageTimeBugSLOInfo     BugSLOInfo `json:"yellow_triage_time_bug_slo_info"`
	RedResponseTimeBugSLOInfo      BugSLOInfo `json:"red_response_time_bug_slo_info"`
	RedResolutionTimeBugSLOInfo    BugSLOInfo `json:"red_resolution_time_bug_slo_info"`
	YellowResolutionTimeBugSLOInfo BugSLOInfo `json:"yellow_resolution_time_bug_slo_info"`
}
