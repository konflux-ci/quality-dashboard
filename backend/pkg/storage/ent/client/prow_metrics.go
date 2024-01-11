package client

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	prowV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/prow/v1alpha1"
	"github.com/redhat-appstudio/quality-studio/api/server/router/prow"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/prowjobs"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/repository"
)

func (d *Database) ObtainProwMetricsByJob(gitOrganization string, repositoryName string, jobName string, startDate string, endDate string) *prowV1Alpha1.JobsMetrics {
	repo, err := d.client.Repository.Query().Where(repository.GitOrganization(gitOrganization)).Where(repository.RepositoryName(repositoryName)).First(context.Background())
	if err != nil {
		return nil
	}

	dbJobs, err := d.client.Repository.QueryProwJobs(repo).
		Where(prowjobs.JobName(jobName)).
		Where(func(s *sql.Selector) { // "created_at BETWEEN ? AND 2022-08-17", "2022-08-16"
			s.Where(sql.ExprP(fmt.Sprintf("created_at BETWEEN '%s' AND '%s'", startDate, endDate)))
		}).
		All(context.Background())
	if err != nil {
		return nil
	}

	numberOfSuccessJobs, _ := d.GetNumberOfSuccessJobs(repo, jobName, startDate, endDate)

	numberOfFailedJobs, _ := d.GetNumberOfFailedJobs(repo, jobName, startDate, endDate)
	numberOfInfraImpactJobs, _ := d.GetNumberOfInfraImpact(repo, jobName, startDate, endDate)
	totalImpact := numberOfFailedJobs + numberOfInfraImpactJobs

	flaky, _ := d.GetSuitesFailureFrequency(gitOrganization, repositoryName, jobName, startDate, endDate)

	extImpct, _ := d.GetJobsImpactedByAnExternalService(repo, jobName, startDate, endDate)

	unknownJobs := totalImpact - (flaky.JobsAffectedByFlayTests + extImpct + numberOfInfraImpactJobs)

	return &prowV1Alpha1.JobsMetrics{
		GitOrganization: gitOrganization,
		RepositoryName:  repositoryName,
		JobName:         jobName,
		StartDate:       startDate,
		EndDate:         endDate,
		JobsRuns: prowV1Alpha1.JobsRuns{
			Total:             len(dbJobs),
			Success:           numberOfSuccessJobs,
			Failures:          totalImpact,
			SuccessPercentage: CalculatePercentage(float64(numberOfSuccessJobs), float64(len(dbJobs))),
			FailedPercentage:  CalculatePercentage(float64(totalImpact), float64(len(dbJobs))),
		},
		JobsImpacts: prowV1Alpha1.JobsImpacts{
			InfrastructureImpact: prowV1Alpha1.InfrastructureImpact{
				Total:      numberOfInfraImpactJobs,
				Percentage: CalculatePercentage(float64(numberOfInfraImpactJobs), float64(len(dbJobs))),
			},
			FlakyTestsImpact: prowV1Alpha1.FlakyTestsImpact{
				Percentage: flaky.GlobalImpact,
				Total:      flaky.JobsAffectedByFlayTests,
			},
			ExternalServicesImpact: prowV1Alpha1.ExternalServicesImpact{
				Percentage: CalculatePercentage(float64(extImpct), float64(len(dbJobs))),
				Total:      extImpct,
			},
			UnknowFailuresImpact: prowV1Alpha1.UnknowFailuresImpact{
				Total:      unknownJobs,
				Percentage: CalculatePercentage(float64(unknownJobs), float64(len(dbJobs))),
			},
		},
	}
}

