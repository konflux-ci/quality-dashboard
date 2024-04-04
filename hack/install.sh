#!/bin/bash

export STORAGE_PASSWORD=""
export STORAGE_USER=""
export GITHUB_TOKEN=""
export JIRA_TOKEN=""
export WORKSPACE=$(dirname $(dirname $(readlink -f "$0")))
export SECRET_DASHBOARD_TMP=${WORKSPACE}/backend/deploy/overlays/local/secrets.txt
export FRONTEND_DEPLOYMENT_TMP=$(mktemp)

while [[ $# -gt 0 ]]
do
    case "$1" in
        -p|--storage-password)
            STORAGE_PASSWORD=$2
            ;;
        -u|--storage-user)
            STORAGE_USER=$2
            ;;
        -g|--github-token)
            GITHUB_TOKEN=$2
            ;;
        -jt|--jira-token)
            JIRA_TOKEN=$2
            ;;
        -di|--dex-issuer)
            DEX_ISSUER=$2
            ;;
        -da|--dex-application)
            DEX_APPLICATION_ID=$2
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

cat << EOF > ${SECRET_DASHBOARD_TMP}
storage-database=quality
storage-user=${STORAGE_USER}
storage-password=${STORAGE_PASSWORD}
github-token=${GITHUB_TOKEN}
rds-endpoint=postgres-service
jira-token=${JIRA_TOKEN}
dex-issuer=${DEX_ISSUER}
dex-application-id=${DEX_APPLICATION_ID}
EOF

# Namespace
NS=quality-dashboard
oc create namespace $NS || true

# BACKEND
echo -e "[INFO] Deploying Quality dashboard backend"

oc apply -k ${WORKSPACE}/backend/deploy/overlays/local

# oc apply -k ${WORKSPACE}/frontend/deploy/openshift


# export BACKEND_ROUTE=$(oc get route backend -n $NS -o json | jq -r '.spec.host')
export BACKEND_ROUTE=$(oc get route backend -n $NS -o json | jq -r '.spec.host') > ${WORKSPACE}/frontend/deploy/overlays/local/configmap.txt
cat <<EOF >${WORKSPACE}/frontend/deploy/overlays/local/configmap.txt
first line
second line
third line
EOF
# FRONTEND
echo -e "[INFO] Deploying Quality dashboard frontend"
oc apply -f ${WORKSPACE}/frontend/deploy/base/route.yaml

export FRONTEND_REDIRECT_URI=$(oc get route frontend -n $NS -o json | jq -r '.spec.host')/login

cat <<EOF > ${WORKSPACE}/frontend/deploy/overlays/local/configmap.txt
FRONTEND_REDIRECT_URI=https://${FRONTEND_REDIRECT_URI}
BACKEND_ROUTE=${BACKEND_ROUTE}
EOF

oc apply -k ${WORKSPACE}/frontend/deploy/overlays/local || true

echo ""
echo "Frontend is accessible from: http://"$(oc get route/frontend -n $NS -o go-template='{{.spec.host}}{{"\n"}}')""
