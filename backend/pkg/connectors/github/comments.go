package github

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	githubV1Alhpa1 "github.com/konflux-ci/quality-studio/api/apis/github/v1alpha1"
	"github.com/konflux-ci/quality-studio/pkg/constants"
	"github.com/shurcooL/githubv4"
)

// ChatopsMergedPRsBetween returns a slice of PRs that were merged in the time
// frame defined by source as parameter with data
// required by chatops tools. EG of source: redhat-appstudio/e2e-tests
func (gc *Github) ChatopsMergedPRsBetween(source string) (githubV1Alhpa1.ChatopsPRList, error) {
	currentTime := time.Now()
	mergedQueryString := fmt.Sprintf("repo:%s type:pr merged:%s..%s",
		source,
		currentTime.AddDate(0, 0, -14).Format(constants.DateFormat),
		currentTime.AddDate(0, 0, 0).Format(constants.DateFormat),
	)

	mergedQueryResult, err := gc.chatopsPRQuery(mergedQueryString)
	if err != nil {
		return nil, err
	}
	return mergedQueryResult, nil
}

/*
{
  search(query: "repo:redhat-appstudio/e2e-tests type:pr merged:2023-02-18T11:02:49Z..2023-03-04T11:02:49Z", type: ISSUE, first: 100) {
    issueCount
    pageInfo {
      hasNextPage
      endCursor
      startCursor
    }
    edges {
      node {
        ... on PullRequest {
          number
          title
          repository {
            nameWithOwner
          }
          createdAt
          mergedAt
          url
          changedFiles
          additions
          deletions
        }
      }
    }
  }
}
*/

func (c *Github) chatopsPRQuery(query string) (githubV1Alhpa1.ChatopsPRList, error) {
	variables := map[string]interface{}{
		"querystring": githubv4.String(query),
	}

	var mergedQuery struct {
		Search struct {
			IssueCount int
			PageInfo   githubV1Alhpa1.PageInfo
			Nodes      githubV1Alhpa1.ChatopsPRList
		} `graphql:"search(query: $querystring, type: ISSUE, first:100)"`
	}

	err := c.graphqlClient.Query(context.Background(), &mergedQuery, variables)
	return mergedQuery.Search.Nodes, err
}

// RetestsToMerge returns average of retest calls it took to land
// each merged PR in the time.
func (gc *Github) RetestsToMerge(source string) (float64, error) {
	var totalRetests, average float64
	items, err := gc.ChatopsMergedPRsBetween(source)

	if err != nil {
		return 0, err
	}

	for _, prItem := range items {
		totalRetests += RetestCommentsAfterLastPush(&prItem.ChatopsPullRequestFragment.TimelineItems)
	}

	average = totalRetests / float64(len(items))

	if math.IsNaN(average) {
		average = 0
	} else {
		average = math.Round(average*100) / 100
	}

	return average, nil
}

// RetestComments returns the number of /retest or /test comments a PR received
func RetestComments(items *githubV1Alhpa1.TimelineItems) float64 {
	var total float64 = 0
	for _, timelineItem := range items.Nodes {
		if isRetestComment(timelineItem) {
			total++
		}
	}
	return total
}

// RetestCommentsAfterLastPush returns the number of /retest or /test comments a PR received
// after the last commit or force push.
func RetestCommentsAfterLastPush(items *githubV1Alhpa1.TimelineItems) float64 {
	var total float64 = 0
	lastPush := determineLastPush(items)

	for _, timelineItem := range items.Nodes {
		if isRetestCommentAfterLastPush(timelineItem, lastPush) {
			total++
		}
	}
	return total
}

func determineLastPush(pr *githubV1Alhpa1.TimelineItems) time.Time {
	lastPush := zeroDate

	var itemDate time.Time
	for _, timelineItem := range pr.Nodes {
		if isCommit(timelineItem) {
			itemDate = timelineItem.PullRequestCommitFragment.Commit.CommittedDate
		} else if isHeadRefForcePush(timelineItem) {
			itemDate = timelineItem.HeadRefForcePushFragment.CreatedAt
		} else if isBaseRefForcePush(timelineItem) {
			itemDate = timelineItem.BaseRefForcePushFragment.CreatedAt
		}
		if itemDate.After(lastPush) {
			lastPush = itemDate
		}
	}
	return lastPush
}

func isCommit(timelineItem githubV1Alhpa1.TimelineItem) bool {
	return timelineItem.PullRequestCommitFragment != githubV1Alhpa1.PullRequestCommitFragment{}
}

func isHeadRefForcePush(timelineItem githubV1Alhpa1.TimelineItem) bool {
	return timelineItem.HeadRefForcePushFragment.Actor.Login != ""
}

func isBaseRefForcePush(timelineItem githubV1Alhpa1.TimelineItem) bool {
	return timelineItem.BaseRefForcePushFragment.Actor.Login != ""

}

func isRetestComment(timelineItem githubV1Alhpa1.TimelineItem) bool {
	return timelineItem.IssueCommentFragment != githubV1Alhpa1.IssueCommentFragment{} &&
		(strings.HasPrefix(timelineItem.IssueCommentFragment.BodyText, "/retest") ||
			strings.HasPrefix(timelineItem.IssueCommentFragment.BodyText, "/test"))
}

func isRetestCommentAfterLastPush(timelineItem githubV1Alhpa1.TimelineItem, lastPush time.Time) bool {
	return isRetestComment(timelineItem) &&
		timelineItem.IssueCommentFragment.CreatedAt.After(lastPush)
}
