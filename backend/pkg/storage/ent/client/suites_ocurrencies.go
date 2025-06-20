package client

import (
	"context"
	"fmt"
	"math"
	"slices"
	"strings"
	"time"

	"github.com/konflux-ci/quality-dashboard/api/apis/prow/v1alpha1"
	"github.com/konflux-ci/quality-dashboard/api/server/router/prow"

	//"github.com/konflux-ci/quality-dashboard/pkg/ml"

	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db/prowjobs"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db/prowsuites"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db/repository"
	util "github.com/konflux-ci/quality-dashboard/pkg/utils"
)

func parseDateRange(startDateStr, endDateStr string) (time.Time, time.Time, error) {
	parsedStartDate, err := time.Parse(time.RFC3339Nano, startDateStr)
	if err != nil {
		// Fallback for 2006-01-02 15:04:05 format if that's also used in other parts of your code
		parsedStartDate, err = time.Parse("2006-01-02 15:04:05", startDateStr)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("failed to parse start date '%s': %w", startDateStr, err)
		}
	}
	parsedEndDate, err := time.Parse(time.RFC3339Nano, endDateStr)
	if err != nil {
		// Fallback
		parsedEndDate, err = time.Parse("2006-01-02 15:04:05", endDateStr)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("failed to parse end date '%s': %w", endDateStr, err)
		}
	}
	return parsedStartDate, parsedEndDate, nil
}

func (d *Database) GetSuitesFailureFrequency(gitOrg string, repoName string, jobName string, startDate string, endDate string) (*v1alpha1.FlakyFrequency, error) {
	flakyFrequency := new(v1alpha1.FlakyFrequency)
	flakyFrequency.JobName = jobName
	flakyFrequency.GitOrganization = gitOrg
	flakyFrequency.RepositoryName = repoName

	parsedStartDate, parsedEndDate, err := parseDateRange(startDate, endDate)
	if err != nil {
		return &v1alpha1.FlakyFrequency{}, fmt.Errorf("failed to parse date range: %w", err)
	}

	repo, err := d.client.Repository.Query().
		Where(repository.RepositoryName(repoName)).
		Where(repository.GitOrganization(gitOrg)).
		Only(context.Background())
	if err != nil {
		return &v1alpha1.FlakyFrequency{}, convertDBError("get repository: %w", err)
	}

	allRelevantProwSuites, err := d.client.Repository.QueryProwSuites(repo).
		Where(prowsuites.JobName(jobName)).
		Where(prowsuites.ExternalServicesImpact(false)).
		Where(prowsuites.CreatedAtGTE(parsedStartDate)).
		Where(prowsuites.CreatedAtLTE(parsedEndDate)).
		All(context.Background())
	if err != nil {
		return &v1alpha1.FlakyFrequency{}, convertDBError("get all relevant prow suites: %w", err)
	}

	// Fetch job IDs and deduplicate them manually
	jobIDs, err := d.client.Repository.QueryProwJobs(repo).
		Where(prowjobs.JobName(jobName)).
		Where(prowjobs.Not(prowjobs.State(string(prow.AbortedState)))).
		Where(prowjobs.CreatedAtGTE(parsedStartDate)).
		Where(prowjobs.CreatedAtLTE(parsedEndDate)).
		Select(prowjobs.FieldJobID).
		Strings(context.Background())
	if err != nil {
		return &v1alpha1.FlakyFrequency{}, convertDBError("get prow job IDs: %w", err)
	}
	uniqueJobIDs := make(map[string]struct{})
	for _, id := range jobIDs {
		uniqueJobIDs[id] = struct{}{}
	}
	flakyFrequency.JobsExecuted = len(uniqueJobIDs)

	suiteData := make(map[string]map[string]int)
	testCaseData := make(map[string]map[string][]v1alpha1.Messages)
	suiteJobIDs := make(map[string]map[string]struct{})
	allImpactedJobIDs := make(map[string]struct{})

	for _, s := range allRelevantProwSuites {
		if _, ok := suiteData[s.SuiteName]; !ok {
			suiteData[s.SuiteName] = make(map[string]int)
		}
		suiteData[s.SuiteName][s.Status]++

		if _, ok := testCaseData[s.SuiteName]; !ok {
			testCaseData[s.SuiteName] = make(map[string][]v1alpha1.Messages)
		}
		if s.ErrorMessage != nil {
			testCaseData[s.SuiteName][s.Name] = append(testCaseData[s.SuiteName][s.Name], v1alpha1.Messages{
				JobId:       s.JobID,
				JobURL:      s.JobURL,
				FailureDate: s.CreatedAt,
				Message:     *s.ErrorMessage,
			})
		}

		if _, ok := suiteJobIDs[s.SuiteName]; !ok {
			suiteJobIDs[s.SuiteName] = make(map[string]struct{})
		}
		suiteJobIDs[s.SuiteName][s.JobID] = struct{}{}
		allImpactedJobIDs[s.JobID] = struct{}{}
	}

	flakyFrequency.JobsAffectedByFlayTests = len(allImpactedJobIDs)

	globalFlakyAvg := (float64(flakyFrequency.JobsAffectedByFlayTests) / float64(flakyFrequency.JobsExecuted)) * 100
	if math.IsNaN(globalFlakyAvg) || flakyFrequency.JobsExecuted == 0 {
		globalFlakyAvg = 0
	}
	flakyFrequency.GlobalImpact = util.RoundTo(globalFlakyAvg, 2)

	for suiteName, statuses := range suiteData {
		totalSuiteRuns := 0
		for _, count := range statuses {
			totalSuiteRuns += count
		}

		var primaryStatus string
		if failedCount, ok := statuses[string(prow.FailureState)]; ok && failedCount > 0 {
			primaryStatus = string(prow.FailureState)
		} else if successCount, ok := statuses[string(prow.SuccessState)]; ok && successCount > 0 {
			primaryStatus = string(prow.SuccessState)
		} else {
			for status := range statuses {
				primaryStatus = status
				break
			}
		}

		flakySuiteAVG := util.RoundTo((float64(totalSuiteRuns)/float64(len(allRelevantProwSuites)))*100, 2)
		if math.IsNaN(flakySuiteAVG) || len(allRelevantProwSuites) == 0 {
			flakySuiteAVG = 0
		}

		flakySuite := v1alpha1.SuitesFailureFrequency{
			SuiteName:     suiteName,
			Status:        primaryStatus,
			AverageImpact: flakySuiteAVG,
			TestCases:     make([]v1alpha1.TestCases, 0),
		}

		if tests, ok := testCaseData[suiteName]; ok {
			for testCaseName, messages := range tests {
				uniqueJobIDsForTestCase := make(map[string]struct{})
				for _, msg := range messages {
					uniqueJobIDsForTestCase[msg.JobId] = struct{}{}
				}
				testCaseUniqueJobCount := len(uniqueJobIDsForTestCase)

				testCaseImpact := util.RoundTo((float64(testCaseUniqueJobCount)/float64(len(allImpactedJobIDs)))*100, 2)
				if math.IsNaN(testCaseImpact) || len(allImpactedJobIDs) == 0 {
					testCaseImpact = 0
				}

				flakySuite.TestCases = append(flakySuite.TestCases, v1alpha1.TestCases{
					Name:           testCaseName,
					Count:          len(messages),
					TestCaseImpact: testCaseImpact,
					Messages:       messages,
				})
			}
		}
		flakyFrequency.SuitesFailureFrequency = append(flakyFrequency.SuitesFailureFrequency, flakySuite)
	}

	return flakyFrequency, nil
}

