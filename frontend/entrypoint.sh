#!/bin/sh
set -ex

URL_TO_REPLACE=$(grep REACT_APP_API_SERVER_URL .env | cut -d '=' -f 2)

sed -i -- "s#$URL_TO_REPLACE#$REACT_APP_API_SERVER_URL#g" "dist/main.bundle.js"
sed -i -- "s#$URL_TO_REPLACE#$REACT_APP_API_SERVER_URL#g" "dist/main.bundle.js.map"

yarn start
