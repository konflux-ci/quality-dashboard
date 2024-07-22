package konflux

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	prowV1Alpha1 "github.com/konflux-ci/quality-dashboard/api/apis/prow/v1alpha1"
	"github.com/konflux-ci/quality-dashboard/api/server/router/prow"
	"github.com/konflux-ci/quality-dashboard/api/types"
	"github.com/konflux-ci/quality-dashboard/pkg/utils/httputils"

	"io"
	"net/http"
	"regexp"
	"strings"
)

func (k *konfluxRouter) receiveMetrics(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	var xunit prow.TestSuites

	xunitFile, xunitHead, err := r.FormFile("xunit")
	if err != nil {
		return httputils.WriteJSON(w, http.StatusBadRequest, &types.ErrorResponse{
			Message:    "failed to get xunit file from request",
			StatusCode: http.StatusBadRequest,
		})
	}

	xunitBytes, err := io.ReadAll(xunitFile)
	if err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    "failed to convert xunit file to bytes",
			StatusCode: http.StatusInternalServerError,
		})
	}

	err = xml.Unmarshal(xunitBytes, &xunit)
	if err != nil {
		return httputils.WriteJSON(w, http.StatusBadRequest, &types.ErrorResponse{
			Message:    "failed to decode the JUnit xunitFile provided in /konflux/metadata/post body",
			StatusCode: http.StatusBadRequest,
		})
	}

	metadataFile, _, err := r.FormFile("metadata")
	if err != nil {
		return httputils.WriteJSON(w, http.StatusBadRequest, &types.ErrorResponse{
			Message:    "failed to get metadata file from request",
			StatusCode: http.StatusBadRequest,
		})
	}

	metadataBytes, err := io.ReadAll(metadataFile)
	if err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    "failed to convert metadata file to bytes",
			StatusCode: http.StatusInternalServerError,
		})
	}

	var metadata KonfluxMetadata
	err = json.Unmarshal(metadataBytes, &metadata)
	if err != nil {
		return httputils.WriteJSON(w, http.StatusBadRequest, &types.ErrorResponse{
			Message:    "failed to decode the json metadata file",
			StatusCode: http.StatusBadRequest,
		})
	}

	repoInfo, err := k.storage.GetRepository(metadata.RepositoryName, metadata.GitOrganization)
	if err != nil {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    fmt.Sprintf("Repository '%s' doesn't exists in quality studio database", metadata.RepositoryName),
			StatusCode: 400,
		})
	}

	for _, suite := range xunit.Suites {
		for _, testCase := range suite.TestCases {
			if testCase.FailureOutput == nil {
				continue
			}

			jobSuite := prowV1Alpha1.JobSuites{
				JobID:                 metadata.JobId,
				JobName:               metadata.JobName,
				TestCaseName:          testCase.Name,
				TestCaseStatus:        testCase.Status,
				TestTiming:            testCase.Duration,
				JobType:               metadata.JobType,
				ErrorMessage:          testCase.FailureOutput.Message,
				JobURL:                metadata.JobUrl,
				ExternalServiceImpact: metadata.ExternalImpact,
				CreatedAt:             metadata.CreatedAt,
			}

			re := regexp.MustCompile(`\[It]\s*\[([^]]+)]`)
			matches := re.FindStringSubmatch(testCase.Name)

			if jobSuite.JobName == "pull-ci-redhat-appstudio-infra-deployments-main-appstudio-hac-e2e-tests" {
				jobSuite.SuiteName = suite.Name
			} else if strings.Contains(xunitHead.Filename, "qd-report") {
				jobSuite.SuiteName = suite.Name
				jobSuite.ErrorMessage = testCase.FailureOutput.Output
			} else if len(matches) > 1 {
				jobSuite.SuiteName = matches[1]
			}

			if jobSuite.SuiteName == "" {
				continue
			}

			err = k.storage.CreateProwJobSuites(jobSuite, repoInfo.ID)
			if err != nil {
				k.logger.Sugar().Info(err)
			}
		}
	}

	if err := k.storage.CreateProwJobResults(prowV1Alpha1.Job{
		JobID:                 metadata.JobId,
		CreatedAt:             metadata.CreatedAt,
		State:                 metadata.State,
		JobType:               metadata.JobType,
		JobName:               metadata.JobName,
		JobURL:                metadata.JobUrl,
		ExternalServiceImpact: metadata.ExternalImpact,
	}, repoInfo.ID); err != nil {
		k.logger.Sugar().Errorf("failed to save prowJob", err)
	}

	return nil
}
