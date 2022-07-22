package github

import (
	"context"
	"time"

	"github.com/google/go-github/v44/github"
	"golang.org/x/oauth2"
)

type Github struct {
	client *github.Client
}

func NewGithubClient(token string) *Github {
	ctx, _ := context.WithTimeout(context.Background(), 90*time.Second)

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	return &Github{
		client: github.NewClient(tc),
	}
}
