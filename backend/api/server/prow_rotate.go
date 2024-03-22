package server

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	prowV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/prow/v1alpha1"
	"github.com/redhat-appstudio/quality-studio/api/server/router/prow"

	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
	"go.uber.org/zap"
)

const (
	RedHatAppStudioOrg = "redhat-appstudio"
	ProwEndpoint       = "https://prow.ci.openshift.org/prowjobs.js"
)

var JunitRegexpSearch = regexp.MustCompile(`(j?unit|e2e|qd_report_)-?[0-9a-z]+\.xml`)
var ExternalServicesSearch = regexp.MustCompile(`services-status.json`)

func (s *Server) UpdateProwStatusByTeam() {
	prowjobs, err := fetchJobsJSON(ProwEndpoint)
	if err != nil {
		s.cfg.Logger.Sugar().Warnf("failed to fetch Prow jobs", err)
	}

	repos, err := s.cfg.Storage.ListAllRepositories()
	if err != nil {
		s.cfg.Logger.Sugar().Warnf("failed to fetch GitHub repo from storage", err)
	}

	for _, repo := range repos {
		filter := FilterProwJobsByRepository(repo.GitOrganization, repo.RepositoryName, prowjobs)
		s.SaveProwJobsInDatabase(filter, repo)
	}

	s.ErrorsUpdateInProwJobs()
}

func (s *Server) SaveProwJobsInDatabase(prowJobs []prow.ProwJob, repo *db.Repository) {
	for _, job := range prowJobs {
		suitesXml := prow.TestSuites{}
		prowJobsInDatabase, _ := s.cfg.Storage.GetProwJobsResultsByJobID(job.Status.BuildID)

		if len(prowJobsInDatabase) > 0 {
			s.cfg.Logger.Sugar().Debugf("Data already exists in database about jobID %v, %v", job.Status.BuildID)
			continue
		}

		diff := job.Status.CompletionTime.Sub(job.Status.StartTime)

		var suites []byte
		var xmlFileName string

		if job.Spec.Type == "periodic" {
			suites, xmlFileName = s.cfg.GCS.GetJobJunitContent("", "", "", job.Status.BuildID, job.Spec.Type, job.Spec.Job, JunitRegexpSearch)
		} else if job.Spec.Type == "presubmit" {
			suites, xmlFileName = s.cfg.GCS.GetJobJunitContent(repo.GitOrganization, repo.RepositoryName, ExtractPullRequestNumberFromLabels(job.Labels),
				job.Status.BuildID, job.Spec.Type, job.Spec.Job, JunitRegexpSearch)
		}

		if len(suites) > 0 {
			if err := xml.Unmarshal(suites, &suitesXml); err != nil {
				s.cfg.Logger.Sugar().Warnf("Failed convert xml file to golang bytes ", err)
				continue
			}
		}

		s.cfg.Logger.Info("successfully pulled data from gcs. Updating database", zap.String("job_name", job.Spec.Job), zap.String("job_id", job.Status.BuildID),
			zap.String("git_org", repo.GitOrganization), zap.String("repo_name", repo.RepositoryName), zap.String("job_status", string(job.Status.State)))

		isImpacted, _ := s.CheckIfJobIsImpactedByExternalServices(job, repo)

		totalErrorMessages := ""
		for _, suite := range suitesXml.Suites {
			for _, testCase := range suite.TestCases {
				// hac prow jobs does not have status field in testCase
				if testCase.FailureOutput != nil {
					errorMsg := testCase.Name + ":\n" + testCase.FailureOutput.Message + "\n\n"
					totalErrorMessages += errorMsg

					jobSuite := prowV1Alpha1.JobSuites{
						JobID:                 job.Status.BuildID,
						JobName:               job.Spec.Job,
						TestCaseName:          testCase.Name,
						TestCaseStatus:        "failed",
						TestTiming:            testCase.Duration,
						JobType:               job.Spec.Type,
						CreatedAt:             job.Status.StartTime,
						ErrorMessage:          testCase.FailureOutput.Message,
						JobURL:                job.Status.URL,
						ExternalServiceImpact: isImpacted,
					}

					re := regexp.MustCompile(`\[It\]\s*\[([^\]]+)\]`)
					matches := re.FindStringSubmatch(testCase.Name)

					if jobSuite.JobName == "pull-ci-redhat-appstudio-infra-deployments-main-appstudio-hac-e2e-tests" {
						jobSuite.SuiteName = suite.Name
					} else if strings.Contains(xmlFileName, "qd_report") {
						jobSuite.SuiteName = suite.Name
						jobSuite.ErrorMessage = testCase.FailureOutput.Output
					} else if len(matches) > 1 {
						jobSuite.SuiteName = matches[1]
					}

					if jobSuite.SuiteName != "" {
						if err := s.cfg.Storage.CreateProwJobSuites(jobSuite, repo.ID); err != nil {
							// nolint:all
							s.cfg.Logger.Sugar().Info(err)
						}
					}
				}
			}
		}

		buildErrorLogs := ""
		// means that job failed probably because of some infra issue
		if totalErrorMessages == "" && job.Status.State == prow.FailureState {
			buildErrorLogs = getBuildLogErrors(job.Status.URL)
		}

		if err := s.cfg.Storage.CreateProwJobResults(prowV1Alpha1.Job{
			JobID:                 job.Status.BuildID,
			CreatedAt:             job.Status.StartTime,
			State:                 string(job.Status.State),
			JobType:               job.Spec.Type,
			JobName:               job.Spec.Job,
			JobURL:                job.Status.URL,
			Duration:              float64(diff),
			ExternalServiceImpact: isImpacted,
			// E2EFailedTestMessages and BuildErrorLogs are used to get the impact of RHTAPBUGS
			E2EFailedTestMessages: totalErrorMessages,
			BuildErrorLogs:        buildErrorLogs,
		}, repo.ID); err != nil {
			s.cfg.Logger.Sugar().Errorf("failed to save prowJob", err)
		}
	}
}

