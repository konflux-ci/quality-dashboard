package github

import (
	"time"

	"github.com/shurcooL/githubv4"
)

type IssueCommentFragment struct {
	CreatedAt time.Time
	BodyText  string
}

type Commit struct {
	CommittedDate time.Time
}

type PullRequestCommitFragment struct {
	Commit Commit
}

type Actor struct {
	Login string
}

type BaseRefForcePushFragment struct {
	Actor     Actor
	CreatedAt time.Time
}

type HeadRefForcePushFragment struct {
	Actor     Actor
	CreatedAt time.Time
}

type TimelineItem struct {
	IssueCommentFragment      `graphql:"... on IssueComment"`
	PullRequestCommitFragment `graphql:"... on PullRequestCommit"`
	BaseRefForcePushFragment  `graphql:"... on BaseRefForcePushedEvent"`
	HeadRefForcePushFragment  `graphql:"... on HeadRefForcePushedEvent"`
}

type TimelineItems struct {
	Nodes []TimelineItem
}

type ChatopsPullRequestFragment struct {
	Number        int
	CreatedAt     time.Time
	MergedAt      time.Time
	TimelineItems `graphql:"timelineItems(first:100, itemTypes:[ISSUE_COMMENT])"`
}

type MergeQueuePullRequestFragment struct {
	Number        int
	CreatedAt     time.Time
	MergedAt      time.Time
	TimelineItems `graphql:"timelineItems(first:100, itemTypes:[LABELED_EVENT, UNLABELED_EVENT])"`
}

type ChatopsPRList []struct {
	ChatopsPullRequestFragment `graphql:"... on PullRequest"`
}

type PageInfo struct {
	StartCursor githubv4.String
	EndCursor   githubv4.String
	HasNextPage bool
}
