package client

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	prV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/github/v1alpha1"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
)

func (d *Database) getMetricByDay(repo *db.Repository, startDate, endDate string) prV1Alpha1.Metrics {
	createdPrs, _ := d.client.Repository.QueryPrs(repo).Select().
		Where(func(s *sql.Selector) { // "created_at BETWEEN ? AND 2022-08-17", "2022-08-16"
			s.Where(sql.ExprP(fmt.Sprintf("created_at BETWEEN '%s' AND '%s'", startDate, endDate)))
		}).All(context.TODO())

	mergedPrs, _ := d.client.Repository.QueryPrs(repo).Select().
		Where(func(s *sql.Selector) { // "merged_at BETWEEN ? AND 2022-08-17", "2022-08-16"
			s.Where(sql.ExprP(fmt.Sprintf("merged_at BETWEEN '%s' AND '%s'", startDate, endDate)))
		}).All(context.TODO())

	return prV1Alpha1.Metrics{
		Date:                     startDate,
		CreatedPullRequestsCount: len(createdPrs),
		MergedPullRequestsCount:  len(mergedPrs),
	}
}

func (d *Database) getMetrics(repository *db.Repository, startDate, endDate string) []prV1Alpha1.Metrics {
	var metrics []prV1Alpha1.Metrics
	dayArr := getDatesBetweenRange(startDate, endDate)

	// range between one day (same day)
	if len(dayArr) == 2 && isSameDay(startDate, endDate) {
		metric := d.getMetricByDay(repository, startDate, endDate)
		metrics = append(metrics, metric)
		return metrics
	}

	// range between more than one day
	for i, day := range dayArr {
		t, _ := time.Parse("2006-01-02 15:04:05", day)
		y, m, dd := t.Date()

		if i == 0 { // first day
			metric := d.getMetricByDay(repository, day, fmt.Sprintf("%04d-%02d-%02d 23:59:59", y, m, dd))
			metrics = append(metrics, metric)
		} else {
			if i == len(dayArr)-1 { // last day
				metric := d.getMetricByDay(repository, fmt.Sprintf("%04d-%02d-%02d 00:00:00", y, m, dd), day)
				metrics = append(metrics, metric)
			} else { // middle days
				metric := d.getMetricByDay(repository, fmt.Sprintf("%04d-%02d-%02d 00:00:00", y, m, dd), fmt.Sprintf("%04d-%02d-%02d 23:59:59", y, m, dd))
				metrics = append(metrics, metric)
			}
		}
	}
	return metrics
}
