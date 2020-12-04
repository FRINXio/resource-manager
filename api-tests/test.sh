set -x
dirname="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd ${dirname}


docker-compose -f docker-compose.api-tests.yaml up -d

trap 'docker-compose -f docker-compose.api-tests.yaml logs resource-manager' err exit

sleep 5
cp .env-LOCAL-DEV-CONFIG .env
yarn install
yarn jest
