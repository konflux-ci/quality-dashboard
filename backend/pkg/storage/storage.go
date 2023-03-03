package storage

import (
	"errors"

	"github.com/google/uuid"
	coverageV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/codecov/v1alpha1"
	repoV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/github/v1alpha1"
	jiraV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/jira/v1alpha1"
	prowV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/prow/v1alpha1"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
)

var (
	// ErrNotFound is the error returned by storages if a resource cannot be found.
	ErrNotFound = errors.New("not found")

	// ErrAlreadyExists is the error returned by storages if a resource ID is taken during a create.
	ErrAlreadyExists = errors.New("already exists")
)

// Storage is the storage interface used by the server. Implementations are
// required to be able to perform atomic compare-and-swap updates and either
// support timezones or standardize on UTC.
type Storage interface {
	Close() error

	// GET
	TotalBugsResolutionTime(priority string) (bugsMetrics jiraV1Alpha1.ResolvedBugsMetrics, err error)
	GetRepository(repositoryName string, gitOrganizationName string) (*db.Repository, error)
	GetLatestProwTestExecution(r *db.Repository, jobType string) (*db.ProwJobs, error)
	GetSuitesByJobID(jobID string) ([]*db.ProwSuites, error)
	GetProwJobsResults(*db.Repository) ([]*db.ProwSuites, error)
	GetProwJobsResultsByJobID(jobID string) ([]*db.ProwJobs, error)
	GetMetrics(gitOrganization string, repoName string, jobType string, startDate string, endDate string) prowV1Alpha1.JobsMetrics
	GetAllTeamsFromDB() ([]*db.Teams, error)
	GetTeamByName(teamName string) (*db.Teams, error)
	ListWorkflowsByRepository(repositoryName string) (w []repoV1Alpha1.Workflow, err error)
	ListRepositories(team *db.Teams) ([]repoV1Alpha1.Repository, error)
	ListRepositoriesQualityInfo(team *db.Teams) ([]RepositoryQualityInfo, error)
	GetAllJiraBugs() ([]*db.Bugs, error)

	// POST
	CreateRepository(p repoV1Alpha1.Repository, team_id uuid.UUID) (*db.Repository, error)
	CreateQualityStudioTeam(teamName string, description string) (*db.Teams, error)
	CreateWorkflows(p repoV1Alpha1.Workflow, repo_id uuid.UUID) error
	CreateCoverage(p coverageV1Alpha1.Coverage, repo_id uuid.UUID) error
	CreateProwJobSuites(prowJobStatus prowV1Alpha1.JobSuites, repo_id uuid.UUID) error
	CreateProwJobResults(prowJobStatus prowV1Alpha1.Job, repo_id uuid.UUID) error
	ReCreateWorkflow(workflow repoV1Alpha1.Workflow, repoName string) error
	UpdateCoverage(codecov coverageV1Alpha1.Coverage, repoName string) error
	CreateJiraBug(bug jiraV1Alpha1.JiraBug) error
	UpdateTeam(t *db.Teams, target string) error
	DeleteTeam(teamName string) (bool, error)
	GetOpenBugsMetricsByStatusAndPriority(priority string) (bugsMetrics jiraV1Alpha1.OpenBugsMetrics, err error)

	// Delete
	DeleteRepository(repositoryName string, gitOrganizationName string) error
}

type RepositoryQualityInfo struct {
	RepositoryName string `json:"repository_name"`

	GitOrganization string `json:"git_organization"`

	Description string `json:"description"`

	GitURL string `json:"git_url"`

	CodeCoverage coverageV1Alpha1.Coverage `json:"code_coverage"`
}
