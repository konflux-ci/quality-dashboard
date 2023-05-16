# Backend Server

Server API solution to talk to quality dashboard.

# Specifications
* Structured logging with zap.
* Use go-cache to save quality repositories in cache.
* Use swaggo to create a specific swagger for all server served endpoints
* Stores data in a PostreSQL database
* Uses entgo to model and query the database

# Build

A proper setup Go workspace using **Go 1.19+ is required**.

Install dependencies:
```
# Go to backend dir and install dependencies
go mod tidy
# Copy the dependencies to vendor folder
go mod vendor
# Build the backend binary
make build
```

# Execution

Setup environment variables:

| Environment Name | Value | Default | Required |
| -- | -- | -- | -- |
| `GITHUB_TOKEN` | Github token to make requests | `` | yes |
| `CODECOV_TOKEN` | CodeCov token to make requests | `` | no |
| `JIRA_TOKEN` | Jira token to read jira issues | yes | --jira-token xxxxx |


You can run the backend server binary at
```
./bin/server-runtime
```