package client

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	prV1Alpha1 "github.com/konflux-ci/quality-dashboard/api/apis/github/v1alpha1"
	"github.com/konflux-ci/quality-dashboard/pkg/connectors/github"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db/predicate"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db/pullrequests"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db/repository"
)

// CreatePullRequest saves provided pull request information in database.
func (d *Database) CreatePullRequests(prs prV1Alpha1.PullRequests, repo_id string) error {
	create := false
	bulkPullRequests := make([]*db.PullRequestsCreate, 0)

	for _, p := range prs {
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
				SetMergeCommit(p.MergeCommit.OID).
				SetRetestCount(github.RetestComments(&p.TimelineItems)).
				SetRetestBeforeMergeCount(github.RetestCommentsAfterLastPush(&p.TimelineItems)).
				Save(context.TODO())
			if err != nil {
				return convertDBError("failed to update pull request: %w", err)
			}
			continue
		}
		bulkPullRequest := d.client.PullRequests.Create().
			SetTitle(p.Title).
			SetCreatedAt(p.CreatedAt).
			SetMergedAt(p.MergedAt).
			SetNumber(p.Number).
			SetClosedAt(p.ClosedAt).
			SetPrsID(repo_id).
			SetState(p.State).
			SetAuthor(p.Author.User.Login).
			SetRepositoryName(p.Repository.Name).
			SetRepositoryOrganization(p.Repository.Owner.Login).
			SetMergeCommit(p.MergeCommit.OID).
			SetRetestCount(github.RetestComments(&p.TimelineItems)).
			SetRetestBeforeMergeCount(github.RetestCommentsAfterLastPush(&p.TimelineItems))
		bulkPullRequests = append(bulkPullRequests, bulkPullRequest)
		create = true
	}

	defer func() {
		if err := recover(); err != nil {
			// Usually occurs when u have network issues
			fmt.Println("Internal panic occurred, check network connection:", err)
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
	totalMergedPrs := 0
	var totalMergeTime, totalRetests, retestAvg, totalRetestBeforeMerge, retestBeforeMergeAvg float64
	info = prV1Alpha1.PullRequestsInfo{}

	repo := d.client.Repository.Query().Where(predicate.Repository(repository.RepositoryName(repositoryName))).FirstX(context.TODO())
	totalPullRequestsMerged, err := d.client.Repository.QueryPrs(repo).
		Select(pullrequests.FieldState).
		Where(predicate.PullRequests(pullrequests.State("MERGED"))).
		All(context.TODO())
	if err != nil {
		return info, convertDBError("failed to get merged pull requests: %w", err)
	}

	totalOpenPullRequests, err := d.client.Repository.QueryPrs(repo).
		Select().
		Where(predicate.PullRequests(pullrequests.State("OPEN"))).
		All(context.TODO())
	if err != nil {
		return info, convertDBError("failed to get open pull requests: %w", err)
	}

	openPullRequestsInTimeRange, _ := d.client.Repository.QueryPrs(repo).Select().
		Where(func(s *sql.Selector) {
			// 0001-01-01 00:00:00 is the default value meaning the pull requests has not yet been closed.
			s.Where(sql.ExprP(fmt.Sprintf("created_at <= '%s' AND (closed_at >= '%s' OR closed_at='0001-01-01 00:00:00')", endDate, startDate)))
		}).
		All(context.TODO())
	if err != nil {
		return info, convertDBError("failed to get open pull requests in time range: %w", err)
	}

	for _, openPr := range openPullRequestsInTimeRange {
		if openPr.RetestCount != nil {
			totalRetests += *openPr.RetestCount
		}
	}

	if len(openPullRequestsInTimeRange) != 0 {
		retestAvg = totalRetests / float64(len(openPullRequestsInTimeRange))
	}

	createdPullRequestsInTimeRange, err := d.client.Repository.QueryPrs(repo).Select().
		Where(func(s *sql.Selector) { // "created_at BETWEEN ? AND 2022-08-17", "2022-08-16"
			s.Where(sql.ExprP(fmt.Sprintf("created_at BETWEEN '%s' AND '%s'", startDate, endDate)))
		}).
		All(context.TODO())
	if err != nil {
		return info, convertDBError("failed to get merged pull requests in time range: %w", err)
	}

	mergedPullRequestsInTimeRange, err := d.client.Repository.QueryPrs(repo).Select().
		Where(predicate.PullRequests(pullrequests.State("MERGED"))).
		Where(func(s *sql.Selector) { // "merged_at BETWEEN ? AND 2022-08-17", "2022-08-16"
			s.Where(sql.ExprP(fmt.Sprintf("merged_at BETWEEN '%s' AND '%s'", startDate, endDate)))
		}).
		All(context.TODO())
	if err != nil {
		return info, convertDBError("failed to get merged pull requests in time range: %w", err)
	}

	for _, merged := range mergedPullRequestsInTimeRange {
		mergeTime := merged.MergedAt.Sub(merged.CreatedAt).Hours() / 24
		var retestBeforeMergeCount float64 = 0
		if merged.RetestBeforeMergeCount != nil {
			retestBeforeMergeCount = *merged.RetestBeforeMergeCount
		}

		if mergeTime > 0 {
			totalMergeTime += mergeTime
			totalMergedPrs++
			totalRetestBeforeMerge += retestBeforeMergeCount
		}
	}

	if totalMergedPrs != 0 {
		totalMergeTime = totalMergeTime / float64(len(mergedPullRequestsInTimeRange))
		retestBeforeMergeAvg = totalRetestBeforeMerge / float64(len(mergedPullRequestsInTimeRange))
	}

	info.Summary = prV1Alpha1.Summary{
		CreatedPrsCountInTimeRange: len(createdPullRequestsInTimeRange),
		OpenPrsCount:               len(totalOpenPullRequests),
		MergedPrsCount:             len(totalPullRequestsMerged),
		MergedPrsCountInTimeRange:  len(mergedPullRequestsInTimeRange),
		MergeAvg:                   math.Round(totalMergeTime*100) / 100,
		RetestAvg:                  math.Round(retestAvg*100) / 100,
		RetestBeforeMergeAvg:       math.Round(retestBeforeMergeAvg*100) / 100,
	}
	info.Metrics = d.getMetrics(repo, startDate, endDate)

	return info, nil
}
