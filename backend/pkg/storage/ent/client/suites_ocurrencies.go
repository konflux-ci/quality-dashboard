package client

import (
	"context"
	"fmt"
	"math"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/redhat-appstudio/quality-studio/api/apis/prow/v1alpha1"

	//"github.com/redhat-appstudio/quality-studio/pkg/ml"

	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/prowjobs"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/prowsuites"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/repository"
	util "github.com/redhat-appstudio/quality-studio/pkg/utils"
	"k8s.io/utils/strings/slices"
)

func (d *Database) GetSuitesFailureFrequency(gitOrg string, repoName string, jobName string, startDate string, endDate string) (*v1alpha1.FlakyFrequency, error) {
	flakyFrequency := new(v1alpha1.FlakyFrequency)
	var suitesFailure []struct {
		SuiteName string `json:"suite_name"`
		Status    string `json:"status"`
		Count     int    `json:"count"`
	}

	repository, err := d.client.Repository.Query().
		Where(repository.RepositoryName(repoName)).Where(repository.GitOrganization(gitOrg)).Only(context.TODO())
	if err != nil {
		return &v1alpha1.FlakyFrequency{}, convertDBError("get repository: %w", err)
	}

	err = d.client.Repository.QueryProwSuites(repository).
		Where(func(s *sql.Selector) { // "merged_at BETWEEN ? AND 2022-08-17", "2022-08-16"
			s.Where(sql.ExprP(fmt.Sprintf("created_at BETWEEN '%s' AND '%s'", startDate, endDate)))
		}).
		Where(prowsuites.JobName(jobName)).
		Where(prowsuites.ExternalServicesImpact(false)).
		GroupBy(prowsuites.FieldSuiteName, prowsuites.FieldStatus).
		Aggregate(db.Count()).
		Scan(context.Background(), &suitesFailure)

	if err != nil {
		return &v1alpha1.FlakyFrequency{}, convertDBError("get suites: %w", err)
	}

	allJobs, err := d.client.Repository.QueryProwJobs(repository).
		Where(prowjobs.JobName(jobName)).
		Where(func(s *sql.Selector) { // "merged_at BETWEEN ? AND 2022-08-17", "2022-08-16"
			s.Where(sql.ExprP(fmt.Sprintf("created_at BETWEEN '%s' AND '%s'", startDate, endDate)))
		}).
		Aggregate(
			db.Count(),
		).
		Int(context.Background())

	if err != nil {
		return &v1alpha1.FlakyFrequency{}, convertDBError("get all prow jobs: %w", err)
	}

	allImpacted, err := d.client.Repository.QueryProwSuites(repository).
		Where(func(s *sql.Selector) { // "merged_at BETWEEN ? AND 2022-08-17", "2022-08-16"
			s.Where(sql.ExprP(fmt.Sprintf("created_at BETWEEN '%s' AND '%s'", startDate, endDate)))
		}).
		Where(prowsuites.JobName(jobName)).
		Where(prowsuites.ExternalServicesImpact(false)).
		Aggregate(
			db.Count(),
		).
		All(context.Background())

	for _, suiteFail := range suitesFailure {
		testCase := make([]v1alpha1.TestCases, 0)

		suites, err := d.client.Repository.QueryProwSuites(repository).Where(func(s *sql.Selector) { // "merged_at BETWEEN ? AND 2022-08-17", "2022-08-16"
			s.Where(sql.ExprP(fmt.Sprintf("created_at BETWEEN '%s' AND '%s'", startDate, endDate)))
		}).Where(prowsuites.JobName(jobName)).
			Where(prowsuites.SuiteName(suiteFail.SuiteName)).
			All(context.Background())
		if err != nil {
			continue
		}

		for _, s := range suites {
			var msg = []v1alpha1.Messages{}
			testcase, err := d.client.Repository.QueryProwSuites(repository).Where(func(s *sql.Selector) { // "merged_at BETWEEN ? AND 2022-08-17", "2022-08-16"
				s.Where(sql.ExprP(fmt.Sprintf("created_at BETWEEN '%s' AND '%s'", startDate, endDate)))
			}).Where(prowsuites.Name(s.Name)).
				Where(prowsuites.JobName(jobName)).
				Where(prowsuites.ExternalServicesImpact(false)).
				Where(prowsuites.SuiteName(suiteFail.SuiteName)).All(context.Background())

			if err != nil {
				return &v1alpha1.FlakyFrequency{}, convertDBError("get impacted jobs: %w", err)
			}

			for _, tc := range testcase {
				msg = append(msg, v1alpha1.Messages{
					JobId:       tc.JobID,
					JobURL:      tc.JobURL,
					FailureDate: tc.CreatedAt,
					Message:     *tc.ErrorMessage,
				})
			}

			if !AlreadyExists(testCase, s.Name) {
				testFlakyAvg := util.RoundTo((float64(getLengthOfJobIdsInPRowSuiteWithoutDuplication(testcase))/float64(len(allImpacted)))*100, 2)

				if math.IsNaN(testFlakyAvg) || math.IsInf(testFlakyAvg, len(allImpacted)) {
					testFlakyAvg = 0
				}

				testCase = append(testCase, v1alpha1.TestCases{
					Name:           s.Name,
					Count:          len(testcase),
					TestCaseImpact: testFlakyAvg,
					Messages:       msg,
				})
			}
		}

		flakySuiteAVG := util.RoundTo((float64(suiteFail.Count)/float64(len(allImpacted)))*100, 2)

		if math.IsNaN(flakySuiteAVG) || math.IsInf(flakySuiteAVG, len(allImpacted)) {
			flakySuiteAVG = 0
		}

		flakyFrequency.SuitesFailureFrequency = append(flakyFrequency.SuitesFailureFrequency, v1alpha1.SuitesFailureFrequency{
			SuiteName:     suiteFail.SuiteName,
			Status:        suiteFail.Status,
			AverageImpact: flakySuiteAVG,
			TestCases:     testCase,
		})

	}

	if err != nil {
		return &v1alpha1.FlakyFrequency{}, convertDBError("get impacted jobs: %w", err)
	}

	globalFlakyAvg := (float64(getLengthOfJobIdsInPRowSuiteWithoutDuplication(allImpacted)) / float64(allJobs)) * 100

	if math.IsNaN(globalFlakyAvg) || math.IsInf(globalFlakyAvg, allJobs) {
		globalFlakyAvg = 0
	}

	flakyFrequency.GlobalImpact = util.RoundTo(globalFlakyAvg, 2)
	flakyFrequency.JobName = jobName
	flakyFrequency.GitOrganization = gitOrg
	flakyFrequency.JobsExecuted = allJobs
	flakyFrequency.JobsAffectedByFlayTests = getLengthOfJobIdsInPRowSuiteWithoutDuplication(allImpacted)
	flakyFrequency.RepositoryName = repoName

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

	p, err := d.client.Repository.QueryProwSuites(repo).Select(prowsuites.FieldJobName).Unique(true).All(context.Background())
	if err != nil {
		return nil, convertDBError("get jobs name: %w", err)
	}

	for _, v := range p {
		if !slices.Contains(flakyJobs, v.JobName) {
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
