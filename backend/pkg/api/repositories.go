package api

import (
	"encoding/json"
	"net/http"

	"github.com/flacatus/qe-dashboard-backend/pkg/storage"
	"github.com/gorilla/mux"
)

type RepositoryDeleteRequest struct {
	GitOrganization string `json:"git_organization"`
	GitRepository   string `json:"repository_name"`
}

// Version godoc
// @Summary Quality Repositories
// @Description Create a new github repository to monitor
// @Tags Github Repository API
// @Produce json
// @Router /api/quality/repositories/create [post]
// @Success 200
func (s *Server) repositoriesCreateHandler(w http.ResponseWriter, r *http.Request) {
	var repository GitRepositoryRequest
	json.NewDecoder(r.Body).Decode(&repository)

	repoInfo, err := s.githubAPI.GetRepositoriesInformation(repository.GitOrganization, repository.GitRepository)
	if err != nil {
		s.logger.Sugar().Errorf("Failed to save repository %v", err)
		s.ErrorResponse(w, r, "Failed to obtain repositories. There are no repository cached", 500)

		return
	}
	repo, err := s.config.Storage.CreateRepository(storage.Repository{
		RepositoryName:  repoInfo.RepositoryName,
		GitOrganization: repository.GitOrganization,
		Description:     repoInfo.Description,
		GitURL:          repoInfo.HTMLUrl,
	})
	if err != nil {
		s.logger.Sugar().Errorf("Failed to save repository %v", err)
		s.ErrorResponse(w, r, "Failed to obtain repositories. There are no repository cached", 500)

		return
	}

	coverage, err := s.codecovAPI.GetCodeCovInfo(repo.GitOrganization, repo.RepositoryName)
	if err != nil {
		s.logger.Sugar().Errorf("Failed to save repository %v", err)
		s.ErrorResponse(w, r, "Failed to obtain repositories. There are no repository cached", 500)

		return
	}
	totalCoverageConverted, _ := coverage.Commit.Totals.TotalCoverage.Float64()
	err = s.config.Storage.CreateCoverage(storage.Coverage{
		GitOrganization:    repo.GitOrganization,
		RepositoryName:     repo.RepositoryName,
		CoveragePercentage: totalCoverageConverted,
	}, repo.ID)

	if err != nil {
		s.logger.Sugar().Errorf("Failed to save repository %v", err)
		s.ErrorResponse(w, r, "Failed to obtain repositories. There are no repository cached", 500)

		return
	}

	workflows, err := s.githubAPI.GetRepositoryWorkflows(repo.GitOrganization, repo.RepositoryName)
	for _, w := range workflows.Workflows {
		s.config.Storage.CreateWorkflows(storage.GithubWorkflows{
			WorkflowName: w.Name,
			BadgeURL:     w.BadgeURL,
			HTMLURL:      w.HTML_URL,
			JobURL:       w.JobURL,
			State:        w.State,
		}, repo.ID)
	}
	if err != nil {
		s.logger.Sugar().Errorf("Failed to save repository %v", err)
		s.ErrorResponse(w, r, "Failed to obtain repositories. There are no repository cached", 500)

		return
	}

	s.JSONResponse(w, r, repo)
}

// Version godoc
// @Summary Quality Repositories
// @Description returns all repository information founded in server configuration
// @Tags Github Repository API
// @Produce json
// @Router /api/quality/repositories/list [get]
// @Success 200
func (s *Server) listRepositoriesHandler(w http.ResponseWriter, r *http.Request) {
	// set a value with a cost of 1
	repos, err := s.config.Storage.ListRepositoriesQualityInfo()

	if err != nil {
		s.ErrorResponse(w, r, "Failed to get repositories", 500)
		return
	}
	s.JSONResponse(w, r, repos)
}

// Version godoc
// @Summary Quality Repositories
// @Description return github repository given org and repo name
// @Tags Github Repository API
// @Produce json
// @Router /api/quality/repositories/get/{gitOrg}/{repo_name} [get]
// @Success 200
func (s *Server) getRepositoryHandler(w http.ResponseWriter, r *http.Request) {
	urlParameters := mux.Vars(r)
	gitOrg := urlParameters["git_org"]
	if len(gitOrg) <= 0 {
		s.ErrorResponse(w, r, "GitOrg missing", 500)
		return
	}
	repoName := urlParameters["repo_name"]
	if len(repoName) <= 0 {
		s.ErrorResponse(w, r, "Repository name missing", 500)
		return
	}
	repository, err := s.config.Storage.GetRepository(repoName, gitOrg)
	if err != nil {
		s.ErrorResponse(w, r, "Failed to get repositories "+gitOrg+repoName, 500)
		return
	}
	s.JSONResponse(w, r, repository)
}

// Version godoc
// @Summary Quality Repositories
// @Description return github workflows from a given repository
// @Tags Github Repository API
// @Produce json
// @Router /api/quality/workflows/get [get]
// @Success 200
func (s *Server) listWorkflowsHandler(w http.ResponseWriter, r *http.Request) {
	repositoryName := r.URL.Query()["repository_name"]
	workflows, err := s.config.Storage.ListWorkflowsByRepository(repositoryName[0])
	if err != nil {
		s.ErrorResponse(w, r, "Failed to get repositories", 500)
		return
	}
	s.JSONResponse(w, r, workflows)
}

// Version godoc
// @Summary Quality Repositories
// @Description delete a given repository from a organization
// @Tags Github Repository API
// @Produce json
// @Router /api/quality/repositories/delete [delete]
// @Success 200
func (s *Server) deleteRepositoryHandler(w http.ResponseWriter, r *http.Request) {
	var repository RepositoryDeleteRequest
	json.NewDecoder(r.Body).Decode(&repository)
	if repository.GitRepository == "" {
		s.ErrorResponse(w, r, "Failed to remove repository. Field 'repository_name' missing", 400)
		return
	}
	if repository.GitOrganization == "" {
		s.ErrorResponse(w, r, "Failed to remove repository. Field 'git_organization' missing", 400)
		return
	}
	err := s.config.Storage.DeleteRepository(repository.GitRepository, repository.GitOrganization)
	if err != nil {
		s.ErrorResponse(w, r, "Failed to remove repository", 400)
		return
	}

	s.JSONResponse(w, r, SuccessResponse{
		Message: "Repository deleted",
	})
}
