package client

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	prV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/github/v1alpha1"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/predicate"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/pullrequests"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/repository"
)

// CreatePullRequest saves provided pull request information in database.
func (d *Database) CreatePullRequests(prs prV1Alpha1.PullRequests, repo_id string) error {
	var create = false
	bulkPullRequests := make([]*db.PullRequestsCreate, len(prs))

	for i, p := range prs {
		prAlreadyExists := d.client.PullRequests.Query().
			Where(pullrequests.Number(p.Number)).
			Where(pullrequests.RepositoryName(p.Repository.Name)).
			Where(pullrequests.RepositoryOrganization(p.Repository.Owner.Login)).
			ExistX(context.TODO())
		if prAlreadyExists {
			_, err := d.client.PullRequests.Update().
				Where(predicate.PullRequests(pullrequests.Number(p.Number))).
				Where(predicate.PullRequests(pullrequests.RepositoryName(p.Repository.Name))).
				Where(predicate.PullRequests(pullrequests.RepositoryOrganization(p.Repository.Owner.Login))).
				SetTitle(p.Title).
				SetCreatedAt(p.CreatedAt).
				SetMergedAt(p.MergedAt).
				SetClosedAt(p.ClosedAt).
				SetState(p.State).
				SetAuthor(p.Author.User.Login).
				Save(context.TODO())
			if err != nil {
				return convertDBError("failed to update pull request: %w", err)
			}
			continue
		}
		bulkPullRequests[i] = d.client.PullRequests.Create().
			SetTitle(p.Title).
			SetCreatedAt(p.CreatedAt).
			SetMergedAt(p.MergedAt).
			SetNumber(p.Number).
			SetClosedAt(p.ClosedAt).
			SetPrsID(repo_id).
			SetState(p.State).
			SetAuthor(p.Author.User.Login).
			SetRepositoryName(p.Repository.Name).
			SetRepositoryOrganization(p.Repository.Owner.Login)
		create = true
	}

	defer func() {
		if err := recover(); err != nil {
			// Ussually occurs when u have network issues
			fmt.Println("Internal panic ocurred, check network connection:", err)
		}
	}()

	// todo: https://github.com/ent/ent/issues/2494 Wait until we can resolve nicely conflicts in psql
	if create {
		if err := d.client.PullRequests.CreateBulk(bulkPullRequests...).OnConflict(sql.ResolveWithNewValues()).DoNothing().Exec(context.TODO()); err != nil {
			return err
		}
	}

	return nil
}

// GetPullRequestsByRepository gets the summary and the metrics of the open and merged pull requests.
func (d *Database) GetPullRequestsByRepository(repositoryName, organization, startDate, endDate string) (info prV1Alpha1.PullRequestsInfo, err error) {
	var averagePRsMergedInTimeRange float64
	totalOpenPrs, totalMergedPrs := 0, 0
	var totalMergeTime float64
	calculateDaysRange := len(getDatesBetweenRange(startDate, endDate))
	info = prV1Alpha1.PullRequestsInfo{}

	repo := d.client.Repository.Query().Where(predicate.Repository(repository.RepositoryName(repositoryName))).FirstX(context.TODO())
	totalPullRequestsMerged, err := d.client.Repository.QueryPrs(repo).
		Select(pullrequests.FieldState).
		Where(predicate.PullRequests(pullrequests.State("MERGED"))).
		All(context.TODO())

	if err != nil {
		return info, convertDBError("failed to get open pull requests: %w", err)
	}
	totalOpenPullRequests, err := d.client.Repository.QueryPrs(repo).Select(pullrequests.FieldState).
		Where(predicate.PullRequests(pullrequests.State("OPEN"))).
		All(context.TODO())

	if err != nil {
		return info, convertDBError("failed to get Opeen pull requests: %w", err)
	}
	mergedPullRequestsInTimeRange, err := d.client.Repository.QueryPrs(repo).Select(pullrequests.FieldState).
		Where(predicate.PullRequests(pullrequests.State("MERGED"))).
		Where(func(s *sql.Selector) { // "created_at BETWEEN ? AND 2022-08-17", "2022-08-16"
			s.Where(sql.ExprP(fmt.Sprintf("created_at BETWEEN '%s' AND '%s'", startDate, endDate)))
		}).
		All(context.TODO())

	if err != nil {
		return info, convertDBError("failed to get Merged pull requests in time range: %w", err)
	}

	// range between one day (same day)
	if len(getDatesBetweenRange(startDate, endDate)) == 2 && isSameDay(startDate, endDate) {
		averagePRsMergedInTimeRange = float64(len(mergedPullRequestsInTimeRange)) / 1

	} else {
		averagePRsMergedInTimeRange = float64(len(mergedPullRequestsInTimeRange)) / float64(calculateDaysRange)
	}

	if math.IsNaN(averagePRsMergedInTimeRange) {
		averagePRsMergedInTimeRange = 0
	}
	pullRequests := make([]prV1Alpha1.PullRequest, 0)
	metrics := d.getMetrics(repo, startDate, endDate)
	prs, _ := d.client.Repository.QueryPrs(repo).Select().
		Where(func(s *sql.Selector) { // "created_at BETWEEN ? AND 2022-08-17", "2022-08-16"
			s.Where(sql.ExprP(fmt.Sprintf("created_at BETWEEN '%s' AND '%s'", startDate, endDate)))
		}).All(context.TODO())

	for _, pr := range prs {
		if pr.State == "CLOSED" {
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

	info.Summary = prV1Alpha1.Summary{
		MergedPrsCount: len(totalPullRequestsMerged),
		OpenPrsCount:   len(totalOpenPullRequests),
		MergeAvg:       averagePRsMergedInTimeRange,
	}
	info.Metrics = metrics

	return info, nil
}
