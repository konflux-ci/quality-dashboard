package client

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	prV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/github/v1alpha1"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/predicate"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/pullrequests"
)

// CreatePullRequest saves provided pull request information in database.
func (d *Database) CreatePullRequest(pr prV1Alpha1.PullRequest, repo_id uuid.UUID) error {
	prAlreadyExists := d.client.PullRequests.Query().
		Where(pullrequests.Number(pr.Number)).
		Where(pullrequests.RepositoryName(pr.RepositoryName)).
		Where(pullrequests.RepositoryOrganization(pr.RepositoryOrganization)).
		ExistX(context.TODO())

	if prAlreadyExists {
		_, err := d.client.PullRequests.Update().
			Where(predicate.PullRequests(pullrequests.Number(pr.Number))).
			Where(predicate.PullRequests(pullrequests.RepositoryName(pr.RepositoryName))).
			Where(predicate.PullRequests(pullrequests.RepositoryOrganization(pr.RepositoryOrganization))).
			SetTitle(pr.Title).
			SetCreatedAt(pr.CreatedAt).
			SetMergedAt(pr.MergedAt).
			SetClosedAt(pr.ClosedAt).
			SetState(pr.State).
			SetAuthor(pr.Author).
			Save(context.TODO())
		if err != nil {
			return convertDBError("failed to create bug: %w", err)
		}
	} else {
		p, err := d.client.PullRequests.Create().
			SetRepositoryName(pr.RepositoryName).
			SetRepositoryOrganization(pr.RepositoryOrganization).
			SetNumber(pr.Number).
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
	}

	return nil
}

// GetPullRequestsByRepository gets the summary and the metrics of the open and merged pull requests.
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

			fmt.Println(prs)

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

	mergeAvg := float64(0)

	if totalMergedPrs != 0 {
		totalMergeTime = totalMergeTime / float64(totalMergedPrs)
		mergeAvg = math.Round(totalMergeTime*100) / 100
	}

	totalMergeTime = totalMergeTime / float64(totalMergedPrs)

	info.Summary = prV1Alpha1.Summary{
		MergedPrsCount: totalMergedPrs,
		OpenPrsCount:   totalOpenPrs,
		MergeAvg:       mergeAvg,
	}
	info.Metrics = metrics

	return info, nil
}
