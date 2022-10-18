# Quality Dashboard
1. [Purpose](#purpose)
2. [Specifications](#ppecifications)
    * [Backend](#backend)
    * [Frontend](#frontend)
3. [Install](#installation)
    * [Install locally](#install-quality-dashboard-locally)
4. [Features](#features)
    * [Teams](#teams)
5. [Connectors](#connectors)
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

## Backend
Quality Dashboard implements a Golang-based backed and stores data in a PostgreSQL database, modelled and queried with [entgo](https://entgo.io/) framework.

Different specific connectors are developed to pull data from different sources:
* Github connector: to pull data from github, such as repositories information and actions status
* Codecov connector: to pull code coverage data from Codecov
* ProwJobs connector: to pull automatically data about prow jobs executions impacting the repositories 
* Jira connector: to pull issues from Jira

The database will retain last 10 days of CI job executions. 

##### About entgo framewrok
Ent is an Object Relational Mapping (ORM) framework for modeling any database schema as Go objects. The only thing you need to do is to define a schema and Ent will handle the rest. Your schema will be validated first, then Ent will generate a well-typed and idiomatic API.
The generated API is used to manage the data and will contain:
* Client objects used to interact with the database
* CRUD builders for each schema type
* Entity object (Go struct) for each the schema type

You can use such generated code to build your endopints and manipulate the database in an easy and programmatic way. 

The schema for Quality Dashboard data types is located [here](https://github.com/redhat-appstudio/quality-dashboard/tree/main/backend/pkg/storage/ent/schema). You can refer to entgo [documentation](https://entgo.io/docs/schema-def) for syntax details. 

After adding new data types to the schema (or editing the existing ones), you have to execute the following command in `backend/pkg/storage/ent` to re-build the model:

```
go run -mod=mod entgo.io/ent/cmd/ent generate ./schema --target ./db
```

The generated code will be saved into the `backend/pkg/storage/ent/db` folder.

The `backend/pkg/storage/ent/client` package implements the database client used to interact with the database. 

In turn, the database client package implements the storage interface used by the server.


##### APIs
The backend server exposes a set of APIs to interact with data. The implementation of the API server is located at `backend/api` and uses a basic HTTP router configuration. 

## Frontend 
The frontend component is a React web application that uses [patternfly project](https://www.patternfly.org/v4/get-started/develop/) to build the UI.
It interacts with the backend via HTTP api endpoints. 

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

To install quality dashboard locally (for development purposes) you will need to run both backend and frontend by your own. 

##### Backend

First, you need to have a PostgreSQL instance running to host local data. You can start one with your favourite container engine (docker or podman)

```bash
    podman run -p 5432:5432 --name some-postgres -e POSTGRES_PASSWORD=postgres -d postgres
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

### Teams
All data is organized by Teams: a team groups a set of repositories to show data in a more concise manner and acts as a global filter. 
All the teams that have been created will be listed in a table in the Teams page, where they can also be managed.
Switching a team from the main toolbar, will update the context for the whole view in the dashboard.

## Connectors

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
