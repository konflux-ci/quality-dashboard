package jira

import (
	"log"

	"github.com/andygrunwald/go-jira"
)

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
		log.Fatal(err)
	}
	return issues
}
