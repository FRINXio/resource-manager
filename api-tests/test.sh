#!/bin/bash
set -xe
dirname="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd ${dirname}

docker compose -f docker-compose.api-tests.yaml up -d

sleep 5 # wait for postgres
cp .env-LOCAL-DEV-CONFIG .env
yarn install
"$@"
