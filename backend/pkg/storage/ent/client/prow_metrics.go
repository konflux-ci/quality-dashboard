package client

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	prowV1Alpha1 "github.com/konflux-ci/quality-dashboard/api/apis/prow/v1alpha1"
	"github.com/konflux-ci/quality-dashboard/api/server/router/prow"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db/prowjobs"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db/repository"
)

func (d *Database) ObtainProwMetricsByJob(gitOrganization string, repositoryName string, jobName string, startDate string, endDate string) (*prowV1Alpha1.JobsMetrics, error) {
	repo, err := d.client.Repository.Query().Where(repository.GitOrganization(gitOrganization)).Where(repository.RepositoryName(repositoryName)).First(context.Background())
	if err != nil {
		return nil, err
	}

	dbJobs, err := d.client.Repository.QueryProwJobs(repo).
		Where(prowjobs.JobName(jobName)).
		Where(func(s *sql.Selector) { // "created_at BETWEEN ? AND 2022-08-17", "2022-08-16"
			s.Where(sql.ExprP(fmt.Sprintf("created_at BETWEEN '%s' AND '%s'", startDate, endDate)))
		}).
		Where(prowjobs.Not(prowjobs.State(string(prow.AbortedState)))).
		All(context.Background())
	if err != nil {
		return nil, err
	}

	numberOfSuccessJobs, err := d.GetNumberOfSuccessJobs(repo, jobName, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get number of success from db %v", err)
	}

	numberOfFailedJobs, err := d.GetNumberOfFailedJobs(repo, jobName, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get number of failed jobs %v", err)
	}

	numberOfInfraImpactJobs, err := d.GetNumberOfInfraImpact(repo, jobName, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get number of jobs impacted by infrastructure %v", err)
	}

	flaky, err := d.GetSuitesFailureFrequency(gitOrganization, repositoryName, jobName, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get impact of flaky tests %v", err)
	}

	extImpct, err := d.GetJobsImpactedByAnExternalService(repo, jobName, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get impact of external impact %v", err)
	}

	totalImpact := numberOfFailedJobs + numberOfInfraImpactJobs
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
	}, nil
}

func (d *Database) GetMetricsSummaryByDay(repo *db.Repository, job, startDate, endDate string) []*prowV1Alpha1.JobsMetrics {
	var metrics []*prowV1Alpha1.JobsMetrics
	dayArr := getDatesBetweenRange(startDate, endDate)

	// range between one day (same day). !TODO: Check error
	if len(dayArr) == 2 && isSameDay(startDate, endDate) {
		metric, _ := d.ObtainProwMetricsByJob(repo.GitOrganization, repo.RepositoryName, job, startDate, endDate)
		metrics = append(metrics, metric)
		return metrics
	}

	// range between more than one day
	for i, day := range dayArr {
		t, _ := time.Parse("2006-01-02 15:04:05", day)
		y, m, dd := t.Date()

		if i == 0 { // first day !TODO: Check error
			metric, _ := d.ObtainProwMetricsByJob(repo.GitOrganization, repo.RepositoryName, job, day, fmt.Sprintf("%04d-%02d-%02d 23:59:59", y, m, dd))
			metrics = append(metrics, metric)
		} else {
			if i == len(dayArr)-1 { // last day !TODO: Check error
				metric, _ := d.ObtainProwMetricsByJob(repo.GitOrganization, repo.RepositoryName, job, fmt.Sprintf("%04d-%02d-%02d 00:00:00", y, m, dd), day)
				metrics = append(metrics, metric)
			} else { // middle days !TODO: Check error
				metric, _ := d.ObtainProwMetricsByJob(repo.GitOrganization, repo.RepositoryName, job, fmt.Sprintf("%04d-%02d-%02d 00:00:00", y, m, dd), fmt.Sprintf("%04d-%02d-%02d 23:59:59", y, m, dd))
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
		Where(func(s *sql.Selector) { // "merged_at BETWEEN ? AND 2022-08-17", "2022-08-16"
			s.Where(sql.ExprP(fmt.Sprintf("created_at BETWEEN '%s' AND '%s'", startDate, endDate)))
		}).
		Aggregate(
			db.Count(),
		).
		Int(context.Background())
}
