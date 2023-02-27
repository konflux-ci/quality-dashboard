package storage

import (
	"errors"
	"time"

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
	GetLatestProwTestExecution(r *db.Repository, jobType string) (*db.ProwJobs, error)
	GetSuitesByJobID(jobID string) ([]*db.ProwSuites, error)
	GetProwJobsResults(*db.Repository) ([]*db.ProwSuites, error)
	GetProwJobsResultsByJobID(jobID string) ([]*db.ProwJobs, error)
	GetMetrics(gitOrganization string, repoName string, jobType string, startDate string, endDate string) ProwJobsMetrics
	GetAllTeamsFromDB() ([]*db.Teams, error)
	GetTeamByName(teamName string) (*db.Teams, error)
	ListWorkflowsByRepository(repositoryName string) (w []GithubWorkflows, err error)
	ListRepositories(team *db.Teams) ([]Repository, error)
	ListRepositoriesQualityInfo(team *db.Teams) ([]RepositoryQualityInfo, error)

	// POST
	CreateRepository(p Repository, team_id uuid.UUID) (*db.Repository, error)
	CreateQualityStudioTeam(teamName string, description string) (*db.Teams, error)
	CreateWorkflows(p GithubWorkflows, repo_id uuid.UUID) error
	CreateCoverage(p Coverage, repo_id uuid.UUID) error
	CreateProwJobSuites(prowJobStatus ProwJobSuites, repo_id uuid.UUID) error
	CreateProwJobResults(prowJobStatus ProwJobStatus, repo_id uuid.UUID) error
	ReCreateWorkflow(workflow GithubWorkflows, repoName string) error

	// UPDATE
	UpdateCoverage(codecov Coverage, repoName string) error
	UpdateTeam(t *db.Teams, target string) error

	// Delete
	DeleteRepository(repositoryName string, gitOrganizationName string) error
	DeleteTeam(teamName string) (bool, error)
}

// Repository is an github repository info managed by the storage.
type Repository struct {
	// RepositoryName identify an github repository
	RepositoryName string `json:"repository_name"`

	GitOrganization string `json:"git_organization"`

	Description string `json:"description"`

	GitURL string `json:"git_url"`

	ID uuid.UUID `json:"id"`
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
type ProwJobSuites struct {
	JobID string `json:"job_id"`

	TestCaseName string `json:"test_name"`

	TestCaseStatus string `json:"test_status"`

	TestTiming float64 `json:"test_timing"`

	JobType string `json:"job_type"`
}

// Repository is an github repository info managed by the storage.
type ProwJobStatus struct {
	JobID string `json:"job_id"`

	CreatedAt time.Time `json:"created_at"`

	State string `json:"state"`

	Duration float64 `json:"duration"`

	TestsCount int64 `json:"tests_count"`

	FailedCount int64 `json:"failed_count"`

	SkippedCount int64 `json:"skipped_count"`

	JobType string `json:"job_type"`

	JobName string `json:"job_name"`

	JobURL string `json:"job_url"`

	CIFailed int16 `json:"ci_failed"`
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

type ProwJobsMetrics struct {
	RepositoryName  string `json:"repository_name"`
	JobType         string `json:"type"`
	GitOrganization string `json:"git_organization"`
	Jobs            []Jobs `json:"jobs"`
}

type Jobs struct {
	Name    string    `json:"name"`
	Summary Summary   `json:"summary"`
	Metrics []Metrics `json:"metrics"`
}

type Metrics struct {
	SuccessRate  float64 `json:"success_rate"`
	FailureRate  float64 `json:"failure_rate"`
	CiFailedRate float64 `json:"ci_failed_rate"`
	Date         string  `json:"date"`
}

type Summary struct {
	DateFrom       string  `json:"date_from"`
	DateTo         string  `json:"date_to"`
	SuccessRateAvg float64 `json:"success_rate_avg"`
	JobFailedAvg   float64 `json:"failure_rate_avg"`
	CIFailedAvg    float64 `json:"ci_failed_rate_avg"`
	TotalJobs      int     `json:"total_jobs"`
}
