package storage

import (
	"context"
	"errors"

	"github.com/andygrunwald/go-jira"
	"github.com/google/uuid"
	coverageV1Alpha1 "github.com/konflux-ci/quality-dashboard/api/apis/codecov/v1alpha1"
	configurationV1Alpha1 "github.com/konflux-ci/quality-dashboard/api/apis/configuration/v1alpha1"
	failureV1Alpha1 "github.com/konflux-ci/quality-dashboard/api/apis/failure/v1alpha1"
	repoV1Alpha1 "github.com/konflux-ci/quality-dashboard/api/apis/github/v1alpha1"
	jiraV1Alpha1 "github.com/konflux-ci/quality-dashboard/api/apis/jira/v1alpha1"
	"github.com/konflux-ci/quality-dashboard/api/apis/prow/v1alpha1"
	prowV1Alpha1 "github.com/konflux-ci/quality-dashboard/api/apis/prow/v1alpha1"
	userV1Alpha1 "github.com/konflux-ci/quality-dashboard/api/apis/user/v1alpha1"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db"
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
	TotalBugsResolutionTime(priority, startDate, endDate string, team *db.Teams) (bugsMetrics jiraV1Alpha1.ResolvedBugsMetrics, err error)
	GetRepository(repositoryName string, gitOrganizationName string) (*db.Repository, error)
	GetLatestProwTestExecution(r *db.Repository, jobType string) (*db.ProwJobs, error)
	GetSuitesByJobID(jobID string) ([]*db.ProwSuites, error)
	GetProwJobsResults(repo *db.Repository, startDate, endDate string) ([]*db.ProwJobs, error)
	GetProwJobsResultsByJobID(jobID string) ([]*db.ProwJobs, error)
	ObtainProwMetricsByJob(gitOrganization string, repositoryName string, jobName string, startDate string, endDate string) (*prowV1Alpha1.JobsMetrics, error)
	GetAllTeamsFromDB() ([]*db.Teams, error)
	GetTeamByName(teamName string) (*db.Teams, error)
	ListWorkflowsByRepository(repositoryName string) (w []repoV1Alpha1.Workflow, err error)
	ListRepositories(team *db.Teams) ([]repoV1Alpha1.Repository, error)
	ListRepositoriesQualityInfo(team *db.Teams, startDate, endDate string) ([]RepositoryQualityInfo, error)
	GetAllJiraBugs() ([]*db.Bugs, error)
	GetAllJiraBugsByProject(project string) ([]*db.Bugs, error)
	GetJobsNameAndType(repo *db.Repository) ([]*db.ProwJobs, error)
	GetMetricsSummaryByDay(repo *db.Repository, job, startDate, endDate string) []*prowV1Alpha1.JobsMetrics
	GetAllOCIArtifacts() ([]*db.OCI, error)

	// GetAllOpenBugs(dateFrom, dateTo string) ([]*db.Bugs, error)
	GetAllOpenBugs(project string) ([]*db.Bugs, error)
	GetAllOpenBugsForSliAlerts(project string) ([]*db.Bugs, error)
	GetPullRequestsByRepository(repositoryName, organization, startDate, endDate string) (repoV1Alpha1.PullRequestsInfo, error)
	GetFrequency(team *db.Teams, errorMessage, startDate, endDate string) (float64, error)
	GetJiraBug(key string) (*db.Bugs, error)
	GetJiraStatus(key string) (string, error)
	GetFailuresByDate(team *db.Teams, startDate, endDate string) ([]*failureV1Alpha1.Failure, error)
	GetAllFailures(team *db.Teams) ([]*db.Failure, error)
	ListAllRepositories() ([]*db.Repository, error)
	BugExists(projectKey string, t *db.Teams) (bool, error)
	GetProwFlakyTrendsMetrics(gitOrg string, repoName string, jobName string, startDate string, endDate string) []prowV1Alpha1.FlakyMetrics
	GetProwJobsByRepoOrg(repo *db.Repository) ([]string, error)
	GetUser(userEmail string) (*db.Users, error)
	ListAllUsers() ([]*db.Users, error)
	GetTasksMetrics(startDate string, endDate string) ([]TaskMetrics, error)
	GetSuitesFailureFrequency(gitOrg string, repoName string, jobName string, startDate string, endDate string) (*v1alpha1.FlakyFrequency, error)
	// POST
	CreateRepository(p repoV1Alpha1.Repository, team_id uuid.UUID) (*db.Repository, error)
	CreateQualityStudioTeam(teamName string, description string, jira_keys string) (*db.Teams, error)
	CreateWorkflows(p repoV1Alpha1.Workflow, repo_id string) error
	CreateCoverage(p coverageV1Alpha1.Coverage, repo_id string) error
	CreateProwJobSuites(prowJobStatus prowV1Alpha1.JobSuites, repo_id string) error
	CreateProwJobResults(job prowV1Alpha1.Job, repo_id string) (*db.ProwJobs, error)
	ReCreateWorkflow(workflow repoV1Alpha1.Workflow, repoName string) error
	UpdateCoverage(codecov coverageV1Alpha1.Coverage, repoName string) error
	CreateJiraBug(bugsArr []jira.Issue, team *db.Teams) error
	UpdateTeam(t *db.Teams, target string) error
	GetOpenBugsMetricsByStatusAndPriority(priority, startDate, endDate string, team *db.Teams) (bugsMetrics jiraV1Alpha1.OpenBugsMetrics, err error)
	CreatePullRequests(prs repoV1Alpha1.PullRequests, repo_id string) error
	CreateFailure(f failureV1Alpha1.Failure, team_id uuid.UUID) error
	UpdateErrorMessages(jobID, buildErrorLogs, e2eErrorMessages string) error
	GetAllProwJobs(startDate, endDate string) ([]*db.ProwJobs, error)
	ListFailedProwJobsByRepository(repo *db.Repository) ([]*db.ProwJobs, error)
	CreateUser(u userV1Alpha1.User) error
	GetOpenBugsAffectingCI(t *db.Teams, startDate, endDate string) ([]*db.Bugs, error)
	CreateConfiguration(c configurationV1Alpha1.Configuration) error
	GetConfiguration(teamName string) (*db.Configuration, error)
	CreateOCIArtifact(ociArtifact *db.OCI, repoID *string) (*db.OCI, error)
	CreateTektonTasksBulk(tasks []*db.TektonTasks, jobID *int) error

	// Update
	UpdateOCIArtifact(ctx context.Context, ociArtifact *db.OCI) (*db.OCI, error)

	// Delete
	DeleteRepository(repositoryName, gitOrganizationName string) error
	DeleteTeam(teamName string) (bool, error)
	DeleteJiraBugsByProject(projectKey string, team *db.Teams) error
	DeleteJiraBugByJiraKey(jiraKey string) error
	DeleteFailure(teamID, failureID uuid.UUID) error
	DeleteAllJiraBugByTeam(teamName string) error
}

// TaskMetrics holds the calculated metrics for a single Tekton task.
type TaskMetrics struct {
	TaskName          string  `json:"task_name"`
	TotalRuns         int     `json:"total_runs"`
	SuccessPercentage float64 `json:"success_percentage"`
	FailurePercentage float64 `json:"failure_percentage"`
	AverageDuration   float64 `json:"average_duration"`
}

type RepositoryQualityInfo struct {
	RepositoryName string `json:"repository_name"`

	GitOrganization string `json:"git_organization"`

	Description string `json:"description"`

	GitURL string `json:"git_url"`

	CodeCoverage coverageV1Alpha1.Coverage `json:"code_coverage"`

	PullRequests repoV1Alpha1.PullRequestsInfo `json:"prs"`

	Workflows []repoV1Alpha1.Workflow `json:"workflows"`
}
