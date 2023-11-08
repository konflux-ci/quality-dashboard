package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v44/github"
	"github.com/redhat-appstudio/quality-studio/pkg/logger"
	"go.uber.org/zap"
)

func (g *Github) CheckIfRepoExistsInOpenshiftCI(organization string, repository string) bool {
	if _, _, resp, err := g.client.Repositories.GetContents(context.Background(), "openshift", "release", fmt.Sprintf("ci-operator/config/%s/%s", organization, repository), &github.RepositoryContentGetOptions{
		Ref: "master",
	}); err != nil {
		logger, _ := logger.InitZap("info")
		if resp.StatusCode == 404 {
			logger.Warn(fmt.Sprintf("repository %s/%s does not exist in OpenShift CI", repository, organization))
		} else {
			logger.Error("repository does not exist in OpenShift CI", zap.String("repository", repository), zap.Error(err))
		}
		return false
	}

	return true
}

func (g *Github) GetJobTypes(organization string, repository string) []string {
	jobTypes := make([]string, 0)

	_, contents, _, err := g.client.Repositories.GetContents(context.Background(), "openshift", "release", fmt.Sprintf("ci-operator/jobs/%s/%s", organization, repository), &github.RepositoryContentGetOptions{
		Ref: "master",
	})
	if err != nil {
		logger, _ := logger.InitZap("info")
		logger.Error("Failed to get job types for repository", zap.String("repository", repository), zap.Error(err))
	}

	for i := range contents {
		fileName := *contents[i].Name

		if strings.HasSuffix(fileName, "-presubmits.yaml") {
			jobTypes = append(jobTypes, "presubmit")
		}

		if strings.HasSuffix(fileName, "-postsubmits.yaml") {
			jobTypes = append(jobTypes, "postsubmit")
		}
		if strings.HasSuffix(fileName, "-periodics.yaml") {
			jobTypes = append(jobTypes, "periodic")
		}
	}

	return jobTypes
}
