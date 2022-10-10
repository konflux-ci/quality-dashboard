# Backend Server

Server API solution to talk to quality dashboard.

# Specifications
* Structured logging with zap.
* Use go-cache to save quality repositories in cache.
* Use swaggo to create a specific swagger for all server served endpoints
* Stores data in a PostreSQL database
* Uses entgo to model and query the database

# Setup

A proper setup Go workspace using **Go 1.19+ is required**.

Install dependencies:
```
# Go to backend dir and install dependencies
go mod tidy
# Copy the dependencies to vendor folder
go mod vendor
# Create qe-dashboard-backend binary in bin folder. Please add the binary to the path or just execute ./bin/qe-dashboard-backend
make build
```

Environments used by the server:
| Environment Name | Value | Default | Required |
| -- | -- | -- | -- |
| `GITHUB_TOKEN` | Github token to make requests | `` | yes |
| `CODECOV_TOKEN` | CodeCov token to make requests | `` | no |
| `JIRA_TOKEN` | Jira token to read jira issues | yes | --jira-token xxxxx |