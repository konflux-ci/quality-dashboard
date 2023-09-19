package v1alpha1

type Failure struct {
	JiraKey      string  `json:"jira_key"`
	JiraStatus   string  `json:"jira_status"`
	ErrorMessage string  `json:"error_message"`
	Frequency    float64 `json:"frequency"`
}

type Failures []Failure
