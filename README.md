# Quality Dashboard

The purpose of the quality dashboard is to collect the status of AppStudio services(codecov, build status, test types, what is in progress, blocked and why).
# Specifications

* Collect github information like github actions status.
* Collect unit test coverage from codecov.
* Use patternfly project[(https://www.patternfly.org/v4/get-started/develop/)]
* Golang backend based.

# Install
Prerequisites:
You need to be logged in OpenShift cluster.
For example with oc command: `oc login -u <user> -p <password> <oc_api_url>.`

The install script will deploy quality dashboard and all resources to the OpenShift.

Script creates namespace `appstudio-qe` and deploy OpenShift resources for [backend](https://github.com/redhat-appstudio/quality-dashboard/tree/main/backend/deploy/openshift) and [frontend](https://github.com/redhat-appstudio/quality-dashboard/tree/main/frontend/deploy/openshift).

```
# Run `install.sh` from hack folder to deploy the dashboard
$ /bin/bash hack/install.sh --storage-user <username> --storage-password <password>  --github-token <token>
```

When running the install script, you need to specify these parameters:

| Parameter Name | Description | Required | Example |
| -- | -- | -- | -- |
| `github-token` | Github token to read repositories | yes | --github-token ghp_xxxxx |
| `storage-user` | Database user name | yes | --storage-user admin |
| `storage-password` | Database user password | yes | --storage-password adminPassword |

## Install quality dashboard locally

Steps to install quality dashboard locally for development purposes

1. Start a postgres instance

```bash
    <container-engine> run -p 5432:5432 --name some-postgres -e POSTGRES_PASSWORD=postgres -d postgres # docker-engine is docker or podman
```

2.Open a new terminal and follow backend [instructions](./backend/README.md) to install binaries

3.Start the backend

```bash
    cd backend
    make build
    ./bin/qe-dashboard-backend
```

4.Open a new terminal and install the frontend

```bash
    cd frontend
    yarn # to install the dependencies
    yarn start:dev
```
