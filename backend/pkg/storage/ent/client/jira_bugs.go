package client

import (
	"context"
	"fmt"
	"math"
	"time"

	"entgo.io/ent/dialect/sql"
	jiraV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/jira/v1alpha1"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/bugs"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/predicate"
)

// CreateJiraBug saves provided jira bug information in database.
func (d *Database) CreateJiraBug(jiraBug jiraV1Alpha1.JiraBug) error {
	bugAlreadyExists := d.client.Bugs.Query().Where(bugs.JiraKey(jiraBug.JiraKey)).ExistX(context.TODO())

	if bugAlreadyExists {
		_, err := d.client.Bugs.Update().Where(predicate.Bugs(bugs.JiraKey(jiraBug.JiraKey))).
			SetCreatedAt(jiraBug.CreatedAt).
			SetUpdatedAt(jiraBug.UpdatedAt).
			SetResolvedAt(jiraBug.ResolvedAt).
			SetResolved(jiraBug.IsResolved).
			SetResolutionTime(jiraBug.ResolutionTime).
			SetJiraKey(jiraBug.JiraKey).
			SetPriority(jiraBug.Priority).
			SetSummary(jiraBug.Summary).
			SetURL(jiraBug.Url).
			SetStatus(jiraBug.Status).
			Save(context.TODO())
		if err != nil {
			return convertDBError("failed to create bug: %w", err)
		}
	} else {
		_, err := d.client.Bugs.Create().
			SetCreatedAt(jiraBug.CreatedAt).
			SetUpdatedAt(jiraBug.UpdatedAt).
			SetResolutionTime(jiraBug.ResolutionTime).
			SetJiraKey(jiraBug.JiraKey).
			SetPriority(jiraBug.Priority).
			SetSummary(jiraBug.Summary).
			SetURL(jiraBug.Url).
			SetResolvedAt(jiraBug.ResolvedAt).
			SetResolved(jiraBug.IsResolved).
			SetStatus(jiraBug.Status).
			Save(context.TODO())
		if err != nil {
			return convertDBError("failed to update bug: %w", err)
		}
	}

	return nil
}

func (d *Database) TotalBugsResolutionTime(priority string) (bugsMetrics jiraV1Alpha1.BugsMetrics, err error) {

	var totalBugs int
	var totalAverage float64
	currentTime := time.Now()

	for month := 0; month < 12; month++ {
		firstDayOfMonth := BeginningOfMonth(currentTime.AddDate(0, -month, 0))
		lastDayOfMonth := EndOfMonth(currentTime.AddDate(0, -month, 0))

		totalMonthResolution, err := d.ResolutionBugByDate(priority, firstDayOfMonth.Format("2006-01-02"), lastDayOfMonth.Format("2006-01-02"))

		if err != nil {
			return jiraV1Alpha1.BugsMetrics{}, convertDBError("failed to return bugs: %w", err)
		}

		bugsAll, err := d.GetBugsByDate(priority, firstDayOfMonth.Format("2006-01-02"), lastDayOfMonth.Format("2006-01-02"))
		if err != nil {
			return jiraV1Alpha1.BugsMetrics{}, convertDBError("failed to return bugs: %w", err)
		}

		bugsMetrics.ResolutionTimeTotal.Months = append(bugsMetrics.ResolutionTimeTotal.Months, jiraV1Alpha1.MonthsResolution{
			Name:                 firstDayOfMonth.Month().String(),
			Total:                totalMonthResolution,
			NumberOfResolvedBugs: len(bugsAll),
			Bugs:                 bugsAll,
		})
	}

	for _, r := range bugsMetrics.ResolutionTimeTotal.Months {
		totalBugs = totalBugs + r.NumberOfResolvedBugs
		totalAverage = (totalAverage + r.Total)
	}

	return jiraV1Alpha1.BugsMetrics{
		ResolutionTimeTotal: jiraV1Alpha1.ResolutionTime{
			Total:             totalAverage / 12,
			Priority:          priority,
			NumberOfTotalBugs: totalBugs,
			Months:            bugsMetrics.ResolutionTimeTotal.Months,
		},
	}, nil
}

func (d *Database) ResolutionBugByDate(priority, dateFrom string, dateTo string) (float64, error) {
	var resolution []struct {
		Sum, Count float64
	}

	if priority == "Global" {
		err := d.client.Bugs.Query().
			Where(predicate.Bugs(bugs.Resolved(true))).
			Where(func(s *sql.Selector) {
				s.Where(sql.ExprP(fmt.Sprintf("resolved_at BETWEEN '%s' AND '%s'", dateFrom, dateTo)))
			}).
			Aggregate(
				db.Sum(bugs.FieldResolutionTime),
				db.Count(),
			).
			Scan(context.TODO(), &resolution)

		if err != nil {
			return 0, err
		}
	} else {
		err := d.client.Bugs.Query().
			Where(predicate.Bugs(bugs.Priority(priority))).
			Where(predicate.Bugs(bugs.Resolved(true))).
			Where(func(s *sql.Selector) {
				s.Where(sql.ExprP(fmt.Sprintf("resolved_at BETWEEN '%s' AND '%s'", dateFrom, dateTo)))
			}).
			Aggregate(
				db.Sum(bugs.FieldResolutionTime),
				db.Count(),
			).
			Scan(context.TODO(), &resolution)

		if err != nil {
			return 0, err
		}
	}

	totalMonthResolution := resolution[0].Sum / resolution[0].Count

	if math.IsNaN(totalMonthResolution) {
		totalMonthResolution = 0
	}

	return totalMonthResolution, nil
}

func (d *Database) GetBugsByDate(priority string, dateFrom string, dateTo string) (bugsArray []*db.Bugs, err error) {
	if priority == "Global" {
		bugsArray, err = d.client.Bugs.Query().
			Where(predicate.Bugs(bugs.Resolved(true))).
			Where(func(s *sql.Selector) { // "created_at BETWEEN ? AND 2022-08-17", "2022-08-16"
				s.Where(sql.ExprP(fmt.Sprintf("resolved_at BETWEEN '%s' AND '%s'", dateFrom, dateTo)))
			}).All(context.TODO())
		if err != nil {
			return nil, err
		}
	} else {
		bugsArray, err = d.client.Bugs.Query().
			Where(predicate.Bugs(bugs.Resolved(true))).
			Where(predicate.Bugs(bugs.Priority(priority))).
			Where(func(s *sql.Selector) { // "created_at BETWEEN ? AND 2022-08-17", "2022-08-16"
				s.Where(sql.ExprP(fmt.Sprintf("resolved_at BETWEEN '%s' AND '%s'", dateFrom, dateTo)))
			}).All(context.TODO())
		if err != nil {
			return nil, err
		}
	}
	return bugsArray, nil
}

func (d *Database) GetAllJiraBugs() ([]*db.Bugs, error) {
	bugsAll, err := d.client.Bugs.Query().All(context.Background())
	if err != nil {
		return nil, convertDBError("failed to return bugs: %w", err)
	}

	return bugsAll, nil
}

func BeginningOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 0, -date.Day()+1)
}

func EndOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 1, -date.Day())
}
