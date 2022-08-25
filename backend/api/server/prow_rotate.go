package server

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/redhat-appstudio/quality-studio/api/server/router/prow"
	"github.com/redhat-appstudio/quality-studio/pkg/storage"
)

const (
	RedHatAppStudioOrg = "redhat-appstudio"
	ProwEndpoint       = "https://prow.ci.openshift.org/prowjobs.js"
)

func (s *Server) ProwStaticUpdate() {
	jobsJSON, err := fetchJobsJSON("https://prow.ci.openshift.org/prowjobs.js")
	if err != nil {
		s.cfg.Logger.Sugar().Warnf("Failed to fetch prow endpoint ", err)
	}
	prowjobs, err := jobsJSONToProwJobs(jobsJSON)
	if err != nil {
		s.cfg.Logger.Sugar().Warnf("Failed to convert prow json endpoint to struct ", err)
	}

	storageRepos, err := s.cfg.Storage.ListRepositories()
	if err != nil {
		s.cfg.Logger.Sugar().Warnf("Failed to get repos from database ", err)
	}

	for _, repo := range storageRepos {
		for _, pj := range prowjobs {
			suitesXml := prow.TestSuites{}
			prowOrg, prowRepo := ExtractOrgAndRepoFromProwJobLabels(pj.Labels)

			if prowOrg == repo.GitOrganization && prowRepo == repo.RepositoryName && pj.Status.State != prow.AbortedState && pj.Status.State != prow.PendingState && !strings.Contains(pj.Status.URL, "-images") && !strings.Contains(pj.Status.URL, "-index") {
				// check if job already in database
				prowJobsInDatabase, _ := s.cfg.Storage.GetProwJobsResultsByJobID(pj.Status.BuildID)
				if len(prowJobsInDatabase) > 0 {
					s.cfg.Logger.Sugar().Debugf("Data already exist in database about jobID %v, %v", pj.Status.BuildID)
					continue
				}

				pj.Status.URL = strings.Replace(pj.Status.URL, "https://prow.ci.openshift.org/view/gs", "https://gcsweb-ci.apps.ci.l2s4.p1.openshiftapps.com/gcs", 1)
				matches := RegexpCompiler.FindStringSubmatch(pj.Status.URL)

				if !strings.Contains(pj.Status.URL, "-main-images") && pj.Status.State != prow.ErrorState && prowOrg == "redhat-appstudio" && len(matches) == 2 {
					// convert the url to get GCS url where are stored the artifacts for appstudio
					suites, _ := fetchSuitesXml(pj.Status.URL + "/" + "artifacts/" + matches[2] + "/" + matches[2] + "/artifacts/e2e-report.xml")
					// we unmarshal our byteArray which contains our
					// xmlFiles content into 'suitesXml' which we defined above
					if err := xml.Unmarshal(suites, &suitesXml); err != nil {
						s.cfg.Logger.Sugar().Warnf("Failed convert xml file to golang bytes ", err)
					}
				}
				if err := SaveProwJobsinDatabase(s.cfg.Storage, pj, suitesXml, repo.ID); err != nil {
					s.cfg.Logger.Sugar().Error("Failed to save job database ", err)
				}
			}
		}
	}
}

func SaveProwJobsinDatabase(s storage.Storage, pj prow.ProwJob, ts prow.TestSuites, repositoryId uuid.UUID) error {
	prowJob := storage.ProwJobStatus{}
	duration, numTests, testFailed, testSkipped := getSuitesData(pj, ts)

	prowJob.JobID = pj.Status.BuildID
	prowJob.CreatedAt = pj.Status.StartTime
	prowJob.Duration = duration
	prowJob.TestsCount = int64(numTests)
	prowJob.FailedCount = int64(testFailed)
	prowJob.SkippedCount = int64(testSkipped)
	prowJob.JobType = pj.Spec.Type
	prowJob.JobName = pj.GetAnnotations()["prow.k8s.io/job"]
	prowJob.State = string(pj.Status.State)
	prowJob.JobURL = pj.Status.URL
	prowJob.CIFailed = getCIFailed(pj, ts)

	if err := s.CreateProwJobResults(prowJob, repositoryId); err != nil {
		return fmt.Errorf("Failed to save job to db %s", err)
	}

	if pj.Spec.Type == "periodic" {
		for _, suite := range ts.Suites {
			for _, testCase := range suite.TestCases {
				s.CreateProwJobSuites(storage.ProwJobSuites{
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

func getSuitesData(pj prow.ProwJob, ts prow.TestSuites) (float64, int, int, int) {
	var duration float64 = 0
	var numTests, testFailed, testSkipped = 0, 0, 0
	if pj.Spec.Type != "periodic" {
		return duration, numTests, testFailed, testSkipped
	}
	for _, suite := range ts.Suites {
		duration = suite.Duration
		numTests = int(suite.NumTests)
		testFailed = int(suite.NumFailed)
		testSkipped = int(suite.NumSkipped)
	}

	return duration, numTests, testFailed, testSkipped
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