func (d *Database) GetProwFlakyTrendsMetrics(gitOrg string, repoName string, jobName string, startDate string, endDate string) []v1alpha1.FlakyMetrics {
	var metrics []v1alpha1.FlakyMetrics

	dayArr := getRangesInISO(startDate, endDate)
	// range between one day (same day)
	if len(dayArr) == 2 && isSameDay(startDate, endDate) {
		metric, _ := d.GetSuitesFailureFrequency(gitOrg, repoName, jobName, startDate, endDate)
		metrics = append(metrics, v1alpha1.FlakyMetrics{
			JobsExecuted: metric.JobsExecuted,
			GlobalImpact: metric.GlobalImpact,
			Date:         startDate,
		})
		return metrics
	}

	// range between more than one day
	for i, day := range dayArr {
		t, _ := time.Parse(time.RFC3339, day)
		y, m, dd := t.Date()

		if i == 0 { // first day
			metric, _ := d.GetSuitesFailureFrequency(gitOrg, repoName, jobName, day, fmt.Sprintf("%04d-%02d-%02d 23:59:59", y, m, dd))
			metrics = append(metrics, v1alpha1.FlakyMetrics{
				GlobalImpact: metric.GlobalImpact,
				JobsExecuted: metric.JobsExecuted,
				Date:         fmt.Sprintf("%04d-%02d-%02d 23:59:59", y, m, dd),
			})
		} else {
			if i == len(dayArr)-1 { // last day
				metric, _ := d.GetSuitesFailureFrequency(gitOrg, repoName, jobName, fmt.Sprintf("%04d-%02d-%02d 00:00:00", y, m, dd), fmt.Sprintf("%04d-%02d-%02d 23:59:59", y, m, dd))
				metrics = append(metrics, v1alpha1.FlakyMetrics{
					JobsExecuted: metric.JobsExecuted,
					GlobalImpact: metric.GlobalImpact,
					Date:         fmt.Sprintf("%04d-%02d-%02d 23:59:59", y, m, dd),
				})
			} else { // middle days
				metric, _ := d.GetSuitesFailureFrequency(gitOrg, repoName, jobName, fmt.Sprintf("%04d-%02d-%02d 00:00:00", y, m, dd), fmt.Sprintf("%04d-%02d-%02d 23:59:59", y, m, dd))
				metrics = append(metrics, v1alpha1.FlakyMetrics{
					JobsExecuted: metric.JobsExecuted,
					GlobalImpact: metric.GlobalImpact,
					Date:         fmt.Sprintf("%04d-%02d-%02d 23:59:59", y, m, dd),
				})
			}
		}
	}

	return metrics
}

func (d *Database) GetProwJobsByRepoOrg(repo *db.Repository) ([]string, error) {
	var flakyJobs []string

	p, err := d.client.Repository.QueryProwJobs(repo).Select(prowjobs.FieldJobName).Unique(true).All(context.Background())
	if err != nil {
		return nil, convertDBError("get jobs name: %w", err)
	}

	for _, v := range p {
		if !slices.Contains(flakyJobs, v.JobName) && !strings.Contains(v.JobName, "main-images") {
			flakyJobs = append(flakyJobs, v.JobName)
		}
	}
	return flakyJobs, nil
}

func AlreadyExists(tc []v1alpha1.TestCases, caseName string) bool {
	for _, t := range tc {
		if t.Name == caseName {
			return true
		}
	}

	return false
}

func getLengthOfJobIdsInPRowSuiteWithoutDuplication(suites []*db.ProwSuites) int {
	var jobIds []string

	for _, s := range suites {
		jobIds = append(jobIds, s.JobID)
	}

	return len(removeDuplicateStr(jobIds))
}
