# Quality Dashboard
- [Quality Dashboard](#quality-dashboard)
- [Purpose](#purpose)
- [Specifications](#specifications)
  - [Backend](#backend)
        - [About entgo framework](#about-entgo-framework)
        - [APIs](#apis)
  - [Frontend](#frontend)
- [Install quality dashboard locally](#install-quality-dashboard-locally)
  - [Prerequisites](#prerequisites)
    - [Dex for oauth](#dex-for-oauth)
    - [Backend](#backend-1)
        - [Frontend](#frontend-1)
  - [Features](#features)
    - [User Configuration](#user-configuration)
    - [Teams](#teams)
    - [Config](#config)
    - [Bug SLIs](#bug-slis)
  - [Connectors](#connectors)
    - [Openshift CI and Prow Jobs](#openshift-ci-and-prow-jobs)
    - [Github](#github)
    - [Codecov](#codecov)
    - [Jira](#jira)

# Purpose
The purpose of the quality dashboard is to collect the status of AppStudio services:

* Github repositories data and action status
* Code coverage from codecov
* Build status and test types
* Openshift CI jobs and their statistics
* Jira issues impacting RHTAP (in progress, blockers)

# Specifications

## Backend
Quality Dashboard implements a Golang-based backed and stores data in a PostgreSQL database, modeled and queried with [entgo](https://entgo.io/) framework.

Different specific connectors are developed to pull data from different sources:
* Github connector: to pull data from GitHub, such as repository information and action status
* Codecov connector: to pull code coverage data from Codecov
* ProwJobs connector: to pull automatically data about prow jobs executions impacting the repositories
* Jira connector: to pull issues from Jira

The database will retain the last 10 days of CI job executions.

##### About entgo framework
Ent is an Object Relational Mapping (ORM) framework for modeling any database schema as Go objects. The only thing you need to do is to define a schema and Ent will handle the rest. Your schema will be validated first, then Ent will generate a well-typed and idiomatic API.
The generated API is used to manage the data and will contain:
* Client objects used to interact with the database
* CRUD builders for each schema type
* Entity object (Go struct) for each the schema type

You can use such generated code to build your endpoints and manipulate the database easily and programmatically.

The schema for Quality Dashboard data types is located [here](https://github.com/redhat-appstudio/quality-dashboard/tree/main/backend/pkg/storage/ent/schema). You can refer to entgo [documentation](https://entgo.io/docs/schema-def) for syntax details.

After adding new data types to the schema (or editing the existing ones), you have to execute the following command in `backend/pkg/storage/ent` to re-build the model:

```
go run -mod=mod entgo.io/ent/cmd/ent generate ./schema --target ./db --feature sql/upsert
```

The generated code will be saved into the `backend/pkg/storage/ent/db` folder.

The `backend/pkg/storage/ent/client` package implements the database client used to interact with the database.

In turn, the database client package implements the storage interface used by the server.

##### APIs
The backend server exposes a set of APIs to interact with data. The implementation of the API server is located at `backend/api` and uses a basic HTTP router configuration.

## Frontend
The frontend component is a React web application that uses [patternfly project](https://www.patternfly.org/v4/get-started/develop/) to build the UI.
It interacts with the backend via HTTP api endpoints.

# Install quality dashboard locally

To install quality dashboard locally (for development purposes) you will need to run both backend and frontend on your own.

## Prerequisites
* Make sure you have Go (Golang) installed on your system, as DEX and backend is written in Go
* You will need a GitHub account and access to create OAuth applications on GitHub

### Dex for oauth
To install dex locally you need to follow the next steps:

* Clone the DEX GitHub repository to your local machine.
```bash
    git clone https://github.com/dexidp/dex.git
```

* Change your working directory to the DEX repository.
``` bash
    cd dex
```
* Configure GitHub OAuth App:
  * Navigate to Settings -> Developer settings -> OAuth Apps.
  * Click on "New OAuth App" and fill in the required information. You will need to specify a "Homepage URL" and a "Callback URL." For local development, you can use http://localhost:5555/callback as the callback URL.
  * After creating the OAuth App, you will receive a Client ID and Client Secret. Keep these values handy.
* Create a configuration file for DEX. You can use the provided examples/config-dev.yaml file as a starting point and modify it according to your needs. Make sure to configure the GitHub connector with your GitHub OAuth App's Client ID and Client Secret.
* Build DEX using the following command:
``` bash
    go build ./cmd/dex
```

Then, run DEX with your configuration file:
``` bash
    ./dex serve <path-to-your-config-file>
    Replace <path-to-your-config-file> with the actual path to your DEX configuration file.
```
Example configuration about GitHub provider can be found [here](https://dexidp.io/docs/connectors/github/#configuration)

Please note that these are general steps, and the exact steps may vary based on your specific requirements and DEX configuration. Make sure to refer to the DEX documentation and GitHub OAuth documentation for more detailed information and troubleshooting if needed.

### Backend

First, you need to have a PostgreSQL instance running to host local data. You can start one with your favorite container engine (docker or podman)

```bash
    podman run -p 5432:5432 --name some-postgres -e POSTGRES_PASSWORD=postgres -d postgres:14
```

After that, you need to build the backend binaries. To do that you can follow the backend [instructions](./backend/README.md).

Once built, run the backend server in a terminal:
```bash
    # from the backend folder
    ./bin/server-runtime
```
If you specified some different values for the database container, you can override the default values by exporting the following environment variables:
* `POSTGRES_ENT_HOST`
* `POSTGRES_ENT_PORT`
* `POSTGRES_ENT_DATABASE`
* `POSTGRES_ENT_USER`
* `POSTGRES_ENT_PASSWORD`
* `GITHUB_TOKEN`

The server runtime will take care of initializing the database structure and pull the data.


##### Frontend

Open a new terminal, navigate to the frontend folder, install dependencies and run:

```bash
    cd frontend
    yarn
    yarn start:dev
```
or with npm:
```bash
    cd frontend
    npm install
    npm run start:dev
```

## Features

### User Configuration
You can set your user's settings by clicking in (your user name) -> Settings.

Available options:
  - set default team

### Teams
All data is organized by Teams: a team groups a set of repositories to show data more concisely and acts as a global filter.
All the teams that have been created will be listed in a table on the Teams page, where they can also be managed.
Switching a team from the main toolbar will update the context for the whole view in the dashboard.

### Config
The Config page provides a quick way of adding multiple teams and repositories in the DB. It presents a code editor where you can set teams and their corresponding JIRA projects and repositories, by downloading an existing config or starting from scratch.

Please, note that:
 - different teams can not have the same repository
 - different teams can not have the same description

The config should conform to the following template:
```
teams:
   - name: team-example
     description: description-example
     jira_projects:
        - STONE
     repositories:
        - name: e2e-tests
          organization: redhat-appstudio
        - name: quality-dashboard
          organization: redhat-appstudio
```


### Bug SLIs

#### Description
Currently, we have defined four Service Level Objectives (SLOs):

- **Resolution Time Bug SLO**: Aims to ensure that Blocker, Critical, and Major bugs are resolved in a reasonable period.
- **Response Time Bug SLO**: Aims to ensure that Blocker and Critical bugs are assigned in the early phase of the bug's life.
- **Priority Triage Time Bug SLO**: Aims to ensure that untriaged bugs are prioritized in the early phase of the bug's life.
- **Component Assignment Triage Time Bug SLO**: Aims to ensure that bugs are assigned to a component in the early phase of the bug's life.

And three Service Level Indicators (SLIs):
- **green**: means that the bug meets the defined SLO.
- **yellow**: means that the bug is almost not meeting the defined SLO.
- **red**: means that the bug is not meeting the defined SLO.

| **SLO**                          | **Target Value**                                                                                                                                    | **SLIs**                                                                                                                                                                                                                                                   |
|----------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Bug Resolution Time              | Resolve blocker bug in < 10 days<br><br><br><br><br><br>Resolve critical bug in < 20 days<br><br><br><br><br><br><br>Resolve major bug in < 40 days | Green:  age < 5 days<br>Yellow: age  > 5 days<br>Red:    age > 10 days<br><br><br><br>Green:  age < 10 days<br>Yellow: age  > 10 days<br>Red:    age > 20 days<br><br><br><br><br>Green:  age < 20 days<br>Yellow: age  > 20 days<br>Red:    age > 40 days |
| Bug Response Time                | Blocker or Critical bug will get assigned in < 2 days                                                                                               | Red:    unassigned > 2 days                                                                                                                                                                                                                                |
| Priority Triage Time             | Bug will get assigned priority in < 2 day                                                                                                           | Yellow: age > 1 days & untriaged<br>Red:    age > 2 days & untriaged                                                                                                                                                                                       |
| Component Assignment Triage Time | Bug will get assigned component in < 1 working day                                                                                                  | Red:     age > 1 working days & component empty             

#### Bug SLIs page
With this new page, you can observe which RHTAPBUGS are not meeting the defined Bug SLOs, which can be helpful to better (re)prioritize them.
You can also find metrics regarding how many bugs are not meeting or almost not meeting the defined SLOs by each component.

#### Bug SLIs alerts
The SLIs alerts are sent daily at 9 AM by the RHTAP QE bot in the [rhtap-bug-slis-alerts](https://app.slack.com/client/E030G10V24F/C062AF1RFK8) channel.
To send the alerts daily, we are using a chron job. You can find all the code on [rhtap_bug_slis_alert.go](https://github.com/redhat-appstudio/quality-dashboard/blob/main/backend/api/server/rhtap_bug_slis_alert.go).
Note that bugs with status as "Waiting" and "Release Pending" are excluded.

#### How Bug SLIs are being got?
As the first step, we need to know how many days the issues have not been resolved, prioritized, or assigned. So, when adding or updating a Jira Issue Bug, we need to keep this information on the bugs db table through ['DaysWithoutAssignee', 'DaysWithoutPriority', and 'DaysWithoutResolution' fields](https://github.com/redhat-appstudio/quality-dashboard/blob/main/backend/pkg/storage/ent/client/jira_bugs.go).

Then, we need to use the information above to get the issues that are not meeting or almost not meeting the defined SLIs. For that, we defined a [function to get each SLI](https://github.com/redhat-appstudio/quality-dashboard/blob/main/backend/api/server/router/jira/bug_slos.go).


### RHTAPBUGS Impact on CI

#### Description
RHTAPBUGS Impact on CI plugin lists the RHTAPBUGS that are impacting CI, by showing the Jira Key, Jira Status, Error Message, and Frequency. You can add, update, or delete them.
To add a new entry, you need to point out the Jira Key of the bug and the associated error message.

#### How the frequency/impact is being measured?
Since we are saving the prow jobs on our db, we can also [save the error messages](https://github.com/redhat-appstudio/quality-dashboard/blob/main/backend/api/server/prow_rotate.go) related to each prow job.
We are saving two types of error messages:
 - E2E Failed Test Messages: Messages related to failed E2E tests.
 - Build Error Logs: If there is no E2E Failed Test Messages, we save the last 50 lines of the build-log.txt in order to help us find infra errors, for example. We are only saving the last 50 days because we believe it is enough to catch the build errors.

To calculate the impact of each failure (in the Failures table), in the date time range selected, we will search for all the prow jobs and verify in how many prow jobs the bug's error message is present in the E2E Failed Test Messages or Build Error Logs. You can find the code (here)[https://github.com/redhat-appstudio/quality-dashboard/blob/main/backend/pkg/storage/ent/client/failure.go].

## Connectors

### Openshift CI and Prow Jobs
The Openshift CI connector will collect and show an overview of the last 10 days of jobs execution, by repository and job type.
Current job types are: presubmit, periodic and postsubmit.
If more than one job per repository and job type is there, the connector will collect all of them.
The dashboard will present the last 10 days of data in a chart, for day to day inspection, and the averages of the whole period of time.
Also, just for periodic jobs, we show the test suites output of the last executed job.

### Github
The Github connector will pull data from Github, such as repositories info and action status.

### Codecov
The codecov connector will pull code coverage data from Codecov.

### Jira
The Jira connector will pull data from Jira. We gather Jira issues that are impacting RHTAP (such as blockers, in progress, etc.) and present them in the dashboard for quick reference.
