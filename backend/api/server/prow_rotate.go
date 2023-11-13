package server

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/redhat-appstudio/quality-studio/api/apis/prow/v1alpha1"
	prowV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/prow/v1alpha1"
	"github.com/redhat-appstudio/quality-studio/api/server/router/prow"

	"github.com/redhat-appstudio/quality-studio/pkg/connectors/gcs"
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
		s.cfg.Logger.Info("updating repository", zap.String("repo_name", repo.RepositoryName), zap.String("git_org", repo.GitOrganization))
		filter := FilterProwJobsByRepository(repo.GitOrganization, repo.RepositoryName, prowjobs)
		s.SaveProwJobsInDatabase(filter, repo)
	}
}

func (s *Server) SaveProwJobsInDatabase(prowJobs []prow.ProwJob, repo *db.Repository) {
	for _, job := range prowJobs {
		prowJobsInDatabase, _ := s.cfg.Storage.GetProwJobsResultsByJobID(job.Status.BuildID)

		if len(prowJobsInDatabase) > 0 {
			s.cfg.Logger.Sugar().Debugf("Data already exists in database about jobID %v, %v", job.Status.BuildID)
			continue
		}

		diff := job.Status.CompletionTime.Sub(job.Status.StartTime)

		if err := s.cfg.Storage.CreateProwJobResults(v1alpha1.Job{
			JobID:     job.Status.BuildID,
			CreatedAt: job.Status.StartTime,
			State:     string(job.Status.State),
			JobType:   job.Spec.Type,
			JobName:   job.Spec.Job,
			JobURL:    job.Status.URL,
			Duration:  float64(diff),
		}, repo.ID); err != nil {
			s.cfg.Logger.Sugar().Errorf("failed to save prowJob", err)
		}
		if job.Spec.Type != "periodic" {
			suitesXml := prow.TestSuites{}

			g := gcs.BucketHandleClient()
			suites := g.GetJobJunitContent(repo.GitOrganization, repo.RepositoryName, ExtractPullRequestNumberFromLabels(job.Labels), job.Status.BuildID, job.Spec.Job, "e2e-report.xml")

			if len(suites) > 0 {
				fmt.Println(job.Status.State, job.Status.BuildID)
				if err := xml.Unmarshal(suites, &suitesXml); err != nil {
					s.cfg.Logger.Sugar().Warnf("Failed convert xml file to golang bytes ", err)
					continue
				}
			}

			for _, suite := range suitesXml.Suites {
				for _, testCase := range suite.TestCases {
					if testCase.Status == "failed" {
						jobSuite := prowV1Alpha1.JobSuites{
							JobID:          job.Status.BuildID,
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
		}
	}
}

func FilterProwJobsByRepository(gitOrg string, repoName string, prowJobs []prow.ProwJob) []prow.ProwJob {
	var filtered = []prow.ProwJob{}

	for _, job := range prowJobs {

		prowOrg, prowRepo := ExtractOrgAndRepoFromProwJobLabels(job.GetLabels())

		if prowOrg == gitOrg && prowRepo == repoName && job.Status.State != prow.AbortedState && job.Status.State != prow.PendingState {
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
