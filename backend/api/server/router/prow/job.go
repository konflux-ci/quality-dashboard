package prow

import (
	"context"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/redhat-appstudio/quality-studio/api/types"
	"github.com/redhat-appstudio/quality-studio/pkg/storage"
	"github.com/redhat-appstudio/quality-studio/pkg/utils/httputils"
	"go.uber.org/zap"
)

type GitRepositoryRequest struct {
	GitOrganization string `json:"git_organization"`
	GitRepository   string `json:"repository_name"`
}

var (
	suitesXml TestXml
)

// version godoc
// @Summary Prow Jobs info
// @Description returns all prow jobs information stored in database
// @Tags Prow Jobs info
// @Accept json
// @Produce json
// @Param repository_name body GitRepositoryRequest true "repository name"
// @Param git_organization body GitRepositoryRequest true "repository name"
// @Param job_id body string true "repository name"
// @Router /prow/results/post [post]
// @Success 200 {Object} types.SuccessResponse
// @Failure 400 {object} types.ErrorResponse
func (s *jobRouter) createProwCIResults(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	repositoryName := r.URL.Query()["repository_name"]
	gitOrgazanitation := r.URL.Query()["git_organization"]
	jobID := r.URL.Query()["job_id"]

	if len(repositoryName) == 0 {
		return httputils.WriteJSON(w, http.StatusOK, types.ErrorResponse{
			Message:    "repository_name value not present in query",
			StatusCode: 400,
		})
	} else if len(gitOrgazanitation) == 0 {
		return httputils.WriteJSON(w, http.StatusOK, types.ErrorResponse{
			Message:    "git_organization value not present in query",
			StatusCode: 400,
		})
	} else if len(jobID) == 0 {
		return httputils.WriteJSON(w, http.StatusOK, types.ErrorResponse{
			Message:    "git_organization value not present in query",
			StatusCode: 400,
		})
	}

	prowJobsInDatabase, _ := s.Storage.GetProwJobsResultsByJobID(jobID[0])
	if len(prowJobsInDatabase) > 0 {
		return httputils.WriteJSON(w, http.StatusOK, types.ErrorResponse{
			Message:    "ci jobid already exist in database",
			StatusCode: 400,
		})
	}

	repoInfo, _ := s.Storage.GetRepository(repositoryName[0], gitOrgazanitation[0])

	s.Logger.Sugar().Info(repoInfo)

	testXml, _ := parseFileFromRequest(r, &s.Logger)

	for _, suite := range testXml.TestSuites.TestSuite {
		s.Storage.CreateProwJobResults(storage.ProwJob{
			JobID:          jobID[0],
			TestCaseName:   suite.Name,
			TestCaseStatus: suite.Status,
			TestTiming:     suite.Time,
		}, repoInfo.ID)
	}

	return httputils.WriteJSON(w, http.StatusOK, types.SuccessResponse{
		Message:    "Successfully stored Prow Job",
		StatusCode: http.StatusCreated,
	})
}

func parseFileFromRequest(r *http.Request, logger *zap.Logger) (TestXml, error) {
	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `file`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("file")
	if err != nil {
		logger.Sugar().Infof("Failed to get file from header %s", err)
		return suitesXml, err
	}

	// copy example
	f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Sugar().Infof("Failed to create temporary file to archive the tests results %s", err)
		return suitesXml, err
	}

	io.Copy(f, file)
	defer file.Close()

	xmlFile, err := os.Open(f.Name())
	// if we os.Open returns an error then handle it
	if err != nil {
		logger.Sugar().Infof("Failed to open xml file %s", err)
		return suitesXml, err
	}

	// defer the closing of our xmlFile so that we can parse it later on
	defer xmlFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(xmlFile)

	// Remove the file after reading
	os.Remove(handler.Filename)

	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'users' which we defined above
	if err := xml.Unmarshal(byteValue, &suitesXml); err != nil {
		if err != nil {
			logger.Sugar().Infof("Failed convert xml file to golang bytes %s", err)
			return suitesXml, err
		}
	}

	return suitesXml, nil
}

// version godoc
// @Summary Prow Jobs info
// @Description returns all prow jobs related to git_organization and repository_name
// @Tags Prow Jobs info
// @Accept json
// @Produce json
// @Param repository body GitRepositoryRequest true "repository name"
// @Param organization body GitRepositoryRequest true "git_organization"
// @Router /prow/results/get [get]
// @Success 200 {Object} []db.Prow
// @Failure 400 {object} types.ErrorResponse
func (s *jobRouter) getProwJobs(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	repositoryName := r.URL.Query()["repository_name"]
	gitOrgazanitation := r.URL.Query()["git_organization"]

	if len(repositoryName) == 0 {
		return httputils.WriteJSON(w, http.StatusOK, types.ErrorResponse{
			Message:    "repository_name value not present in query",
			StatusCode: 400,
		})
	} else if len(gitOrgazanitation) == 0 {
		return httputils.WriteJSON(w, http.StatusOK, types.ErrorResponse{
			Message:    "git_organization value not present in query",
			StatusCode: 400,
		})
	}

	repoInfo, err := s.Storage.GetRepository(repositoryName[0], gitOrgazanitation[0])
	if err != nil {
		return httputils.WriteJSON(w, http.StatusOK, types.ErrorResponse{
			Message:    "failed to get repository from database; check if repository exist in quality studio",
			StatusCode: 400,
		})
	}

	prows, err := s.Storage.GetProwJobsResults(repoInfo)
	if err != nil {
		return httputils.WriteJSON(w, http.StatusOK, types.ErrorResponse{
			Message:    "failed to get repository from database; check if repository exist in quality studio",
			StatusCode: 400,
		})
	}

	return httputils.WriteJSON(w, http.StatusOK, prows)
}
