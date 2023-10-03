#!/bin/sh
set -ex

ENVFILE=$(pwd)/.env
if [ -f "$ENVFILE" ]; then
    rm -rf $ENVFILE
fi

# creates env file to start frontend
cat <<EOT >> .env
REACT_APP_API_SERVER_URL=$REACT_APP_API_SERVER_URL
DEX_ISSUER=$DEX_ISSUER
FRONTEND_REDIRECT_URI=$FRONTEND_REDIRECT_URI
DEX_APPLICATION_ID=$DEX_APPLICATION_ID
EOT

yarn install --network-timeout 1000000
yarn build && yarn start
