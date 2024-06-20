package client

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/andygrunwald/go-jira"
	jiraV1Alpha1 "github.com/konflux-ci/quality-dashboard/api/apis/jira/v1alpha1"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db/bugs"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db/predicate"
)

type JiraBugMetricsInfo struct {
	// number of days that took to a bug to be assigned (for closed bugs)
	AssignmentTime float64

	// number of days that took to a bug to be prioritized (for closed bugs)
	PrioritizationTime float64

	// number of days that took to a bug to be resolved (for closed bugs)
	ResolutionTime float64

	// current number of days that a bug is not assigned (for open bugs)
	DaysWithoutAssignee float64

	// current number of days that a bug is not prioritized (for open bugs)
	DaysWithoutPriority float64

	// current number of days that a bug is not resolved (for open bugs)
	DaysWithoutResolution float64

	// current number of days that a bug does not have component assigned (for open bugs)
	DaysWithoutComponent float64

	// if the bub is resolved
	BugIsResolved bool
}

func getComponent(components []*jira.Component) string {
	if len(components) != 0 {
		component := components[0].Name
		if component == "docs" && len(components) > 1 {
			component = components[1].Name
		}
		return component
	}

	return "undefined"
}

// CreateJiraBug saves provided jira bugs information in database.
func (d *Database) CreateJiraBug(bugsArr []jira.Issue, team *db.Teams) error {
	bulkSize := 2000

	if len(bugsArr) > bulkSize {
		// number of issues is too high
		// probably will hit a similar error like: 'Update failed: insert nodes to table "bugs": pq: got 554645 parameters but PostgreSQL only supports 65535 parameters'
		// we will need to split in smaller bulks
		for start := 0; start < len(bugsArr); start += bulkSize {
			end := start + bulkSize
			if end > len(bugsArr) {
				end = len(bugsArr)
			}

			err := d.CreateBug(bugsArr[start:end], team)
			if err != nil {
				return err
			}
		}
		return nil
	}

	err := d.CreateBug(bugsArr, team)
	return err
}

