package client

import (
	"github.com/redhat-appstudio/quality-studio/pkg/storage"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
)

func toStorageRepository(p *db.Repository) storage.Repository {
	return storage.Repository{
		RepositoryName:  p.RepositoryName,
		GitOrganization: p.GitOrganization,
		Description:     p.Description,
		GitURL:          p.GitURL,
		ID:              p.ID,
	}
}

func toStorageWorkflows(p *db.Workflows) storage.GithubWorkflows {
	return storage.GithubWorkflows{
		WorkflowName: p.WorkflowName,
		BadgeURL:     p.BadgeURL,
		HTMLURL:      p.HTMLURL,
		State:        p.State,
	}
}

func toStorageRepositoryAllInfo(p *db.Repository, w []*db.Workflows, c *db.CodeCov) storage.RepositoryQualityInfo {
	storageRepositories := make([]storage.GithubWorkflows, 0, len(w))
	for _, workflow := range w {
		storageRepositories = append(storageRepositories, toStorageWorkflows(workflow))
	}
	return storage.RepositoryQualityInfo{
		GitOrganization: p.GitOrganization,
		RepositoryName:  p.RepositoryName,
		GitURL:          p.GitURL,
		Description:     p.Description,
		CI:              storageRepositories,
		CodeCoverage: storage.Coverage{
			RepositoryName:     c.RepositoryName,
			GitOrganization:    c.GitOrganization,
			CoveragePercentage: c.CoveragePercentage,
		},
	}
}