func (d *Database) GetMetricsSummaryByDay(repo *db.Repository, job, startDate, endDate string) []*prowV1Alpha1.JobsMetrics {
	var metrics []*prowV1Alpha1.JobsMetrics
	dayArr := getDatesBetweenRange(startDate, endDate)

	// range between one day (same day)
	if len(dayArr) == 2 && isSameDay(startDate, endDate) {
		metric := d.ObtainProwMetricsByJob(repo.GitOrganization, repo.RepositoryName, job, startDate, endDate)
		metrics = append(metrics, metric)
		return metrics
	}

	// range between more than one day
	for i, day := range dayArr {
		t, _ := time.Parse("2006-01-02 15:04:05", day)
		y, m, dd := t.Date()

		if i == 0 { // first day
			metric := d.ObtainProwMetricsByJob(repo.GitOrganization, repo.RepositoryName, job, day, fmt.Sprintf("%04d-%02d-%02d 23:59:59", y, m, dd))
			metrics = append(metrics, metric)
		} else {
			if i == len(dayArr)-1 { // last day
				metric := d.ObtainProwMetricsByJob(repo.GitOrganization, repo.RepositoryName, job, fmt.Sprintf("%04d-%02d-%02d 00:00:00", y, m, dd), day)
				metrics = append(metrics, metric)
			} else { // middle days
				metric := d.ObtainProwMetricsByJob(repo.GitOrganization, repo.RepositoryName, job, fmt.Sprintf("%04d-%02d-%02d 00:00:00", y, m, dd), fmt.Sprintf("%04d-%02d-%02d 23:59:59", y, m, dd))
				metrics = append(metrics, metric)
			}
		}
	}
	return metrics
}

func (d *Database) GetJobsNameAndType(repo *db.Repository) ([]*db.ProwJobs, error) {
	dbJobs, err := d.client.Repository.QueryProwJobs(repo).
		Select(prowjobs.FieldJobName).
		Select(prowjobs.FieldJobType).
		All(context.Background())

	if err != nil {
		return nil, err
	}

	// Create a map to store unique values of the jobNames
	uniqueProwJobs := make(map[string]*db.ProwJobs)

	// Iterate through the original array
	for _, s := range dbJobs {
		if _, found := uniqueProwJobs[s.JobName]; !found {
			uniqueProwJobs[s.JobName] = s
		}
	}

	// Create a new array to store unique structs
	var prowJobs []*db.ProwJobs

	// Iterate through the map and append unique structs to the new array
	for _, value := range uniqueProwJobs {
		prowJobs = append(prowJobs, value)
	}

	return prowJobs, nil
}

func (d *Database) GetNumberOfSuccessJobs(repo *db.Repository, jobName string, startDate string, endDate string) (totalSuccess int, err error) {
	return d.client.Repository.QueryProwJobs(repo).
		Where(prowjobs.JobName(jobName)).
		Where(prowjobs.State(string(prow.SuccessState))).
		Where(func(s *sql.Selector) { // "merged_at BETWEEN ? AND 2022-08-17", "2022-08-16"
			s.Where(sql.ExprP(fmt.Sprintf("created_at BETWEEN '%s' AND '%s'", startDate, endDate)))
		}).
		Aggregate(
			db.Count(),
		).
		Int(context.Background())
}

func (d *Database) GetNumberOfFailedJobs(repo *db.Repository, jobName string, startDate string, endDate string) (totalSuccess int, err error) {
	return d.client.Repository.QueryProwJobs(repo).
		Where(prowjobs.JobName(jobName)).
		Where(prowjobs.State(string(prow.FailureState))).
		Where(func(s *sql.Selector) { // "merged_at BETWEEN ? AND 2022-08-17", "2022-08-16"
			s.Where(sql.ExprP(fmt.Sprintf("created_at BETWEEN '%s' AND '%s'", startDate, endDate)))
		}).
		Aggregate(
			db.Count(),
		).
		Int(context.Background())
}

func (d *Database) GetNumberOfInfraImpact(repo *db.Repository, jobName string, startDate string, endDate string) (totalSuccess int, err error) {
	return d.client.Repository.QueryProwJobs(repo).
		Where(prowjobs.JobName(jobName)).
		Where(prowjobs.State(string(prow.ErrorState))).
		Where(func(s *sql.Selector) { // "merged_at BETWEEN ? AND 2022-08-17", "2022-08-16"
			s.Where(sql.ExprP(fmt.Sprintf("created_at BETWEEN '%s' AND '%s'", startDate, endDate)))
		}).
		Aggregate(
			db.Count(),
		).
		Int(context.Background())
}

