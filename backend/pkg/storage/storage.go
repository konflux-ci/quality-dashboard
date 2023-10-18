package storage

import (
	"errors"

	"github.com/andygrunwald/go-jira"
	"github.com/google/uuid"
	coverageV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/codecov/v1alpha1"
	failureV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/failure/v1alpha1"
	repoV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/github/v1alpha1"
	jiraV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/jira/v1alpha1"
	v1alphaPlugins "github.com/redhat-appstudio/quality-studio/api/apis/plugins/v1alpha1"
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
	TotalBugsResolutionTime(priority, startDate, endDate string, team *db.Teams) (bugsMetrics jiraV1Alpha1.ResolvedBugsMetrics, err error)
	GetRepository(repositoryName string, gitOrganizationName string) (*db.Repository, error)
	GetLatestProwTestExecution(r *db.Repository, jobType string) (*db.ProwJobs, error)
	GetSuitesByJobID(jobID string) ([]*db.ProwSuites, error)
	GetProwJobsResults(*db.Repository) ([]*db.ProwJobs, error)
	GetProwJobsResultsByJobID(jobID string) ([]*db.ProwJobs, error)
	GetMetrics(gitOrganization, repoName, jobType, startDate, endDate string) (prowV1Alpha1.JobsMetrics, error)
	GetAllTeamsFromDB() ([]*db.Teams, error)
	GetTeamByName(teamName string) (*db.Teams, error)
	ListWorkflowsByRepository(repositoryName string) (w []repoV1Alpha1.Workflow, err error)
	ListRepositories(team *db.Teams) ([]repoV1Alpha1.Repository, error)
	ListRepositoriesQualityInfo(team *db.Teams, startDate, endDate string) ([]RepositoryQualityInfo, error)
	GetAllJiraBugs() ([]*db.Bugs, error)
	GetPullRequestsByRepository(repositoryName, organization, startDate, endDate string) (repoV1Alpha1.PullRequestsInfo, error)
	GetFrequency(team *db.Teams, errorMessage, startDate, endDate string) (float64, error)
	GetJiraStatus(key string) (string, error)
	GetFailuresByDate(team *db.Teams, startDate, endDate string) ([]*failureV1Alpha1.Failure, error)
	GetAllFailures(team *db.Teams) ([]*db.Failure, error)
	BugExists(projectKey string, t *db.Teams) (bool, error)
	GetPluginByName(pluginName string) (*db.Plugins, error)
	GetPluginsByTeam(team *db.Teams) (storagePlugin []*v1alphaPlugins.Plugin, err error)

	// POST
	CreateRepository(p repoV1Alpha1.Repository, team_id uuid.UUID) (*db.Repository, error)
	CreateQualityStudioTeam(teamName string, description string, jira_keys string) (*db.Teams, error)
	CreateWorkflows(p repoV1Alpha1.Workflow, repo_id string) error
	CreateCoverage(p coverageV1Alpha1.Coverage, repo_id string) error
	CreateProwJobSuites(prowJobStatus prowV1Alpha1.JobSuites, repo_id string) error
	CreateProwJobResults(prowJobStatus prowV1Alpha1.Job, repo_id string) error
	ReCreateWorkflow(workflow repoV1Alpha1.Workflow, repoName string) error
	UpdateCoverage(codecov coverageV1Alpha1.Coverage, repoName string) error
	CreateJiraBug(bugsArr []jira.Issue, team *db.Teams) error
	UpdateTeam(t *db.Teams, target string) error
	GetOpenBugsMetricsByStatusAndPriority(priority, startDate, endDate string, team *db.Teams) (bugsMetrics jiraV1Alpha1.OpenBugsMetrics, err error)
	CreatePullRequests(prs repoV1Alpha1.PullRequests, repo_id string) error
	CreateFailure(f failureV1Alpha1.Failure, team_id uuid.UUID) error
	UpdateBuildLogErrors(jobID, buildErrorLogs string) error
	GetAllProwJobs(startDate, endDate string) ([]*db.ProwJobs, error)
	CreatePlugin(plugin *v1alphaPlugins.PluginSpec) (*db.Plugins, error)
	InstallPlugin(team *db.Teams, plugin *db.Plugins) (db *db.Teams, err error)

	// Delete
	DeleteRepository(repositoryName, gitOrganizationName string) error
	DeleteTeam(teamName string) (bool, error)
	DeleteJiraBugsByProject(projectKey string, team *db.Teams) error
	DeleteFailure(jiraKey string) error
	RemovePlugin(team *db.Teams, plugin *db.Plugins) (*db.Teams, error)

	ListPlugins() ([]*db.Plugins, error)
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
