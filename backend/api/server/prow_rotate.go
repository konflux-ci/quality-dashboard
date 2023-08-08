package server

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	repoV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/github/v1alpha1"
	prowV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/prow/v1alpha1"
	"github.com/redhat-appstudio/quality-studio/api/server/router/prow"
	"github.com/redhat-appstudio/quality-studio/pkg/constants"
	"github.com/redhat-appstudio/quality-studio/pkg/storage"
)

const (
	RedHatAppStudioOrg = "redhat-appstudio"
	ProwEndpoint       = "https://prow.ci.openshift.org/prowjobs.js"
)

func (s *Server) UpdateProwStatusByTeam() {
	jobsJSON, err := fetchJobsJSON(ProwEndpoint)
	if err != nil {
		s.cfg.Logger.Sugar().Warnf("Failed to fetch prow endpoint ", err)
	}
	prowjobs, err := jobsJSONToProwJobs(jobsJSON)
	if err != nil {
		s.cfg.Logger.Sugar().Warnf("Failed to convert prow json endpoint to struct ", err)
	}
	teamArr, _ := s.cfg.Storage.GetAllTeamsFromDB()

	for _, team := range teamArr {
		repo, _ := s.cfg.Storage.ListRepositories(team)
		s.ProwStaticUpdate(repo, prowjobs)
	}

	s.BuildLogErrorsUpdate()
}

// temporary func to update build log error messages from last 2 months saved prow jobs
func (s *Server) BuildLogErrorsUpdate() {
	startDate := time.Now().AddDate(0, -2, 0).Format(constants.DateFormat)
	endDate := time.Now().Format(constants.DateFormat)

	prowJobs, err := s.cfg.Storage.GetAllProwJobs(startDate, endDate)
	if err != nil {
		s.cfg.Logger.Sugar().Error("Failed to get all prow jobs ", err)
	}

	for _, pj := range prowJobs {
		if *pj.E2eFailedTestMessages == "" &&
			pj.State == string(prow.FailureState) &&
			pj.BuildErrorLogs == nil {
			buildErrors := getBuildLogErrors(pj.JobURL)
			if buildErrors != "" {
				s.cfg.Storage.UpdateBuildLogErrors(pj.JobID, buildErrors)
			}
		}
	}

}

func (s *Server) ProwStaticUpdate(storageRepos []repoV1Alpha1.Repository, prowjobs []prow.ProwJob) {
	for _, repo := range storageRepos {
		for _, pj := range prowjobs {
			suitesXml := prow.TestSuites{}
			suitesXmlUrl := ""
			prowOrg, prowRepo := ExtractOrgAndRepoFromProwJobLabels(pj.Labels)

			if prowOrg == repo.Owner.Login && prowRepo == repo.Name && pj.Status.State != prow.AbortedState && pj.Status.State != prow.PendingState && !strings.Contains(pj.Status.URL, "-images") && !strings.Contains(pj.Status.URL, "-index") {
				// check if job already in database
				prowJobsInDatabase, _ := s.cfg.Storage.GetProwJobsResultsByJobID(pj.Status.BuildID)

				if len(prowJobsInDatabase) > 0 {
					s.cfg.Logger.Sugar().Debugf("Data already exists in database about jobID %v, %v", pj.Status.BuildID)
					continue
				}

				pj.Status.URL = strings.Replace(pj.Status.URL, "https://prow.ci.openshift.org/view/gs", "https://gcsweb-ci.apps.ci.l2s4.p1.openshiftapps.com/gcs", 1)
				matches := RegexpCompiler.FindStringSubmatch(pj.Status.URL)

				if !strings.Contains(pj.Status.URL, "-main-images") && pj.Status.State != prow.ErrorState && prowOrg == "redhat-appstudio" && len(matches) > 1 {
					// convert the url to get GCS url where are stored the artifacts for appstudio
					suitesXmlUrl = pj.Status.URL + "/" + "artifacts/" + matches[2] + "/redhat-appstudio-e2e/artifacts/e2e-report.xml"
					suites, _ := fetchSuitesXml(suitesXmlUrl)
					// we unmarshal our byteArray which contains our
					// xmlFiles content into 'suitesXml' which we defined above
					if err := xml.Unmarshal(suites, &suitesXml); err != nil {
						s.cfg.Logger.Sugar().Warnf("Failed convert xml file to golang bytes ", err)
					}
				}

				if err := SaveProwJobsinDatabase(s.cfg.Storage, pj, suitesXml, repo.ID, suitesXmlUrl); err != nil {
					s.cfg.Logger.Sugar().Error("Failed to save job database ", err)
				}
			}
		}
	}
}

