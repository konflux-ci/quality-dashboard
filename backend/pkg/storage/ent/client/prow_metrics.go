package client

import (
	"context"
	"fmt"
	"math"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/redhat-appstudio/quality-studio/api/server/router/prow"
	"github.com/redhat-appstudio/quality-studio/pkg/storage"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/prowjobs"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/repository"
)

func (d *Database) GetMetrics(gitOrganization, repoName, jobType, startDate, endDate string) storage.ProwJobsMetrics {
	var metrics storage.ProwJobsMetrics
	metrics.GitOrganization = gitOrganization
	metrics.JobType = jobType
	metrics.RepositoryName = repoName

	repo, _ := d.client.Repository.Query().Where(repository.GitOrganization(gitOrganization)).Where(repository.RepositoryName(repoName)).First(context.Background())

	dbJobs, _ := d.client.Repository.QueryProwJobs(repo).Where(prowjobs.JobType(jobType)).All(context.Background())

	for _, job := range ReturnJobNames(dbJobs) {
		jMetric, _ := d.client.Repository.QueryProwJobs(repo).Select().
			Where(prowjobs.JobName(job)).
			Where(prowjobs.JobType(jobType)).
			Where(func(s *sql.Selector) { // "created_at BETWEEN ? AND 2022-08-17", "2022-08-16"
				s.Where(sql.ExprP(fmt.Sprintf("created_at BETWEEN '%s' AND '%s'", startDate, endDate)))
			}).All(context.TODO())

		metrics.Jobs = append(metrics.Jobs, d.getProwJobSummary(jMetric, repo, job, jobType, startDate, endDate))
	}

	return metrics
}

func (d *Database) getMetric(repo *db.Repository, job, jobType, startDate, endDate string) storage.Metrics {
	jMetric, _ := d.client.Repository.QueryProwJobs(repo).Select().
		Where(prowjobs.JobName(job)).
		Where(prowjobs.JobType(jobType)).
		Where(func(s *sql.Selector) { // "created_at BETWEEN ? AND 2022-08-17", "2022-08-16"
			s.Where(sql.ExprP(fmt.Sprintf("created_at BETWEEN '%s' AND '%s'", startDate, endDate)))
		}).All(context.TODO())

	return getProwMetricsByDay(jMetric, startDate)
}

func (d *Database) getMetricsSummaryByDay(repo *db.Repository, job, jobType, startDate, endDate string) []storage.Metrics {
	var metrics []storage.Metrics
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
			metric := d.getMetric(repo, job, jobType, day, fmt.Sprintf("%d-%d-%d 23:59:59", y, m, dd))
			metrics = append(metrics, metric)
		} else {
			if i == len(dayArr)-1 { // last day
				metric := d.getMetric(repo, job, jobType, fmt.Sprintf("%d-%d-%d 00:00:00", y, m, dd), day)
				metrics = append(metrics, metric)
			} else { // middle days
				metric := d.getMetric(repo, job, jobType, fmt.Sprintf("%d-%d-%d 00:00:00", y, m, dd), fmt.Sprintf("%d-%d-%d 23:59:59", y, m, dd))
				metrics = append(metrics, metric)
			}
		}
	}
	return metrics
}

func getProwMetricsByDay(jobs []*db.ProwJobs, date string) storage.Metrics {
	job_nums := float64(len(jobs))
	var success_rate_total, failed_rate_total, ci_failed_total float64

	for _, j := range jobs {
		if j.State == string(prow.SuccessState) {
			success_rate_total = success_rate_total + 1
		}

		if j.State == string(prow.ErrorState) {
			ci_failed_total = ci_failed_total + 1
		}

		if j.State == string(prow.FailureState) {
			failed_rate_total = failed_rate_total + 1
		}
	}

	success_rate := success_rate_total / job_nums * 100
	if math.IsNaN(success_rate) {
		success_rate = 0
	}

	failed_rate := failed_rate_total / job_nums * 100
	if math.IsNaN(failed_rate) {
		failed_rate = 0
	}

	ci_failed_rate := ci_failed_total / job_nums * 100
	if math.IsNaN(ci_failed_rate) {
		ci_failed_rate = 0
	}

	return storage.Metrics{
		Date:         date,
		SuccessRate:  success_rate,
		FailureRate:  failed_rate,
		CiFailedRate: ci_failed_rate,
	}

}

func (d *Database) getProwJobSummary(jobs []*db.ProwJobs, repo *db.Repository, jobName, jobType, startDate, endDate string) storage.Jobs {
	var success_rate_total, failed_rate_total, ci_failed_total float64
	finishedJobsNumber := 0

	for _, j := range jobs {
		if j.State == string(prow.SuccessState) {
			success_rate_total++
			finishedJobsNumber++
		}

		if j.State == string(prow.ErrorState) {
			ci_failed_total++
			finishedJobsNumber++
		}

		if j.State == string(prow.FailureState) {
			failed_rate_total++
			finishedJobsNumber++
		}
	}
	job_nums := float64(len(jobs))
	metricsByDat := d.getMetricsSummaryByDay(repo, jobName, jobType, startDate, endDate)

	return storage.Jobs{
		Name:    jobName,
		Metrics: metricsByDat,
		Summary: storage.Summary{
			DateFrom:       startDate,
			DateTo:         endDate,
			SuccessRateAvg: success_rate_total / job_nums * 100,
			JobFailedAvg:   failed_rate_total / job_nums * 100,
			CIFailedAvg:    ci_failed_total / job_nums * 100,
			TotalJobs:      finishedJobsNumber,
		},
	}
}

func ReturnJobNames(j []*db.ProwJobs) []string {
	var jobsArr []string

	for _, jobs := range j {
		jobsArr = append(jobsArr, jobs.JobName)
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
