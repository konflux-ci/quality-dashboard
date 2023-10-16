package plugins

// Desired definition of a Plugin object
type Plugin struct {
	// Indicate the name of the desired plugin
	Name string `json:"name"`

	// Indicate the Plugin Category: Prow CI/Openshift CI, Jira, GitHub
	Category string `json:"category"`

	// The name of the logo to be used
	Logo string `json:"logo"`

	// Information about what is to suposed to do the plugin
	Description string `json:"description"`

	// Desired status of the plugin. Can be "Available or Deprecated"
	Status string `json:"status"`
}
