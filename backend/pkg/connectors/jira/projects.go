package jira

import (
	jira "github.com/andygrunwald/go-jira"
	"github.com/redhat-appstudio/quality-studio/pkg/logger"
	"go.uber.org/zap"
)

func (t *clientFactory) GetJiraProjects() (list *jira.ProjectList, err error) {

	// In this example, we'll search for all the issues with the provided JQL filter and Print the Story Points
	projectList, _, err := t.Client.Project.GetList()
	if err != nil {
		logger, _ := logger.InitZap("info")
		logger.Error("Failed to search pages", zap.Error(err))

		return nil, err
	}
	return projectList, nil
}
