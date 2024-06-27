package jira

import (
	"context"
	"fmt"
	"time"

	jira "github.com/andygrunwald/go-jira"
	"github.com/konflux-ci/quality-dashboard/pkg/logger"
	"go.uber.org/zap"
)

// project= Stonesoup  and type = Bug
func (t *clientFactory) GetIssueByJQLQuery(JQLQuery string) []jira.Issue {
	var issues []jira.Issue

	// append the jira issues to []jira.Issue
	appendFunc := func(i jira.Issue) (err error) {
		issues = append(issues, i)
		return err
	}

	// In this example, we'll search for all the issues with the provided JQL filter and Print the Story Points
	err := t.Client.Issue.SearchPages(JQLQuery, nil, appendFunc)
	if err != nil {
		logger, _ := logger.InitZap("info")
		logger.Error("Failed to search pages", zap.Error(err))
	}
	return issues
}

func (t *clientFactory) GetBugsByJQLQuery(JQLQuery string) []jira.Issue {
	JQLQuery = getLastSixMonths(JQLQuery)
	var issues []jira.Issue

	// append the jira issues to []jira.Issue
	appendFunc := func(i jira.Issue) (err error) {
		issues = append(issues, i)
		return err
	}

	options := &jira.SearchOptions{
		Expand: "changelog",
	}

	// In this example, we'll search for all the issues with the provided JQL filter and Print the Story Points
	err := t.Client.Issue.SearchPages(JQLQuery, options, appendFunc)
	if err != nil {
		logger, _ := logger.InitZap("info")
		logger.Error("Failed to search pages", zap.Error(err))
	}
	return issues
}

func (t *clientFactory) IsJQLQueryValid(jqlQuery string) error {
	jqlQuery = getLastSixMonths(jqlQuery)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, _, err := t.Client.Issue.Search(jqlQuery, &jira.SearchOptions{})
	select {
	case <-ctx.Done():
		cancel()
		return nil
	default:
		// function finished sooner. very likely because jql is invalid
		return err
	}
}

func getLastSixMonths(jqlQuery string) string {
	// grab bugs created in the last 6 months
	from := time.Now().AddDate(0, -6, 0)
	to := time.Now()
	jqlQuery = fmt.Sprintf("%s AND created >= %s AND created <= %s", jqlQuery, from.Format("2006-01-02"), to.Format("2006-01-02"))

	return jqlQuery
}
