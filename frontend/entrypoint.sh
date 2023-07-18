#!/bin/sh
set -ex

REACT_APP_API_SERVER_URL_REPLACE=$(grep REACT_APP_API_SERVER_URL .env | cut -d '=' -f 2)
DEX_ISSUER_REPLACE=$(grep DEX_ISSUER .env | cut -d '=' -f 2)
FRONTEND_REDIRECT_URI_REPLACE=$(grep FRONTEND_REDIRECT_URI .env | cut -d '=' -f 2)
DEX_APPLICATION_ID_REPLACE=$(grep FRONTEND_REDIRECT_URI .env | cut -d '=' -f 2)

sed -i -- "s#$REACT_APP_API_SERVER_URL_REPLACE#$REACT_APP_API_SERVER_URL#g" "dist/main.bundle.js"
sed -i -- "s#$REACT_APP_API_SERVER_URL_REPLACE#$REACT_APP_API_SERVER_URL#g" "dist/main.bundle.js.map"

sed -i -- "s#$DEX_ISSUER_REPLACE#$DEX_ISSUER#g" "dist/main.bundle.js"
sed -i -- "s#$DEX_ISSUER_REPLACE#$DEX_ISSUER#g" "dist/main.bundle.js.map"

sed -i -- "s#$FRONTEND_REDIRECT_URI_REPLACE#$FRONTEND_REDIRECT_URI#g" "dist/main.bundle.js"
sed -i -- "s#$FRONTEND_REDIRECT_URI_REPLACE#$FRONTEND_REDIRECT_URI#g" "dist/main.bundle.js.map"

sed -i -- "s#$DEX_APPLICATION_ID_REPLACE#$DEX_APPLICATION_ID#g" "dist/main.bundle.js"
sed -i -- "s#$DEX_APPLICATION_ID_REPLACE#$DEX_APPLICATION_ID#g" "dist/main.bundle.js.map"

yarn start