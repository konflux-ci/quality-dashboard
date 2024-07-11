package konflux

import (
	"context"
	"encoding/json"
	"encoding/xml"
	prowV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/prow/v1alpha1"
	"github.com/redhat-appstudio/quality-studio/api/server/router/prow"
	"github.com/redhat-appstudio/quality-studio/api/types"
	"github.com/redhat-appstudio/quality-studio/pkg/utils/httputils"
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

	var metadata prowV1Alpha1.Job
	err = json.Unmarshal(metadataBytes, &metadata)
	if err != nil {
		return httputils.WriteJSON(w, http.StatusBadRequest, &types.ErrorResponse{
			Message:    "failed to decode the json metadata file",
			StatusCode: http.StatusBadRequest,
		})
	}

	totalErrorMessages := ""
	for _, suite := range xunit.Suites {
		for _, testCase := range suite.TestCases {
			if testCase.FailureOutput == nil {
				continue
			}

			errorMsg := testCase.Name + ":\n" + testCase.FailureOutput.Message + "\n\n"
			totalErrorMessages += errorMsg

			jobSuite := prowV1Alpha1.JobSuites{
				JobID:                 metadata.JobID,
				JobName:               metadata.JobName,
				TestCaseName:          testCase.Name,
				TestCaseStatus:        testCase.Status,
				TestTiming:            testCase.Duration,
				JobType:               metadata.JobType,
				CreatedAt:             metadata.CreatedAt,
				ErrorMessage:          testCase.FailureOutput.Message,
				JobURL:                metadata.JobURL,
				ExternalServiceImpact: metadata.ExternalServiceImpact,
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

			repoId := ""
			if len(r.URL.Query()["repo_id"]) > 0 {
				repoId = r.URL.Query()["repo_id"][0]
			}
			err = k.storage.CreateProwJobSuites(jobSuite, repoId)
			if err != nil {
				k.logger.Sugar().Info(err)
			}
		}
	}
	return nil
}
