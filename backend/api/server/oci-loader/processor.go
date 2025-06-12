package ociloader

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"time"

	prowV1Alpha1 "github.com/konflux-ci/quality-dashboard/api/apis/prow/v1alpha1"
	"github.com/konflux-ci/quality-dashboard/api/server/router/prow"
	"github.com/konflux-ci/quality-dashboard/pkg/storage"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db"
	"go.uber.org/zap"
)

// ProcessArtifactSet parses and persists data from a discovered artifact set.
// It reads the pipeline status JSON (required) and optionally parses the xUnit test report (if present).
func ProcessArtifactSet(set ArtifactSet, storage storage.Storage, logger *zap.Logger) error {
	jsonData, err := os.ReadFile(set.PipelineStatusPath)
	if err != nil {
		return fmt.Errorf("read pipeline status: %w", err)
	}
	var pipelineData PipelineStatus
	if err := json.Unmarshal(jsonData, &pipelineData); err != nil {
		return fmt.Errorf("unmarshal pipeline status: %w", err)
	}

	// Optionally parse xUnit test report if it exists
	var xunitData *prow.TestSuites
	if set.E2EReportPath != "" {
		xunitData, _ = ParseXUnitReport(set.E2EReportPath) // ignore parse errors silently
	}

	return saveResultsToDB(pipelineData, xunitData, storage, logger)
}

// saveResultsToDB persists pipeline metadata, Tekton tasks, and xUnit test results into the database.
// It performs input validation, associates results with a repository, and gracefully handles missing or malformed data.
// NOTE: OCI artifacts are only used if the CI system is Konflux.
func saveResultsToDB(p PipelineStatus, x *prow.TestSuites, storage storage.Storage, logger *zap.Logger) error {
	if p.PipelineRunName == "" || p.Status == "" || p.Scenario == "" || p.EventType == "" {
		logger.Sugar().Warnf("Skipping save: missing required pipeline fields for run '%s'", p.PipelineRunName)
		return nil
	}

	// Attempt to resolve the repository
	var repoID string
	if p.Git.Repository != "" && p.Git.Organization != "" {
		repo, err := storage.GetRepository(p.Git.Repository, p.Git.Organization)
		if err != nil {
			return fmt.Errorf("get repo: %w", err)
		}
		repoID = repo.ID
	} else {
		logger.Sugar().Warnf("Missing repository or organization info for PipelineRun '%s'", p.PipelineRunName)
		return nil
	}

	// Create and store the Prow job summary
	job, err := storage.CreateProwJobResults(prowV1Alpha1.Job{
		JobID:                 p.PipelineRunName,
		CreatedAt:             time.Now(),
		State:                 transformStatus(p.Status),
		JobType:               p.EventType,
		JobName:               p.Scenario,
		JobURL:                "none",
		ExternalServiceImpact: false,
	}, repoID)
	if err != nil {
		return fmt.Errorf("create job: %w", err)
	}

	// Validate and convert Tekton task runs before bulk insert
	validTasks := make([]*db.TektonTasks, 0, len(p.TaskRuns))
	for _, tr := range p.TaskRuns {
		if tr.Name == "" || tr.Status == "" || tr.Duration == "" {
			logger.Sugar().Warnf("Skipping invalid task run in '%s': %+v", p.PipelineRunName, tr)
			continue
		}
		validTasks = append(validTasks, &db.TektonTasks{
			TaskName:        tr.Name,
			Status:          transformStatus(tr.Status),
			DurationSeconds: tr.Duration,
		})
	}
	if len(validTasks) > 0 {
		if err := storage.CreateTektonTasksBulk(validTasks, &job.ID); err != nil {
			logger.Sugar().Warnf("create task runs: %w", err)
		}
	}

	// Skip test suite saving if xUnit report is not available
	if x == nil {
		return nil
	}

	re := regexp.MustCompile(`\[It]\s*\[([^]]+)]`)
	// Iterate through each test suite and its test cases
	for _, suite := range x.Suites {
		for _, tc := range suite.TestCases {
			if tc.FailureOutput == nil {
				continue
			}

			suiteName := suite.Name
			if suiteName == "" {
				if m := re.FindStringSubmatch(tc.Name); len(m) > 1 {
					suiteName = m[1]
				}
			}
			if suiteName == "" {
				logger.Sugar().Warnf("Skip test case: no suite name: %s", tc.Name)
				continue
			}

			js := prowV1Alpha1.JobSuites{
				JobID:          p.PipelineRunName,
				JobName:        p.Scenario,
				TestCaseName:   tc.Name,
				TestCaseStatus: tc.Status,
				TestTiming:     tc.Duration,
				JobType:        p.EventType,
				ErrorMessage:   tc.FailureOutput.Message,
				SuiteName:      suiteName,
				JobURL:         "tekton/konflux",
				CreatedAt:      time.Now(),
			}

			if err := storage.CreateProwJobSuites(js, repoID); err != nil {
				logger.Sugar().Errorf("store job suite %s: %v", tc.Name, err)
			}
		}
	}

	return nil
}

// transformStatus normalizes Tekton status values to match internal storage expectations like Prow.
func transformStatus(status string) string {
	switch status {
	case "Succeeded":
		return "success"
	case "Failed":
		return "failure"
	case "Cancelled":
		return "aborted"
	case "TaskRunCancelled":
		return "aborted"
	default:
		return status
	}
}