func (s *Server) CheckIfJobIsImpactedByExternalServices(job prow.ProwJob, repo *db.Repository) (bool, error) {
	var externalService HealthCheckStatus
	var externalByteContent []byte

	if job.Spec.Type == "periodic" {
		externalByteContent, _ = s.cfg.GCS.GetJobJunitContent("", "", "", job.Status.BuildID, job.Spec.Type, job.Spec.Job, ExternalServicesSearch)
	} else if job.Spec.Type == "presubmit" {
		externalByteContent, _ = s.cfg.GCS.GetJobJunitContent(repo.GitOrganization, repo.RepositoryName, ExtractPullRequestNumberFromLabels(job.Labels),
			job.Status.BuildID, job.Spec.Type, job.Spec.Job, ExternalServicesSearch)
	}

	if externalByteContent == nil {
		return false, nil
	}

	if err := json.Unmarshal(externalByteContent, &externalService); err != nil {
		return false, fmt.Errorf("failed to parse external service status %v", err)
	}

	for _, service := range externalService.ExternalServices {
		if service.CurrentStatus.Status.Indicator != "none" {
			return true, nil
		}
	}

	return false, nil
}

func FilterProwJobsByRepository(gitOrg string, repoName string, prowJobs []prow.ProwJob) []prow.ProwJob {
	var filtered = []prow.ProwJob{}

	for _, job := range prowJobs {

		prowOrg, prowRepo := ExtractOrgAndRepoFromProwJobLabels(job.GetLabels())

		if prowOrg == gitOrg && prowRepo == repoName && job.Status.State != prow.AbortedState && job.Status.State != prow.PendingState && job.Status.State != prow.TriggeredState {
			filtered = append(filtered, job)
		}
	}

	return filtered
}

func ExtractOrgAndRepoFromProwJobLabels(labels map[string]string) (org string, repo string) {
	return labels["prow.k8s.io/refs.org"], labels["prow.k8s.io/refs.repo"]
}

func ExtractPullRequestNumberFromLabels(labels map[string]string) (number string) {
	return labels["prow.k8s.io/refs.pull"]
}

func fetchJobsJSON(prowURL string) ([]prow.ProwJob, error) {
	resp, err := http.Get(prowURL) // #nosec G107
	if err != nil {
		return nil, err
	}
	bodyBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to parse prow body %v", err)
	}

	results := make(map[string][]prow.ProwJob)
	if err := json.Unmarshal(bodyBytes, &results); err != nil {
		return nil, err
	}

	return results["items"], nil
}

