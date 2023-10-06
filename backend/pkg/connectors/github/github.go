package github

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v44/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type Github struct {
	client        *github.Client
	graphqlClient *githubv4.Client
}

func NewGithubClient(token string) (*Github, error) {
	ctx, err := context.WithTimeout(context.Background(), 90*time.Second)

	if err != nil {
		return &Github{}, fmt.Errorf("error initializing context with timeout %v", err)
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	return &Github{
		client:        github.NewClient(tc),
		graphqlClient: githubv4.NewClient(tc),
	}, nil
}
