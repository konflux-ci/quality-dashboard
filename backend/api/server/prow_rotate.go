package server

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/redhat-appstudio/quality-studio/api/apis/prow/v1alpha1"
	prowV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/prow/v1alpha1"
	"github.com/redhat-appstudio/quality-studio/api/server/router/prow"

	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
	"go.uber.org/zap"
)

const (
	RedHatAppStudioOrg = "redhat-appstudio"
	ProwEndpoint       = "https://prow.ci.openshift.org/prowjobs.js"
)

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

		s.cfg.Logger.Info("compiling job for analysis", zap.String("job_name", job.Spec.Job), zap.String("job_id", job.Status.BuildID),
			zap.String("git_org", repo.GitOrganization), zap.String("repo_name", repo.RepositoryName), zap.String("job_status", string(job.Status.State)))

		diff := job.Status.CompletionTime.Sub(job.Status.StartTime)

		var suites []byte
		if job.Spec.Type == "periodic" {
			suites = s.cfg.GCS.GetJobJunitContent("", "", "", job.Status.BuildID, job.Spec.Type, job.Spec.Job, "e2e-report.xml")
		} else if job.Spec.Type == "presubmit" {
			suites = s.cfg.GCS.GetJobJunitContent(repo.GitOrganization, repo.RepositoryName, ExtractPullRequestNumberFromLabels(job.Labels),
				job.Status.BuildID, job.Spec.Type, job.Spec.Job, "e2e-report.xml")
		}

		if len(suites) > 0 {
			if err := xml.Unmarshal(suites, &suitesXml); err != nil {
				s.cfg.Logger.Sugar().Warnf("Failed convert xml file to golang bytes ", err)
				continue
			}
		}

		s.cfg.Logger.Info("successfully pulled data from gcs. Updating database", zap.String("job_name", job.Spec.Job), zap.String("job_id", job.Status.BuildID),
			zap.String("git_org", repo.GitOrganization), zap.String("repo_name", repo.RepositoryName), zap.String("job_status", string(job.Status.State)))

		totalErrorMessages := ""
		for _, suite := range suitesXml.Suites {
			for _, testCase := range suite.TestCases {
				if testCase.Status == "failed" {
					errorMsg := testCase.Name + ":\n" + testCase.FailureOutput.Message + "\n\n"
					totalErrorMessages += errorMsg

					jobSuite := prowV1Alpha1.JobSuites{
						JobID:          job.Status.BuildID,
						JobName:        job.Spec.Job,
						TestCaseName:   testCase.Name,
						TestCaseStatus: testCase.Status,
						TestTiming:     testCase.Duration,
						JobType:        job.Spec.Type,
						CreatedAt:      job.Status.StartTime,
					}

					re := regexp.MustCompile(`\[It\]\s*\[([^\]]+)\]`)
					matches := re.FindStringSubmatch(testCase.Name)

					if len(matches) > 1 {
						jobSuite.SuiteName = matches[1]
					}

					jobSuite.JobURL = job.Status.URL

					if testCase.Status == "failed" {
						jobSuite.ErrorMessage = testCase.FailureOutput.Message
					}

					jobSuite.ErrorMessage = testCase.FailureOutput.Message

					if err := s.cfg.Storage.CreateProwJobSuites(jobSuite, repo.ID); err != nil {
						// nolint:all
						s.cfg.Logger.Sugar().Info(err)
					}
				}
			}
		}

		buildErrorLogs := ""
		// means that job failed probably because of some infra issue
		if totalErrorMessages == "" && job.Status.State == prow.FailureState {
			buildErrorLogs = getBuildLogErrors(job.Status.URL)
		}

		if err := s.cfg.Storage.CreateProwJobResults(v1alpha1.Job{
			JobID:     job.Status.BuildID,
			CreatedAt: job.Status.StartTime,
			State:     string(job.Status.State),
			JobType:   job.Spec.Type,
			JobName:   job.Spec.Job,
			JobURL:    job.Status.URL,
			Duration:  float64(diff),
			// E2EFailedTestMessages and BuildErrorLogs are used to get the impact of RHTAPBUGS
			E2EFailedTestMessages: totalErrorMessages,
			BuildErrorLogs:        buildErrorLogs,
		}, repo.ID); err != nil {
			s.cfg.Logger.Sugar().Errorf("failed to save prowJob", err)
		}
	}
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
	return ioutil.ReadAll(resp.Body)
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

// temporary func to update build log error messages and error messages from "old" prow jobs
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
			if pj.E2eFailedTestMessages == nil {
				var suites []byte
				if pj.JobType == "periodic" {
					suites = s.cfg.GCS.GetJobJunitContent("", "", "", pj.JobID, pj.JobType, pj.JobName, "e2e-report.xml")
				} else if pj.JobType == "presubmit" {
					prNumber := getPRNumber(pj.JobURL)

					if prNumber != "" {
						suites = s.cfg.GCS.GetJobJunitContent(repo.GitOrganization, repo.RepositoryName, prNumber,
							pj.JobID, pj.JobType, pj.JobName, "e2e-report.xml")
					} else {
						s.cfg.Logger.Sugar().Warnf("Failed to get pr number from ", pj.JobURL)
					}
				}

				suitesXml := prow.TestSuites{}
				if len(suites) > 0 {
					if err := xml.Unmarshal(suites, &suitesXml); err != nil {
						s.cfg.Logger.Sugar().Warnf("Failed convert xml file to golang bytes ", err)
						continue
					}
				}

				for _, suite := range suitesXml.Suites {
					errorMessages = getFailureMessages(suite)
				}

				if err := s.cfg.Storage.UpdateErrorMessages(pj.JobID, "", errorMessages); err != nil {
					s.cfg.Logger.Sugar().Error("failed to update error messages ", err)
				}

			} else {
				errorMessages = *pj.E2eFailedTestMessages
			}

			if errorMessages == "" && pj.BuildErrorLogs == nil {
				buildErrors := getBuildLogErrors(pj.JobURL)
				if buildErrors != "" {
					if err := s.cfg.Storage.UpdateErrorMessages(pj.JobID, buildErrors, ""); err != nil {
						s.cfg.Logger.Sugar().Error("failed to update build errors ", err)
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
		if testCase.Status == "failed" {
			msg := testCase.Name + ":\n" + testCase.FailureOutput.Message + "\n\n"
			msgs += msg
		}
	}
	return msgs
}
