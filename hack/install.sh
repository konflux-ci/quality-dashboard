#!/bin/bash

export STORAGE_PASSWORD=""
export STORAGE_USER=""
export GITHUB_TOKEN=""
export JIRA_TOKEN=""
export WORKSPACE=$(dirname $(dirname $(readlink -f "$0")))
export SECRET_DASHBOARD_TMP=$(mktemp)
export FRONTEND_DEPLOYMENT_TMP=$(mktemp)

while [[ $# -gt 0 ]]
do
    case "$1" in
        -p|--storage-password)
            STORAGE_PASSWORD=$(echo -n $2 | base64)
            ;;
        -u|--storage-user)
            STORAGE_USER=$(echo -n $2 | base64)
            ;;
        -g|--github-token)
            GITHUB_TOKEN=$(echo -n $2 | base64)
            ;;
        -jt|--jira-token)
            JIRA_TOKEN=$(echo -n $2 | base64)
            ;;
        *)
            ;;
    esac
    shift  # Shift each argument out after processing them
done

if [[ "${STORAGE_PASSWORD}" == "" ]]; then
  echo "[ERROR] Storage password flag is missing. Use '--storage-password <value>' or '-p <value>' to create a storage password for the quality dashboard"
  exit 1
fi

if [[ "${STORAGE_USER}" == "" ]]; then
  echo "[ERROR] Storage database flag is missing. Use '--storage-user <value>' or '-u <value>' to create a storage database for the quality dashboard"
  exit 1
fi

if [[ "${GITHUB_TOKEN}" == "" ]]; then
  echo "[ERROR] Github Token flag is missing. Use '--github-token <value>' or '-g <value>' to allow quality dashboard to make request to github"
  exit 1
fi

echo "[INFO] Starting Quality dashboard..."
echo "   Storage Password   : "${STORAGE_PASSWORD}""
echo "   Storage Database   : "${STORAGE_USER}""
echo "   Github Token       : "${GITHUB_TOKEN}""
echo ""

# Postgres service name from file quality-dashboard/backend/deploy/openshift/service.yaml
export POSTGRES_SERVICE="cG9zdGdyZXMtc2VydmljZQ=="

# Replace variables in /backend/deploy/openshift/secret.yaml
cat "${WORKSPACE}/backend/deploy/openshift/secret.yaml" |
    sed -e "s#REPLACE_STORAGE_PASSWORD#${STORAGE_PASSWORD}#g" |
    sed -e "s#REPLACE_STORAGE_USER#${STORAGE_USER}#g" |
    sed -e "s#REPLACE_GITHUB_TOKEN#${GITHUB_TOKEN}#g" |
    sed -e "s#REPLACE_JIRA_TOKEN#${JIRA_TOKEN}#g" |
    sed -e "s#REPLACE_WITH_RDS_ENDPOINT#${POSTGRES_SERVICE}#g" |
    cat > ${SECRET_DASHBOARD_TMP}

# Namespace
oc create namespace appstudio-qe || true

# BACKEND
echo -e "[INFO] Deploying Quality dashboard backend"

oc apply -f ${SECRET_DASHBOARD_TMP}
oc apply -f "${WORKSPACE}/backend/deploy/openshift/postgres.yaml"
oc apply -f "${WORKSPACE}/backend/deploy/openshift/deployment.yaml"
oc apply -f "${WORKSPACE}/backend/deploy/openshift/service.yaml"
oc apply -f "${WORKSPACE}/backend/deploy/openshift/route.yaml"

export BACKEND_ROUTE=$(oc get route quality-backend-route -n appstudio-qe -o json | jq -r '.spec.host')

# FRONTEND
echo -e "[INFO] Deploying Quality dashboard frontend"

cat "${WORKSPACE}/frontend/deploy/openshift/deployment.yaml" |
    sed -e "s#REPLACE_BACKEND_URL#${BACKEND_ROUTE}#g" |
    cat > ${FRONTEND_DEPLOYMENT_TMP}

oc apply -f ${FRONTEND_DEPLOYMENT_TMP}
oc apply -f "${WORKSPACE}/frontend/deploy/openshift/service.yaml"
oc apply -f "${WORKSPACE}/frontend/deploy/openshift/route.yaml"

echo ""
echo "Frontend is accessible from: http://"$(oc get route/quality-frontend-route -n appstudio-qe -o go-template='{{.spec.host}}{{"\n"}}')""
