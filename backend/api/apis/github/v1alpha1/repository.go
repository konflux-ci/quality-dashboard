package v1alpha1

import "github.com/google/uuid"

// RepositoriesService handles communication with the repository related
// methods of the GitHub API.
//
// GitHub API docs: https://docs.github.com/en/rest/repos/
// Repository represents a GitHub repository.
type Repository struct {
	// Indicate a GitHub repository name. "e2e-tests", "kubernetes"
	Name string `json:"name,omitempty"`

	// Indicate a GitHub organization name
	Organization string `json:"organization,omitempty"`

	// A valid description for a GitHub Repository
	Description string `json:"description,omitempty"`

	// Link to a GitHub url
	HTMLURL string `json:"html_url,omitempty"`

	// A valid ID of a GitHub Repository
	ID uuid.UUID `json:"id,omitempty"`
}