func getBuildLogErrors(url string) string {
	buildErrorLogs := ""
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

func SaveProwJobsinDatabase(s storage.Storage, pj prow.ProwJob, ts prow.TestSuites, repositoryId, suitesXmlUrl string) error {
	prowJob := prowV1Alpha1.Job{}
	testSuiteSummary := getSuitesData(pj, ts)
	buildErrorLogs := ""

	if testSuiteSummary.E2EFailedMessages == "" && pj.Status.State == prow.FailureState {
		buildErrorLogs = getBuildLogErrors(pj.Status.URL)
	}

	prowJob.JobID = pj.Status.BuildID
	prowJob.CreatedAt = pj.Status.StartTime
	prowJob.Duration = testSuiteSummary.Duration
	prowJob.TestsCount = int64(testSuiteSummary.NumTests)
	prowJob.FailedCount = int64(testSuiteSummary.TestFailed)
	prowJob.SkippedCount = int64(testSuiteSummary.TestSkipped)
	prowJob.JobType = pj.Spec.Type
	prowJob.JobName = pj.GetAnnotations()["prow.k8s.io/job"]
	prowJob.State = string(pj.Status.State)
	prowJob.JobURL = pj.Status.URL
	prowJob.CIFailed = getCIFailed(pj, ts)
	prowJob.E2EFailedTestMessages = testSuiteSummary.E2EFailedMessages
	prowJob.SuitesXmlUrl = suitesXmlUrl
	prowJob.BuildErrorLogs = buildErrorLogs

	if err := s.CreateProwJobResults(prowJob, repositoryId); err != nil {
		return fmt.Errorf("failed to save job to db %s", err)
	}

	if pj.Spec.Type == "periodic" {
		for _, suite := range ts.Suites {
			for _, testCase := range suite.TestCases {
				s.CreateProwJobSuites(prowV1Alpha1.JobSuites{
					JobID:          pj.Status.BuildID,
					TestCaseName:   testCase.Name,
					TestCaseStatus: testCase.Status,
					TestTiming:     testCase.Duration,
					JobType:        pj.Spec.Type,
				}, repositoryId)
			}
		}
	}

	return nil
}

func getCIFailed(p prow.ProwJob, s prow.TestSuites) int16 {
	if p.Status.State == prow.ErrorState {
		return 1
	}
	return 0
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

type TestSuiteSummary struct {
	Duration          float64
	NumTests          int
	TestFailed        int
	TestSkipped       int
	E2EFailedMessages string
}

func getSuitesData(pj prow.ProwJob, ts prow.TestSuites) TestSuiteSummary {
	testSuiteSummary := TestSuiteSummary{}

	for _, suite := range ts.Suites {
		testSuiteSummary.Duration = suite.Duration
		testSuiteSummary.NumTests = int(suite.NumTests)
		testSuiteSummary.TestFailed = int(suite.NumFailed)
		testSuiteSummary.TestSkipped = int(suite.NumSkipped)
		testSuiteSummary.E2EFailedMessages = getFailureMessages(suite)
	}

	return testSuiteSummary
}

func ExtractOrgAndRepoFromProwJobLabels(labels map[string]string) (org string, repo string) {
	return labels["prow.k8s.io/refs.org"], labels["prow.k8s.io/refs.repo"]
}

func jobsJSONToProwJobs(jobJSON []byte) ([]prow.ProwJob, error) {
	results := make(map[string][]prow.ProwJob)
	if err := json.Unmarshal(jobJSON, &results); err != nil {
		return nil, err
	}
	return results["items"], nil
}

func fetchJobsJSON(prowURL string) ([]byte, error) {
	resp, err := http.Get(prowURL) // #nosec G107
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

func fetchSuitesXml(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

func GetTestCaseByDefault(test int) int {
	if test == 0 {
		return 0
	}
	return test
}
