package client

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/andygrunwald/go-jira"
	jiraV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/jira/v1alpha1"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/bugs"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/predicate"
)

// CreateJiraBug saves provided jira bugs information in database.
func (d *Database) CreateJiraBug(bugsArr []jira.Issue, team *db.Teams) error {
	create := false
	createBulk := make([]*db.BugsCreate, 0)
	for _, bug := range bugsArr {
		bugAlreadyExists := d.client.Bugs.Query().Where(bugs.JiraKey(bug.Key)).ExistX(context.TODO())

		var bugIsResolved bool
		var diff float64
		if bug.Fields.Status.Name == "Closed" || bug.Fields.Status.Name == "Resolved" || bug.Fields.Status.Name == "Done" {
			resolvedTime := time.Time(bug.Fields.Resolutiondate).UTC()
			createdTime := time.Time(bug.Fields.Created).UTC()

			diff = resolvedTime.Sub(createdTime).Hours()
			bugIsResolved = true
		}

		if bugAlreadyExists {
			_, err := d.client.Bugs.Update().Where(predicate.Bugs(bugs.JiraKey(bug.Key))).
				SetCreatedAt(time.Time(bug.Fields.Created)).
				SetUpdatedAt(time.Time(bug.Fields.Updated)).
				SetResolvedAt(time.Time(bug.Fields.Resolutiondate)).
				SetResolved(bugIsResolved).
				SetResolutionTime(diff).
				SetJiraKey(bug.Key).
				SetPriority(bug.Fields.Priority.Name).
				SetSummary(bug.Fields.Summary).
				SetURL(fmt.Sprintf("https://issues.redhat.com/browse/%s", bug.Key)).
				SetStatus(bug.Fields.Status.Description).
				SetBugs(team).
				SetProjectKey(getProjectKey(bug.Key)).
				Save(context.TODO())
			if err != nil {
				return convertDBError("failed to create bug: %w", err)
			}
		} else {
			createBulkByBug := d.client.Bugs.Create().
				SetCreatedAt(time.Time(bug.Fields.Created)).
				SetUpdatedAt(time.Time(bug.Fields.Updated)).
				SetResolvedAt(time.Time(bug.Fields.Resolutiondate)).
				SetResolved(bugIsResolved).
				SetResolutionTime(diff).
				SetJiraKey(bug.Key).
				SetPriority(bug.Fields.Priority.Name).
				SetSummary(bug.Fields.Summary).
				SetURL(fmt.Sprintf("https://issues.redhat.com/browse/%s", bug.Key)).
				SetStatus(bug.Fields.Status.Description).
				SetBugs(team).
				SetProjectKey(getProjectKey(bug.Key))
			createBulk = append(createBulk, createBulkByBug)
			create = true
		}
	}

	defer func() {
		if err := recover(); err != nil {
			// Usually occurs when u have network issues
			fmt.Println("Internal panic occurred, check network connection:", err)
		}
	}()

	// todo: https://github.com/ent/ent/issues/2494 Wait until we can resolve nicely conflicts in psql
	if create {
		if err := d.client.Bugs.CreateBulk(createBulk...).OnConflict(sql.ResolveWithNewValues()).DoNothing().Exec(context.TODO()); err != nil {
			return err
		}
	}

	return nil
}

func (d *Database) TotalBugsResolutionTime(priority string, team *db.Teams) (bugsMetrics jiraV1Alpha1.ResolvedBugsMetrics, err error) {
	var totalBugs int
	var totalAverage float64
	currentTime := time.Now()

	for month := 0; month < 12; month++ {
		firstDayOfMonth := BeginningOfMonth(currentTime.AddDate(0, -month, 0))
		lastDayOfMonth := EndOfMonth(currentTime.AddDate(0, -month, 0))

		totalMonthResolution, err := d.ResolutionBugByDate(team, priority, firstDayOfMonth.Format("2006-01-02"), lastDayOfMonth.Format("2006-01-02"))

		if err != nil {
			return jiraV1Alpha1.ResolvedBugsMetrics{}, convertDBError("failed to return bugs: %w", err)
		}

		bugsAll, err := d.GetResolvedBugsByPriorityAndStatus(priority, team, true, firstDayOfMonth.Format("2006-01-02"), lastDayOfMonth.Format("2006-01-02"))
		if err != nil {
			return jiraV1Alpha1.ResolvedBugsMetrics{}, convertDBError("failed to return bugs: %w", err)
		}

		bugsMetrics.ResolutionTimeTotal.Months = append(bugsMetrics.ResolutionTimeTotal.Months, jiraV1Alpha1.MonthsResolution{
			Name:                 fmt.Sprintf("%s_%v", firstDayOfMonth.Month().String(), firstDayOfMonth.Year()),
			Total:                totalMonthResolution,
			NumberOfResolvedBugs: len(bugsAll),
			Bugs:                 bugsAll,
		})
	}

	for _, r := range bugsMetrics.ResolutionTimeTotal.Months {
		totalBugs = totalBugs + r.NumberOfResolvedBugs
		totalAverage = (totalAverage + r.Total)
	}

	return jiraV1Alpha1.ResolvedBugsMetrics{
		ResolutionTimeTotal: jiraV1Alpha1.ResolutionTime{
			Total:             totalAverage / 12,
			Priority:          priority,
			NumberOfTotalBugs: totalBugs,
			Months:            bugsMetrics.ResolutionTimeTotal.Months,
		},
	}, nil
}

