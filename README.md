# Quality Dashboard
1. [Purpose](#purpose)
2. [Specifications](#ppecifications)
3. [Install](#installation)
    * [Install locally](#install-quality-dashboard-locally)
4. [Connectors and features](#connectors-and-features)
    * [Teams](#teams)
    * [Openshift CI and Prow Jobs](#openshift-ci-and-prow-jobs)
    * [Github](#github)
    * [Codecov](#codecov)
    * [Jira](#jira)

# Purpose
The purpose of the quality dashboard is to collect the status of AppStudio services:

* Github repositories data and action status
* Code coverage from codecov
* Build status and test types
* Openshift CI jobs and their statistics
* Jira issues impacting Appstudio (in progress, blockers)

# Specifications

Quality Dashboard implements a Golang-based backed and stores data in a PostgreSQL database, modelled and queried with [entgo](https://entgo.io/) framework.
Different specific connectors are developed to pull data from different sources:
* Github connector: to pull data from github, such as repositories information and actions status
* Codecov connector: to pull code coverage data from Codecov
* ProwJobs connector: to pull automatically data about prow jobs executions impacting the repositories 
* Jira connector: to pull issues from Jira

The database will retain last 10 days of CI job executions. 

The frontend uses [patternfly project](https://www.patternfly.org/v4/get-started/develop/).

# Install

You need to be logged in to an OpenShift cluster first.
Example (oc command): `oc login -u <user> -p <password> <oc_api_url>.`

The install script will deploy quality dashboard and all resources to your OpenShift cluster.
The script will create a namespace `appstudio-qe` and deploy all OpenShift resources needed by [backend](https://github.com/redhat-appstudio/quality-dashboard/tree/main/backend/deploy/openshift) and [frontend](https://github.com/redhat-appstudio/quality-dashboard/tree/main/frontend/deploy/openshift).

To run the script:

```
# Run install.sh from hack folder to deploy the dashboard
/bin/bash hack/install.sh --storage-user <username> --storage-password <password> --github-token <token> --jira-token <token>
```

When running the install script, you need to specify these parameters:

| Parameter Name | Description | Required | Example |
| -- | -- | -- | -- |
| `github-token` | Github token to read repositories | yes | --github-token ghp_xxxxx |
| `jira-token` | Jira token to read jira issues | yes | --jira-token xxxxx |
| `storage-user` | Database user name | yes | --storage-user admin |
| `storage-password` | Database user password | yes | --storage-password adminPassword |

## Install quality dashboard locally

To install quality dashboard locally (for development purposes):

 - Start a postgres instance with your container engine (docker or podman for example
```bash
    <container-engine> run -p 5432:5432 --name some-postgres -e POSTGRES_PASSWORD=postgres -d postgres
```

 - Open a new terminal and follow backend [instructions](./backend/README.md) to install binaries.

 - Start the backend:

```bash
    cd backend # Compile following step 2
    ./bin/server-runtime
```

- Open a new terminal and install the frontend:

```bash
    cd frontend
    yarn # to install the dependencies
    yarn start:dev
```
## Connectors and features

### Teams
All data is organized by Teams: a team groups a set of repositories to show data in a more concise manner and acts as a global filter. 
All the teams that have been created will be listed in a table in the Teams page, where they can also be managed.
Switching a team from the main toolbar, will update the context for the whole view in the dashboard.

### Openshift CI and Prow Jobs
The Openshift CI connector will collect and show an overview of the last 10 days of jobs execution, by repository and job type. 
Current job types are: presubmit, periodic and postsubmit.
If more than one job per repository and job type is there, the connector will collect all of them.
The dashboard will present the last 10 days of data in a chart, for day to day inspection, and the averages of the whole period of time.
Also, just for periodic jobs, we show the test suites output of the last executed job. 

### Github
The Github connector will pull data from Github, such has repositories info and action status.

### Codecov
The codecov connector will pull code coverage data from Codecov.

### Jira
The Jira connector will pull data from Jira. We gather Jira issues that are impacting Appstudio (such as blockers, in progress, etc.) and present them in the dashboard for quick reference. 