// CreateBug saves provided jira bugs information in database.
func (d *Database) CreateBug(bugsArr []jira.Issue, team *db.Teams) error {
	create := false
	createBulk := make([]*db.BugsCreate, 0)
	for _, bug := range bugsArr {
		bugAlreadyExists := d.client.Bugs.Query().Where(bugs.JiraKey(bug.Key)).ExistX(context.TODO())
		jiraBugMetricsInfo := d.getJiraBugMetrics(bug)

		if bugAlreadyExists {
			_, err := d.client.Bugs.Update().Where(predicate.Bugs(bugs.JiraKey(bug.Key))).
				SetCreatedAt(time.Time(bug.Fields.Created)).
				SetUpdatedAt(time.Time(bug.Fields.Updated)).
				SetResolvedAt(time.Time(bug.Fields.Resolutiondate)).
				SetResolved(jiraBugMetricsInfo.BugIsResolved).
				SetResolutionTime(jiraBugMetricsInfo.ResolutionTime).
				SetJiraKey(bug.Key).
				SetPriority(bug.Fields.Priority.Name).
				SetSummary(bug.Fields.Summary).
				SetURL(fmt.Sprintf("https://issues.redhat.com/browse/%s", bug.Key)).
				SetStatus(bug.Fields.Status.Name).
				SetBugs(team).
				SetProjectKey(getProjectKey(bug.Key)).
				SetAssignmentTime(jiraBugMetricsInfo.AssignmentTime).
				SetPrioritizationTime(jiraBugMetricsInfo.PrioritizationTime).
				SetDaysWithoutAssignee(jiraBugMetricsInfo.DaysWithoutAssignee).
				SetDaysWithoutPriority(jiraBugMetricsInfo.DaysWithoutPriority).
				SetDaysWithoutResolution(jiraBugMetricsInfo.DaysWithoutResolution).
				SetDaysWithoutComponent(jiraBugMetricsInfo.DaysWithoutComponent).
				SetLabels(strings.Join(bug.Fields.Labels, ",")).
				SetComponent(getComponent(bug.Fields.Components)).
				SetAssignee(getAssignee(bug.Fields.Assignee)).
				SetAge(getDays(bug.Fields.Created, bug.Fields.Resolutiondate, bug.Fields.Status.Name)).
				Save(context.TODO())
			if err != nil {
				return convertDBError("failed to create bug: %w", err)
			}
		} else {
			createBulkByBug := d.client.Bugs.Create().
				SetCreatedAt(time.Time(bug.Fields.Created)).
				SetUpdatedAt(time.Time(bug.Fields.Updated)).
				SetResolvedAt(time.Time(bug.Fields.Resolutiondate)).
				SetResolved(jiraBugMetricsInfo.BugIsResolved).
				SetResolutionTime(jiraBugMetricsInfo.ResolutionTime).
				SetJiraKey(bug.Key).
				SetPriority(bug.Fields.Priority.Name).
				SetSummary(bug.Fields.Summary).
				SetURL(fmt.Sprintf("https://issues.redhat.com/browse/%s", bug.Key)).
				SetStatus(bug.Fields.Status.Name).
				SetBugs(team).
				SetProjectKey(getProjectKey(bug.Key)).
				SetAssignmentTime(jiraBugMetricsInfo.AssignmentTime).
				SetPrioritizationTime(jiraBugMetricsInfo.PrioritizationTime).
				SetDaysWithoutAssignee(jiraBugMetricsInfo.DaysWithoutAssignee).
				SetDaysWithoutPriority(jiraBugMetricsInfo.DaysWithoutPriority).
				SetDaysWithoutResolution(jiraBugMetricsInfo.DaysWithoutResolution).
				SetDaysWithoutComponent(jiraBugMetricsInfo.DaysWithoutComponent).
				SetLabels(strings.Join(bug.Fields.Labels, ",")).
				SetAssignee(getAssignee(bug.Fields.Assignee)).
				SetAge(getDays(bug.Fields.Created, bug.Fields.Resolutiondate, bug.Fields.Status.Name)).
				SetComponent(getComponent(bug.Fields.Components))
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

func (d *Database) TotalBugsResolutionTime(priority, startDate, endDate string, team *db.Teams) (bugsMetrics jiraV1Alpha1.ResolvedBugsMetrics, err error) {
	var totalBugs int
	var totalAverage float64
	dayArr := getDatesBetweenRange(startDate, endDate)

	// range between one day (same day)
	if len(dayArr) == 2 && isSameDay(startDate, endDate) {
		t, _ := time.Parse("2006-01-02 15:04:05", startDate)
		y, m, dd := t.Date()

		totalDayResolution, err := d.ResolutionBugByDate(team, priority, startDate, endDate)
		if err != nil {
			return jiraV1Alpha1.ResolvedBugsMetrics{}, convertDBError("failed to return bugs: %w", err)
		}

		bugsAll, err := d.GetBugsByPriorityAndStatus(priority, team, true, "resolved_at", startDate, endDate)
		if err != nil {
			return jiraV1Alpha1.ResolvedBugsMetrics{}, convertDBError("failed to return bugs: %w", err)
		}

		bugsMetrics.ResolutionTimeTotal.Days = append(bugsMetrics.ResolutionTimeTotal.Days, jiraV1Alpha1.DaysResolution{
			Name:                 fmt.Sprintf("%04d-%02d-%02d", y, m, dd),
			Total:                totalDayResolution,
			NumberOfResolvedBugs: len(bugsAll),
			Bugs:                 bugsAll,
		})

		return bugsMetrics, nil
	}

	// range between more than one day
	for i, day := range dayArr {
		start, end := getRange(i, day, dayArr)

		totalDayResolution, err := d.ResolutionBugByDate(team, priority, start, end)
		if err != nil {
			return jiraV1Alpha1.ResolvedBugsMetrics{}, convertDBError("failed to return bugs: %w", err)
		}

		bugsAll, err := d.GetBugsByPriorityAndStatus(priority, team, true, "resolved_at", start, end)
		if err != nil {
			return jiraV1Alpha1.ResolvedBugsMetrics{}, convertDBError("failed to return bugs: %w", err)
		}

		bugsMetrics.ResolutionTimeTotal.Days = append(bugsMetrics.ResolutionTimeTotal.Days, jiraV1Alpha1.DaysResolution{
			Name:                 day,
			Total:                totalDayResolution,
			NumberOfResolvedBugs: len(bugsAll),
			Bugs:                 bugsAll,
		})
	}

	for _, r := range bugsMetrics.ResolutionTimeTotal.Days {
		totalBugs = totalBugs + r.NumberOfResolvedBugs
		totalAverage = (totalAverage + r.Total)
	}

	totalResolutionTime := totalAverage / float64(len(bugsMetrics.ResolutionTimeTotal.Days))

	return jiraV1Alpha1.ResolvedBugsMetrics{
		ResolutionTimeTotal: jiraV1Alpha1.ResolutionTime{
			Total:             math.Round(totalResolutionTime*100) / 100,
			Priority:          priority,
			NumberOfTotalBugs: totalBugs,
			Days:              bugsMetrics.ResolutionTimeTotal.Days,
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

	totalResolution := resolution[0].Sum / resolution[0].Count

	if math.IsNaN(totalResolution) {
		totalResolution = 0
	}

	return math.Round(totalResolution*100) / 100, nil
}

func (d *Database) GetBugsByPriorityAndStatus(priority string, t *db.Teams, IsResolved bool, exp string, dateFrom string, dateTo string) (bugsArray []*db.Bugs, err error) {
	if priority == "Global" {
		bugsArray, err = d.client.Teams.QueryBugs(t).
			Where(predicate.Bugs(bugs.Resolved(IsResolved))).
			Where(func(s *sql.Selector) { // "created_at BETWEEN ? AND 2022-08-17", "2022-08-16"
				s.Where(sql.ExprP(fmt.Sprintf("%s BETWEEN '%s' AND '%s'", exp, dateFrom, dateTo)))
			}).All(context.TODO())
		if err != nil {
			return nil, err
		}
	} else {
		bugsArray, err = d.client.Teams.QueryBugs(t).
			Where(predicate.Bugs(bugs.Resolved(IsResolved))).
			Where(predicate.Bugs(bugs.Priority(priority))).
			Where(func(s *sql.Selector) { // "created_at BETWEEN ? AND 2022-08-17", "2022-08-16"
				s.Where(sql.ExprP(fmt.Sprintf("%s BETWEEN '%s' AND '%s'", exp, dateFrom, dateTo)))
			}).All(context.TODO())
		if err != nil {
			return nil, err
		}
	}
	return bugsArray, nil
}

func (d *Database) GetOpenBugsMetricsByStatusAndPriority(priority, startDate, endDate string, team *db.Teams) (bugsMetrics jiraV1Alpha1.OpenBugsMetrics, err error) {
	var totalBugs int
	dayArr := getDatesBetweenRange(startDate, endDate)

	// range between one day (same day)
	if len(dayArr) == 2 && isSameDay(startDate, endDate) {
		t, _ := time.Parse("2006-01-02 15:04:05", startDate)
		y, m, dd := t.Date()
		bugsAll, err := d.GetBugsByPriorityAndStatus(priority, team, false, "created_at", startDate, endDate)
		if err != nil {
			return jiraV1Alpha1.OpenBugsMetrics{}, convertDBError("failed to return bugs: %w", err)
		}
		bugsMetrics.TotalOpenBugs.Days = append(bugsMetrics.TotalOpenBugs.Days, jiraV1Alpha1.DaysOpen{
			Name:     fmt.Sprintf("%04d-%02d-%02d", y, m, dd),
			OpenBugs: len(bugsAll),
			Bugs:     bugsAll,
		})
		return bugsMetrics, nil
	}

	// range between more than one day
	for i, day := range dayArr {
		start, end := getRange(i, day, dayArr)

		bugsAll, err := d.GetBugsByPriorityAndStatus(priority, team, false, "created_at", start, end)
		if err != nil {
			return jiraV1Alpha1.OpenBugsMetrics{}, convertDBError("failed to return bugs: %w", err)
		}
		bugsMetrics.TotalOpenBugs.Days = append(bugsMetrics.TotalOpenBugs.Days, jiraV1Alpha1.DaysOpen{
			Name:     day,
			OpenBugs: len(bugsAll),
			Bugs:     bugsAll,
		})
	}

	for _, r := range bugsMetrics.TotalOpenBugs.Days {
		totalBugs = totalBugs + r.OpenBugs
	}

	return jiraV1Alpha1.OpenBugsMetrics{
		TotalOpenBugs: jiraV1Alpha1.OpenBugs{
			Priority:         priority,
			NumberOfOpenBugs: totalBugs,
			Days:             bugsMetrics.TotalOpenBugs.Days,
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

func (d *Database) GetAllJiraBugsByProject(project string) ([]*db.Bugs, error) {
	bugsAll, err := d.client.Bugs.Query().
		Where(predicate.Bugs(bugs.ProjectKey(project))).
		All(context.Background())
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

func (d *Database) DeleteJiraBugByJiraKey(jiraKey string) error {
	_, err := d.client.Bugs.Delete().
		Where(predicate.Bugs(bugs.JiraKey(jiraKey))).
		Exec(context.TODO())

	if err != nil {
		return convertDBError("failed to delete jira bug: %w", err)
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

func (d *Database) GetJiraBug(key string) (*db.Bugs, error) {
	bug, err := d.client.Bugs.Query().
		Where(bugs.JiraKey(key)).
		First(context.Background())
	if err != nil {
		return nil, err
	}

	return bug, nil
}

func (d *Database) GetJiraStatus(key string) (string, error) {
	bug, err := d.client.Bugs.Query().
		Where(bugs.JiraKey(key)).
		First(context.Background())
	if err != nil {
		return "", err
	}

	return bug.Status, nil
}

func (d *Database) BugExists(projectKey string, t *db.Teams) (bool, error) {
	jiraKeys, err := d.client.Teams.QueryBugs(t).
		Where(predicate.Bugs(bugs.JiraKey(projectKey))).
		All(context.TODO())

	if err != nil {
		return false, err
	}

	if len(jiraKeys) == 0 {
		return false, fmt.Errorf("no jira key '%s' found", projectKey)
	}

	return true, nil
}

func (d *Database) getJiraBugMetrics(bug jira.Issue) JiraBugMetricsInfo {
	jiraBugMetric := JiraBugMetricsInfo{
		AssignmentTime:        -1,
		PrioritizationTime:    -1,
		ResolutionTime:        -1,
		DaysWithoutAssignee:   -1,
		DaysWithoutPriority:   -1,
		DaysWithoutResolution: -1,
		DaysWithoutComponent:  -1,
		BugIsResolved:         false,
	}

	createdTime := time.Time(bug.Fields.Created).UTC()
	workingDaysSinceCreation := getWorkingDays(createdTime, time.Now().UTC())

	if bug.Fields.Status.Name == "Closed" || bug.Fields.Status.Name == "Resolved" || bug.Fields.Status.Name == "Done" {
		// issue was closed
		jiraBugMetric.ResolutionTime = getDaysBetweenDates(createdTime, time.Time(bug.Fields.Resolutiondate).UTC())
		jiraBugMetric.BugIsResolved = true
	} else {
		// issue was not resolved
		jiraBugMetric.DaysWithoutResolution = workingDaysSinceCreation
	}

	foundFirstAssignee := false
	foundFirstPriority := false
	for _, history := range bug.Changelog.Histories {
		for _, item := range history.Items {
			if item.Field == "assignee" && item.FromString == "Undefined" && bug.Fields.Assignee != nil && !foundFirstAssignee {
				historyTime, err := history.CreatedTime()
				if err != nil {
					fmt.Println("error getting the CreatedTime of Jira bug's history")
				}

				jiraBugMetric.AssignmentTime = getDaysBetweenDates(createdTime, historyTime.UTC())
				foundFirstAssignee = true
			}
			if item.Field == "priority" && item.FromString == "Undefined" && bug.Fields.Priority != nil && !foundFirstPriority {
				historyTime, err := history.CreatedTime()
				if err != nil {
					fmt.Println("error getting the CreatedTime of Jira bug's history")
				}

				jiraBugMetric.PrioritizationTime = getDaysBetweenDates(createdTime, historyTime.UTC())
				foundFirstPriority = true
			}
		}
	}

	// assignee was defined during the issue creation
	if jiraBugMetric.AssignmentTime == -1 && bug.Fields.Assignee != nil {
		jiraBugMetric.AssignmentTime = 0
	}

	// priority was defined during the issue creation
	if jiraBugMetric.PrioritizationTime == -1 && bug.Fields.Priority != nil && bug.Fields.Priority.Name != "Undefined" {
		jiraBugMetric.PrioritizationTime = 0
	}

	// assignee was not defined
	if bug.Fields.Assignee == nil {
		jiraBugMetric.DaysWithoutAssignee = workingDaysSinceCreation
	}

	// priority was not defined
	if bug.Fields.Priority == nil || bug.Fields.Priority.Name == "Undefined" {
		jiraBugMetric.DaysWithoutPriority = workingDaysSinceCreation
	}

	// component was not defined
	if len(bug.Fields.Components) == 0 {
		jiraBugMetric.DaysWithoutComponent = workingDaysSinceCreation
	}

	return jiraBugMetric
}

// GetAllOpenBugs gets all the bugs that are open for the given project
// func (d *Database) GetAllOpenBugs(dateFrom, dateTo string) ([]*db.Bugs, error) {
func (d *Database) GetAllOpenBugs(project string) ([]*db.Bugs, error) {
	b, err := d.client.Bugs.Query().
		Where(predicate.Bugs(bugs.Resolved(false))).
		Where(predicate.Bugs(bugs.ProjectKey(project))).
		All(context.TODO())

	if err != nil {
		return nil, convertDBError("failed to return bugs: %w", err)
	}

	return b, nil
}

// GetAllOpenBugsForSliAlerts gets all the bugs that are open for the given project
// except bugs with status as "Waiting" or "Release Pending"
func (d *Database) GetAllOpenBugsForSliAlerts(project string) ([]*db.Bugs, error) {
	b, err := d.client.Bugs.Query().
		Where(predicate.Bugs(bugs.StatusNotIn("Waiting", "Release Pending", "Closed"))).
		Where(predicate.Bugs(bugs.ProjectKey(project))).
		All(context.TODO())

	if err != nil {
		return nil, convertDBError("failed to return bugs: %w", err)
	}

	return b, nil
}

func getAssignee(user *jira.User) string {
	if user != nil {
		return user.Key
	}

	return "unassigned"
}

func getDays(createdDate, resolutionDate jira.Time, status string) string {
	firstDate := time.Time(createdDate)
	secondDate := time.Now()

	if status == "Closed" {
		secondDate = time.Time(resolutionDate)
	}

	diff := secondDate.Sub(firstDate).Hours()

	// diff in days
	diff = diff / 24

	// round diff to 2 decimal places
	diff = math.Round(diff*100) / 100

	return fmt.Sprintf("%.2f", diff)
}

// GetOpenBugsAffectingCI gets all the open issues with 'ci-fail' as label
func (d *Database) GetOpenBugsAffectingCI(t *db.Teams) ([]*db.Bugs, error) {
	bugs, err := d.client.Teams.QueryBugs(t).
		Where(predicate.Bugs(bugs.Resolved(false))).
		Where(predicate.Bugs(bugs.LabelsContains("ci-fail"))).
		All(context.TODO())

	if err != nil {
		return nil, convertDBError("failed to return bugs: %w", err)
	}

	return bugs, nil
}
