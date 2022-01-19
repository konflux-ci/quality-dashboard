# Quality Dashboard

The purpose of the quality dashboard is to collect the status of AppStudio services(codecov, build status, test types, what is in progress, blocked and why).
# Specifications

* Collect github information like github actions status.
* Collect unit test coverage from codecov.
* Use patternfly project[(https://www.patternfly.org/v4/get-started/develop/)]
* Golang backend based.

# Install

```
# Run `install.sh` from hack folder to deploy the dashboard
$ /bin/bash hack/install.sh --storage-password <value> --storage-user admin --github-token <value>
```

In progress.