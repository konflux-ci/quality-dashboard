package github

import (
	"context"
	"fmt"
	"strings"
	"time"

	grafanaGH "github.com/grafana/github-datasource/pkg/github"
	githubV1Alhpa1 "github.com/konflux-ci/quality-studio/api/apis/github/v1alpha1"
	"github.com/pkg/errors"
	"github.com/shurcooL/githubv4"
)

// QueryListPullRequests lists all pull requests in a repository
//
//	{
//	  search(query: "is:pr repo:redhat-appstudio/e2e-tests merged:2020-08-19..*", type: ISSUE, first: 100) {
//	    nodes {
//	      ... on PullRequest {
//	        id
//	        title
//	      }
//	  }
//	}
type QueryListPullRequests struct {
	Search struct {
		Nodes []struct {
			PullRequest githubV1Alhpa1.PullRequest `graphql:"... on PullRequest"`
		}
		PageInfo githubV1Alhpa1.PageInfo
		// Do not set the number of PRs per page too high to avoid query timeout
	} `graphql:"search(query: $query, type: ISSUE, first: 50, after: $cursor)"`
}

// GetAllPullRequests uses the graphql search endpoint API to search all pull requests in the repository
func (gh *Github) GetAllPullRequests(ctx context.Context, opts githubV1Alhpa1.ListPullRequestsOptions) (githubV1Alhpa1.PullRequests, error) {
	var (
		variables = map[string]interface{}{
			"cursor": (*githubv4.String)(nil),
			"query":  githubv4.String(buildQuery(opts)),
		}

		pullRequests = []githubV1Alhpa1.PullRequest{}
	)

	for {
		q := &QueryListPullRequests{}
		if err := gh.graphqlClient.Query(ctx, q, variables); err != nil {
			return nil, errors.WithStack(err)
		}
		prs := make([]githubV1Alhpa1.PullRequest, len(q.Search.Nodes))
		for i, v := range q.Search.Nodes {
			prs[i] = v.PullRequest
		}

		pullRequests = append(pullRequests, prs...)

		if !q.Search.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = q.Search.PageInfo.EndCursor
	}

	return pullRequests, nil
}

// GetPullRequestsInRange uses the graphql search endpoint API to find pull requests in the given time range.
func (gh *Github) GetPullRequestsInRange(ctx context.Context, opts githubV1Alhpa1.ListPullRequestsOptions, from time.Time, to time.Time) (githubV1Alhpa1.PullRequests, error) {
	var q string

	if opts.TimeField != githubV1Alhpa1.PullRequestNone {
		q = fmt.Sprintf("%s:%s..%s", opts.TimeField.String(), from.Format(time.RFC3339), to.Format(time.RFC3339))
	}

	if opts.Query != nil {
		q = fmt.Sprintf("%s %s", *opts.Query, q)
	}

	return gh.GetAllPullRequests(ctx, githubV1Alhpa1.ListPullRequestsOptions{
		Repository: opts.Repository,
		Owner:      opts.Owner,
		TimeField:  opts.TimeField,
		Query:      &q,
	})
}

// buildQuery builds the "query" field for Pull Request searches
func buildQuery(opts githubV1Alhpa1.ListPullRequestsOptions) string {
	search := []string{
		"is:pr",
	}

	if opts.Repository == "" {
		search = append(search, fmt.Sprintf("org:%s", opts.Owner))
	} else {
		search = append(search, fmt.Sprintf("repo:%s/%s", opts.Owner, opts.Repository))
	}

	if opts.Query != nil {
		queryString, err := grafanaGH.InterPolateMacros(*opts.Query)
		if err != nil {
			return strings.Join(search, " ")
		}
		search = append(search, queryString)
	}

	return strings.Join(search, " ")
}
