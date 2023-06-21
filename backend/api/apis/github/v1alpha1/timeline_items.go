package v1alpha1

import "time"

type Actor struct {
	Login string
}

type IssueCommentFragment struct {
	CreatedAt time.Time
	BodyText  string
}

type PullRequestCommitFragment struct {
	Commit Commit
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
