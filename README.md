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
