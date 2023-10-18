package plugins

// Desired definition of a Plugin object
type PluginSpec struct {
	// Indicate the name of the desired plugin
	Name string `json:"name"`

	// Indicate the Plugin Category: Prow CI/Openshift CI, Jira, GitHub
	Category string `json:"category"`

	// The name of the logo to be used
	Logo string `json:"logo"`

	// Information about what is to suposed to do the plugin
	Description string `json:"description"`

	// Return if a plugin is available or not. Can be "Available or Deprecated"
	Reason string `json:"reason"`
}

type PluginStatus struct {
	// Indicate if the plugin is installed by a team
	Installed bool `json:"installed"`
}

type Plugin struct {
	Spec   PluginSpec   `json:"plugin"`
	Status PluginStatus `json:"status"`
}
