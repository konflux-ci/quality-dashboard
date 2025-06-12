// Package ociloader defines the core data structures used throughout the artifact loading
// and processing pipeline, including Tekton pipeline metadata, configuration, and grouped workloads.
package ociloader

import (
	"time"

	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db"
)

// PipelineStatus represents the metadata and result summary for a single Tekton PipelineRun.
// It is extracted from the `pipeline-status.json` artifact generated during pipeline execution.
type PipelineStatus struct {
	// Name of the Tekton PipelineRun
	PipelineRunName string `json:"pipelineRunName"`

	// Total execution time (human-readable)
	Duration string `json:"duration"`

	// Final status (e.g. "Succeeded", "Failed")
	Status string `json:"status"`

	// Source event type (e.g. Push, Pull Request)
	EventType string `json:"eventType"`

	// Test scenario name (e.g. specific E2E suite)
	Scenario string `json:"scenario"`

	// Associated Git repository information
	Git GitInfo `json:"git"`

	// Breakdown of each TaskRun executed in the pipeline
	TaskRuns []TaskRun `json:"taskRuns"`
}

// GitInfo captures repository metadata for the source of a Tekton PipelineRun.
// It is part of the `pipeline-status.json` structure.
type GitInfo struct {
	// GitHub organization
	Organization string `json:"gitOrganization,omitempty"`

	// GitHub repository name (optional)
	Repository string `json:"gitRepository,omitempty"`
}

// TaskRun represents the status of a single Tekton TaskRun within a PipelineRun.
type TaskRun struct {
	// Name of the TaskRun
	Name string `json:"name"`

	// Final status of the TaskRun
	Status string `json:"status"`

	// Duration of TaskRun execution
	Duration string `json:"duration"`
}

// ArtifactSet defines the presence of required and optional artifacts in a discovered directory.
// This structure is used after scanning the filesystem for Tekton result sets.
type ArtifactSet struct {
	// Filesystem path to the directory containing artifacts
	Directory string

	// Path to the required `pipeline-status.json` file
	PipelineStatusPath string

	// Optional path to `e2e-report.xml` (may be empty)
	E2EReportPath string
}

// LoaderConfig defines runtime parameters used by the OCI artifact loader.
// These are typically passed in during instantiation of the Loader.
type LoaderConfig struct {
	// Number of concurrent worker goroutines
	NumWorkers int

	// Minimum age of artifacts before reprocessing is allowed
	ProcessingThreshold time.Duration

	// Root directory for storing unpacked artifacts
	ArtifactsDir string

	// Directory for caching OCI downloads
	CacheDir string

	// Whether to remove temporary directories after processing
	CleanupOnCompletion bool
}

// repoJob groups a collection of artifacts by their originating repository,
// forming a unit of work to be processed by a loader worker.
type repoJob struct {
	// Pointer to the repository entity
	Repository *db.Repository

	// List of associated OCI artifact metadata
	Artifacts []*db.OCI
}
