package storage

import (
	"errors"

	"github.com/google/uuid"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
)

var (
	// ErrNotFound is the error returned by storages if a resource cannot be found.
	ErrNotFound = errors.New("not found")

	// ErrAlreadyExists is the error returned by storages if a resource ID is taken during a create.
	ErrAlreadyExists = errors.New("Already exists")
)

// Storage is the storage interface used by the server. Implementations are
// required to be able to perform atomic compare-and-swap updates and either
// support timezones or standardize on UTC.
type Storage interface {
	Close() error

	// GET
	GetRepository(repositoryName string, gitOrganizationName string) (*db.Repository, error)
	GetProwJobsResults(*db.Repository) ([]*db.Prow, error)
	GetProwJobsResultsByJobID(jobID string) ([]*db.Prow, error)
	ListRepositories() ([]Repository, error)
	ListWorkflowsByRepository(repositoryName string) (w []GithubWorkflows, err error)
	ListRepositoriesQualityInfo() ([]RepositoryQualityInfo, error)

	// POST
	CreateRepository(p Repository) (*db.Repository, error)

	CreateWorkflows(p GithubWorkflows, repo_id uuid.UUID) error

	// POST
	CreateCoverage(p Coverage, repo_id uuid.UUID) error
	CreateProwJobResults(repository ProwJob, repo_id uuid.UUID) error
	// Delete
	ReCreateWorkflow(workflow GithubWorkflows, repoName string) error
	UpdateCoverage(codecov Coverage, repoName string) error
	DeleteRepository(repositoryName string, gitOrganizationName string) error
}

// Repository is an github repository info managed by the storage.
type Repository struct {
	// RepositoryName identify an github repository
	RepositoryName string `json:"repository_name"`

	GitOrganization string `json:"git_organization"`

	Description string `json:"description"`

	GitURL string `json:"git_url"`
}

type RepositoryQualityInfo struct {
	// RepositoryName identify an github repository
	RepositoryName string `json:"repository_name"`

	GitOrganization string `json:"git_organization"`

	Description string `json:"description"`

	GitURL string `json:"git_url"`

	//Coverage
	CI []GithubWorkflows `json:"github_actions"`

	CodeCoverage Coverage `json:"code_coverage"`
}

// Repository is an github repository info managed by the storage.
type Coverage struct {
	// RepositoryName identify an github repository
	RepositoryName string `json:"repository_name"`

	GitOrganization string `json:"git_organization"`

	CoveragePercentage float64 `json:"coverage_percentage"`
}

// Repository is an github repository info managed by the storage.
type ProwJob struct {
	JobID string `json:"job_id"`

	TestCaseName string `json:"test_name"`

	TestCaseStatus string `json:"test_status"`

	TestTiming float64 `json:"test_timing"`
}

// Repository is an github repository info managed by the storage.
type GithubWorkflows struct {
	// RepositoryName identify an github repository
	WorkflowName string `json:"workflow_name"`

	BadgeURL string `json:"badge_url"`

	HTMLURL string `json:"html_url"`

	JobURL string `json:"job_url"`

	State string `json:"state"`
}
