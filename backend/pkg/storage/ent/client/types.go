package client

import (
	coverageV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/codecov/v1alpha1"

	repoV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/github/v1alpha1"
	v1alphaPlugins "github.com/redhat-appstudio/quality-studio/api/apis/plugins/v1alpha1"
	"github.com/redhat-appstudio/quality-studio/pkg/storage"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
)

func toStorageRepository(p *db.Repository) repoV1Alpha1.Repository {
	return repoV1Alpha1.Repository{
		Name: p.RepositoryName,
		Owner: repoV1Alpha1.Owner{
			Login: p.GitOrganization,
		},
		Description: p.Description,
		URL:         p.GitURL,
		ID:          p.ID,
	}
}

func installedPlugin(plugin *db.Plugins, installed bool) *v1alphaPlugins.Plugin {
	return &v1alphaPlugins.Plugin{
		Spec: v1alphaPlugins.PluginSpec{
			Name:        plugin.Name,
			Logo:        plugin.Logo,
			Category:    plugin.Category,
			Description: plugin.Description,
			Reason:      plugin.Reason,
		},
		Status: v1alphaPlugins.PluginStatus{
			Installed: installed,
		},
	}
}

func toStorageWorkflows(p *db.Workflows) repoV1Alpha1.Workflow {
	return repoV1Alpha1.Workflow{
		Name:     p.WorkflowName,
		BadgeURL: p.BadgeURL,
		HTMLURL:  p.HTMLURL,
		State:    p.State,
	}
}

func toStorageRepositoryAllInfo(p *db.Repository, c *db.CodeCov, prs repoV1Alpha1.PullRequestsInfo, workflows []repoV1Alpha1.Workflow) storage.RepositoryQualityInfo {
	covTrend := "n/a"
	if c.CoverageTrend != nil {
		covTrend = *c.CoverageTrend
	}

	return storage.RepositoryQualityInfo{
		GitOrganization: p.GitOrganization,
		RepositoryName:  p.RepositoryName,
		GitURL:          p.GitURL,
		Description:     p.Description,
		CodeCoverage: coverageV1Alpha1.Coverage{
			RepositoryName:     p.RepositoryName,
			GitOrganization:    p.GitOrganization,
			CoveragePercentage: c.CoveragePercentage,
			CoverageTrend:      covTrend,
		},
		PullRequests: prs,
		Workflows:    workflows,
	}
}
