package jira

type Alert struct {
	AlertMessage *string `json:"alert_message"`
	Signal       string  `json:"signal"`
}

type Bug struct {
	JiraKey               string   `json:"jira_key"`
	JiraURL               string   `json:"jira_url"`
	Status                string   `json:"status"`
	Summary               string   `json:"summary"`
	Priority              string   `json:"priority"`
	Labels                string   `json:"labels"`
	Component             string   `json:"component"`
	TriageSLI             *Alert   `json:"triage_sli"`
	ResponseSLI           *Alert   `json:"response_sli"`
	ResolutionSLI         *Alert   `json:"resolution_sli"`
	GlobalSLI             string   `json:"global_sli"`
	DaysWithoutAssignee   *float64 `json:"days_without_assignee"`
	DaysWithoutPriority   *float64 `json:"days_without_priority"`
	DaysWithoutResolution *float64 `json:"days_without_resolution"`
}

type Metric struct {
	Average float64 `json:"average"`
	Sum     float64 `json:"sum"`
	Total   int     `json:"total"`
}

type GlobalSLI struct {
	GreenSLI  int `json:"green_sli"`
	RedSLI    int `json:"red_sli"`
	YellowSLI int `json:"yellow_sli"`
}

type SLI struct {
	Bugs   []Bug `json:"bugs"`
	Red    Metric
	Yellow Metric
}

type BugSlisInfo struct {
	GlobalSLI         GlobalSLI `json:"global_sli"`
	TriageTimeSLI     SLI       `json:"triage_time_sli"`
	ResponseTimeSLI   SLI       `json:"response_time_sli"`
	ResolutionTimeSLI SLI       `json:"resolution_time_sli"`
}
