package client

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	prV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/github/v1alpha1"
)

// CreatePullRequest saves provided pull request information in database.
func (d *Database) CreatePullRequest(pr prV1Alpha1.PullRequest, repo_id uuid.UUID) error {
	p, err := d.client.PullRequests.Create().
		SetTitle(pr.Title).
		SetCreatedAt(pr.CreatedAt).
		SetMergedAt(pr.MergedAt).
		SetClosedAt(pr.ClosedAt).
		SetState(pr.State).
		SetAuthor(pr.Author).
		Save(context.TODO())
	if err != nil {
		return convertDBError("create pull request: %w", err)
	}

	_, err = d.client.Repository.UpdateOneID(repo_id).AddPrs(p).Save(context.TODO())
	if err != nil {
		return convertDBError("create pull request: %w", err)
	}

	return nil
}

// GetPullRequestsByRepository list all the pull requests of a repository.
// In addition to that, it also provides information regarding the number of open and closed pull requests.
func (d *Database) GetPullRequestsByRepository(repositoryName, organization, startDate, endDate string) (info prV1Alpha1.PullRequestsInfo, err error) {
	metrics := []prV1Alpha1.Metrics{}
	info = prV1Alpha1.PullRequestsInfo{}
	totalOpenPrs, totalMergedPrs := 0, 0
	var totalMergeTime float64

	repositories, err := d.client.Repository.Query().All(context.TODO())
	if err != nil {
		return info, convertDBError("list repositories: %w", err)
	}

	pullRequests := make([]prV1Alpha1.PullRequest, 0, len(repositories))
	for _, r := range repositories {
		if r.RepositoryName == repositoryName && r.GitOrganization == organization {
			metrics = d.getMetrics(r, startDate, endDate)
			prs, _ := d.client.Repository.QueryPrs(r).Select().
				Where(func(s *sql.Selector) { // "created_at BETWEEN ? AND 2022-08-17", "2022-08-16"
					s.Where(sql.ExprP(fmt.Sprintf("created_at BETWEEN '%s' AND '%s'", startDate, endDate)))
				}).All(context.TODO())

			for _, pr := range prs {
				if pr.State == "closed" {
					mergeTime := pr.MergedAt.Sub(pr.CreatedAt).Hours() / 24

					// avoid prs that were close but not merged
					if mergeTime > 0 {
						totalMergeTime += mergeTime
						totalMergedPrs++
					}
				} else {
					totalOpenPrs++
				}

				pullRequests = append(pullRequests, toStoragePrs(pr))
			}
		}
	}

	totalMergeTime = totalMergeTime / float64(totalMergedPrs)

	info.Summary = prV1Alpha1.Summary{
		MergedPrsCount: totalMergedPrs,
		OpenPrsCount:   totalOpenPrs,
		MergeAvg:       math.Round(totalMergeTime*100) / 100,
	}
	info.Metrics = metrics

	return info, nil
}
