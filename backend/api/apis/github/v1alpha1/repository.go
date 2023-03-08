package v1alpha1

import (
	"github.com/shurcooL/githubv4"
)

type Owner struct {
	Login string
}

// Repository is a code repository
// RepositoriesService handles communication with the repository related
// methods of the GitHub API.
//
// GitHub API docs: https://docs.github.com/en/rest/repos/
// Repository represents a GitHub repository.
type Repository struct {
	ID string
	// Indicate a GitHub repository name. "e2e-tests", "kubernetes"
	Name string

	Description string
	// Indicate a GitHub organization name
	NameWithOwner string
	Owner         Owner
	// Link to a GitHub url
	URL       string
	ForkCount int64
	IsFork    bool
	IsMirror  bool
	IsPrivate bool
	CreatedAt githubv4.DateTime
}