func fetchSuitesXml(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(resp.Body)
}
func getBuildLogErrors(url string) string {
	buildErrorLogs := ""
	url = strings.Replace(url, "https://prow.ci.openshift.org/view/gs", "https://gcsweb-ci.apps.ci.l2s4.p1.openshiftapps.com/gcs", 1)
	buildFileUrl := url + "/" + "build-log.txt"
	buildContent, _ := fetchSuitesXml(buildFileUrl)
	buildErrorLogs = string(buildContent)

	// keep last 50 lines
	lines := strings.Split(buildErrorLogs, "\n")
	if len(lines) > 50 {
		lastIdx := len(lines) - 1
		firstIdx := lastIdx - 50
		lines = lines[firstIdx:lastIdx]
		buildErrorLogs = strings.Join(lines, "\n")
	}

	return buildErrorLogs
}

// temporary func to update build log error messages and error messages from "old" hac prow jobs
// since we were not collecting the junit-(...).xml for hac jobs
func (s *Server) ErrorsUpdateInProwJobs() {
	repos, err := s.cfg.Storage.ListAllRepositories()
	if err != nil {
		s.cfg.Logger.Sugar().Warnf("failed to fetch GitHub repo from storage", err)
	}

	for _, repo := range repos {
		prowJobs, err := s.cfg.Storage.ListFailedProwJobsByRepository(repo)
		if err != nil {
			s.cfg.Logger.Sugar().Warnf("failed to get prow jobs", err)
		}

		for _, pj := range prowJobs {
			errorMessages := ""
			// to cover the cases where we were not collecting e2e failed error messages like in hac jobs
			var suites []byte
			if pj.JobType == "periodic" {
				suites, _ = s.cfg.GCS.GetJobJunitContent("", "", "", pj.JobID, pj.JobType, pj.JobName, JunitRegexpSearch)
			} else if pj.JobType == "presubmit" {
				prNumber := getPRNumber(pj.JobURL)

				if prNumber != "" {
					suites, _ = s.cfg.GCS.GetJobJunitContent("redhat-appstudio", "infra-deployments", prNumber,
						pj.JobID, pj.JobType, pj.JobName, JunitRegexpSearch)
				} else {
					s.cfg.Logger.Sugar().Debug("Failed to get pr number from ", pj.JobURL)
				}
			}

			suitesXml := prow.TestSuites{}
			if len(suites) > 0 {
				if err := xml.Unmarshal(suites, &suitesXml); err != nil {
					s.cfg.Logger.Sugar().Warnf("Failed convert xml file to golang bytes ", err)
					continue
				} else {
					for _, suite := range suitesXml.Suites {
						errorMessages = getFailureMessages(suite)

						for _, testCase := range suite.TestCases {
							if testCase.FailureOutput != nil {

								jobSuite := prowV1Alpha1.JobSuites{
									JobID:          pj.JobID,
									JobName:        pj.JobName,
									TestCaseName:   testCase.Name,
									TestCaseStatus: "failed",
									TestTiming:     suite.Duration,
									JobType:        pj.JobType,
									CreatedAt:      pj.CreatedAt,
									ErrorMessage:   testCase.FailureOutput.Message,
									JobURL:         pj.JobURL,
								}

								re := regexp.MustCompile(`\[It\]\s*\[([^\]]+)\]`)
								matches := re.FindStringSubmatch(suite.Name)

								if jobSuite.JobName == "pull-ci-redhat-appstudio-infra-deployments-main-appstudio-hac-e2e-tests" {
									jobSuite.SuiteName = suite.Name
								} else if len(matches) > 1 {
									jobSuite.SuiteName = matches[1]
								}

								if jobSuite.SuiteName != "" {
									if err := s.cfg.Storage.CreateProwJobSuites(jobSuite, repo.ID); err != nil {
										// nolint:all
										s.cfg.Logger.Sugar().Info(err)
									}
								}
							}
						}
					}

					if errorMessages != "" {
						if err := s.cfg.Storage.UpdateErrorMessages(pj.JobID, "", errorMessages); err != nil {
							s.cfg.Logger.Sugar().Error("failed to update error messages ", err)
						}
					}
				}
			}
		}
	}
}

func getPRNumber(url string) string {
	regex := regexp.MustCompile("pull/[a-zA-Z-_0-9]*/[0-9]*")
	rawPRInfo := regex.FindString(url)
	splitPRInfo := strings.Split(rawPRInfo, "/")

	if len(splitPRInfo) == 3 {
		return splitPRInfo[2]
	}

	return ""
}

func getFailureMessages(suite *prow.TestSuite) string {
	msgs := ""

	for _, testCase := range suite.TestCases {
		if testCase.FailureOutput != nil {
			msg := testCase.Name + ":\n" + testCase.FailureOutput.Message + "\n\n"
			msgs += msg
		}
	}
	return msgs
}