func (d *Database) ResolutionBugByDate(team *db.Teams, priority, dateFrom string, dateTo string) (float64, error) {
	var resolution []struct {
		Sum, Count float64
	}

	if priority == "Global" {
		err := d.client.Teams.QueryBugs(team).
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
		err := d.client.Teams.QueryBugs(team).
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

func (d *Database) GetResolvedBugsByPriorityAndStatus(priority string, t *db.Teams, IsResolved bool, dateFrom string, dateTo string) (bugsArray []*db.Bugs, err error) {
	if priority == "Global" {
		bugsArray, err = d.client.Teams.QueryBugs(t).
			Where(predicate.Bugs(bugs.Resolved(IsResolved))).
			Where(func(s *sql.Selector) { // "created_at BETWEEN ? AND 2022-08-17", "2022-08-16"
				s.Where(sql.ExprP(fmt.Sprintf("resolved_at BETWEEN '%s' AND '%s'", dateFrom, dateTo)))
			}).All(context.TODO())
		if err != nil {
			return nil, err
		}
	} else {
		bugsArray, err = d.client.Teams.QueryBugs(t).
			Where(predicate.Bugs(bugs.Resolved(IsResolved))).
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

func (d *Database) GetOpenBugsByPriorityAndStatus(priority string, t *db.Teams, IsResolved bool, dateFrom string, dateTo string) (bugsArray []*db.Bugs, err error) {
	if priority == "Global" {
		bugsArray, err = d.client.Teams.QueryBugs(t).
			Where(predicate.Bugs(bugs.Resolved(IsResolved))).
			Where(func(s *sql.Selector) { // "created_at BETWEEN ? AND 2022-08-17", "2022-08-16"
				s.Where(sql.ExprP(fmt.Sprintf("created_at BETWEEN '%s' AND '%s'", dateFrom, dateTo)))
			}).All(context.TODO())
		if err != nil {
			return nil, err
		}
	} else {
		bugsArray, err = d.client.Teams.QueryBugs(t).
			Where(predicate.Bugs(bugs.Resolved(IsResolved))).
			Where(predicate.Bugs(bugs.Priority(priority))).
			Where(func(s *sql.Selector) { // "created_at BETWEEN ? AND 2022-08-17", "2022-08-16"
				s.Where(sql.ExprP(fmt.Sprintf("created_at BETWEEN '%s' AND '%s'", dateFrom, dateTo)))
			}).All(context.TODO())
		if err != nil {
			return nil, err
		}
	}
	return bugsArray, nil
}

func (d *Database) GetOpenBugsMetricsByStatusAndPriority(priority string, team *db.Teams) (bugsMetrics jiraV1Alpha1.OpenBugsMetrics, err error) {
	var totalBugs int
	currentTime := time.Now()

	for month := 0; month < 12; month++ {
		firstDayOfMonth := BeginningOfMonth(currentTime.AddDate(0, -month, 0))
		lastDayOfMonth := EndOfMonth(currentTime.AddDate(0, -month, 0))
		bugsAll, err := d.GetOpenBugsByPriorityAndStatus(priority, team, false, firstDayOfMonth.Format("2006-01-02"), lastDayOfMonth.Format("2006-01-02"))
		if err != nil {
			return jiraV1Alpha1.OpenBugsMetrics{}, convertDBError("failed to return bugs: %w", err)
		}

		bugsMetrics.TotalOpenBugs.Months = append(bugsMetrics.TotalOpenBugs.Months, jiraV1Alpha1.MonthsOpen{
			Name:     fmt.Sprintf("%s_%v", firstDayOfMonth.Month().String(), firstDayOfMonth.Year()),
			OpenBugs: len(bugsAll),
			Bugs:     bugsAll,
		})
	}

	for _, r := range bugsMetrics.TotalOpenBugs.Months {
		totalBugs = totalBugs + r.OpenBugs
	}

	return jiraV1Alpha1.OpenBugsMetrics{
		TotalOpenBugs: jiraV1Alpha1.OpenBugs{
			Priority:         priority,
			NumberOfOpenBugs: totalBugs,
			Months:           bugsMetrics.TotalOpenBugs.Months,
		},
	}, nil

}

func (d *Database) GetAllJiraBugs() ([]*db.Bugs, error) {
	bugsAll, err := d.client.Bugs.Query().All(context.Background())
	if err != nil {
		return nil, convertDBError("failed to return bugs: %w", err)
	}

	return bugsAll, nil
}

func (d *Database) DeleteJiraBugsByProject(projectKey string, team *db.Teams) error {
	_, err := d.client.Bugs.Delete().Where(predicate.Bugs(bugs.ProjectKey(projectKey))).Exec(context.TODO())

	if err != nil {
		return convertDBError("failed to delete jira bugs: %w", err)
	}

	return nil
}

func BeginningOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 0, -date.Day()+1)
}

func EndOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 1, -date.Day())
}

func getProjectKey(bugKey string) string {
	bugKeySplit := strings.Split(bugKey, "-")

	if len(bugKeySplit) > 0 {
		return bugKeySplit[0]
	}

	return ""
}