func (d *Database) GetJobsImpactedByAnExternalService(repo *db.Repository, jobName string, startDate string, endDate string) (totalSuccess int, err error) {
	return d.client.Repository.QueryProwJobs(repo).
		Where(prowjobs.JobName(jobName)).
		Where(prowjobs.ExternalServicesImpact(true)).
		// Jobs are executed after verifying external service step. That mean an outage can be fixed when run the job and success.
		Where(prowjobs.Not(prowjobs.State("success"))).
		Where(prowjobs.Not(prowjobs.State("success"))).
		Where(func(s *sql.Selector) { // "merged_at BETWEEN ? AND 2022-08-17", "2022-08-16"
			s.Where(sql.ExprP(fmt.Sprintf("created_at BETWEEN '%s' AND '%s'", startDate, endDate)))
		}).
		Aggregate(
			db.Count(),
		).
		Int(context.Background())
}

/*
func (d *Database) GetMetrics(gitOrganization, repoName, jobType, startDate, endDate string) (prowV1Alpha1.JobsMetrics, error) {
	var metrics prowV1Alpha1.JobsMetrics
	metrics.GitOrganization = gitOrganization
	metrics.JobType = jobType
	metrics.RepositoryName = repoName

	repo, err := d.client.Repository.Query().Where(repository.GitOrganization(gitOrganization)).Where(repository.RepositoryName(repoName)).First(context.Background())
	if err != nil {
		return metrics, err
	}

	dbJobs, err := d.client.Repository.QueryProwJobs(repo).Where(prowjobs.JobType(jobType)).All(context.Background())
	if err != nil {
		return metrics, err
	}

	for _, job := range ReturnJobNames(dbJobs) {
		jMetric, err := d.client.Repository.QueryProwJobs(repo).Select().
			Where(prowjobs.JobName(job)).
			Where(prowjobs.JobType(jobType)).
			Where(func(s *sql.Selector) { // "created_at BETWEEN ? AND 2022-08-17", "2022-08-16"
				s.Where(sql.ExprP(fmt.Sprintf("created_at BETWEEN '%s' AND '%s'", startDate, endDate)))
			}).All(context.TODO())
		if err != nil {
			return metrics, err
		}
		metrics.Jobs = append(metrics.Jobs, d.getProwJobSummary(jMetric, repo, job, jobType, startDate, endDate))
	}

	return metrics, nil
}

func (d *Database) getMetric(repo *db.Repository, job, jobType, startDate, endDate string) prowV1Alpha1.Metrics {
	jMetric, _ := d.client.Repository.QueryProwJobs(repo).Select().
		Where(prowjobs.JobName(job)).
		Where(prowjobs.JobType(jobType)).
		Where(func(s *sql.Selector) { // "created_at BETWEEN ? AND 2022-08-17", "2022-08-16"
			s.Where(sql.ExprP(fmt.Sprintf("created_at BETWEEN '%s' AND '%s'", startDate, endDate)))
		}).All(context.TODO())

	return getProwMetricsByDay(jMetric, startDate)
}

func (d *Database) getMetricsSummaryByDay(repo *db.Repository, job, jobType, startDate, endDate string) []prowV1Alpha1.Metrics {
	var metrics []prowV1Alpha1.Metrics
	dayArr := getDatesBetweenRange(startDate, endDate)

	// range between one day (same day)
	if len(dayArr) == 2 && isSameDay(startDate, endDate) {
		metric := d.getMetric(repo, job, jobType, startDate, endDate)
		metrics = append(metrics, metric)
		return metrics
	}

	// range between more than one day
	for i, day := range dayArr {
		t, _ := time.Parse("2006-01-02 15:04:05", day)
		y, m, dd := t.Date()

		if i == 0 { // first day
			metric := d.getMetric(repo, job, jobType, day, fmt.Sprintf("%04d-%02d-%02d 23:59:59", y, m, dd))
			metrics = append(metrics, metric)
		} else {
			if i == len(dayArr)-1 { // last day
				metric := d.getMetric(repo, job, jobType, fmt.Sprintf("%04d-%02d-%02d 00:00:00", y, m, dd), day)
				metrics = append(metrics, metric)
			} else { // middle days
				metric := d.getMetric(repo, job, jobType, fmt.Sprintf("%04d-%02d-%02d 00:00:00", y, m, dd), fmt.Sprintf("%04d-%02d-%02d 23:59:59", y, m, dd))
				metrics = append(metrics, metric)
			}
		}
	}
	return metrics
}

func getProwMetricsByDay(jobs []*db.ProwJobs, date string) prowV1Alpha1.Metrics {
	var success_total, failed_rate_total, not_scheduled_total, failure_by_e2e_tests_total, failure_by_build_errors_total float64

	for _, j := range jobs {
		if j.State == string(prow.SuccessState) {
			success_total++
		}

		if j.State == string(prow.ErrorState) {
			not_scheduled_total++
		}

		if j.State == string(prow.FailureState) {
			if j.BuildErrorLogs != nil && *j.BuildErrorLogs != "" {
				failure_by_build_errors_total++
			}
			if j.E2eFailedTestMessages != nil && *j.E2eFailedTestMessages != "" {
				failure_by_e2e_tests_total++
			}

			failed_rate_total++
		}
	}

	return prowV1Alpha1.Metrics{
		Date:                        date,
		SuccessCount:                success_total,
		FailureCount:                failed_rate_total,
		JobFailedByE2ETestsCount:    failure_by_e2e_tests_total,
		JobFailedByBuildErrorsCount: failure_by_build_errors_total,
		JobNotScheduledCount:        not_scheduled_total,
		TotalJobs:                   success_total + not_scheduled_total + failed_rate_total,
	}

}

func (d *Database) getProwJobSummary(jobs []*db.ProwJobs, repo *db.Repository, jobName, jobType, startDate, endDate string) prowV1Alpha1.Jobs {
	var success_total, failed_rate_total, not_scheduled_total, failure_by_e2e_tests_total, failure_by_build_errors_total float64

	for _, j := range jobs {
		if j.State == string(prow.SuccessState) {
			success_total++
		}

		if j.State == string(prow.ErrorState) {
			not_scheduled_total++
		}

		if j.State == string(prow.FailureState) {
			if j.BuildErrorLogs != nil && *j.BuildErrorLogs != "" {
				failure_by_build_errors_total++
			}
			if j.E2eFailedTestMessages != nil && *j.E2eFailedTestMessages != "" {
				failure_by_e2e_tests_total++
			}

			failed_rate_total++
		}
	}
	metricsByDat := d.getMetricsSummaryByDay(repo, jobName, jobType, startDate, endDate)

	return prowV1Alpha1.Jobs{
		Name:    jobName,
		Metrics: metricsByDat,
		Summary: prowV1Alpha1.Summary{
			DateFrom:                    startDate,
			DateTo:                      endDate,
			SuccessCount:                success_total,
			JobFailedCount:              failed_rate_total,
			JobFailedByE2ETestsCount:    failure_by_e2e_tests_total,
			JobFailedByBuildErrorsCount: failure_by_build_errors_total,
			JobNotScheduledCount:        not_scheduled_total,
			TotalJobs:                   success_total + failed_rate_total + not_scheduled_total,
		},
	}
}

func ReturnJobNames(j []*db.ProwJobs) []string {
	var jobsArr []string

	for _, jobs := range j {
		// periodic-ci-redhat-appstudio-infra-deployments-main-hacbs-e2e-periodic does not exist anymore
		// stopped to be collected from 14/06/2023
		if jobs.JobName != "periodic-ci-redhat-appstudio-infra-deployments-main-hacbs-e2e-periodic" {
			jobsArr = append(jobsArr, jobs.JobName)
		}
	}

	return removeDuplicateStr(jobsArr)
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
*/
