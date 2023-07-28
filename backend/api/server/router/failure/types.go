package failure

type FailureRequest struct {
	Team         string `json:"team"`
	ErrorMessage string `json:"error_message"`
	JiraKey      string `json:"jira_key"`
}
