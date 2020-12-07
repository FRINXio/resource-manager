# Resource-Manager API tests

This subproject contains various API tests like creating and
executing pools, creating and searching resources and resource types etc.

To run the tests:

1. Start RM, e.g. using `docker-compose -f docker-compose.api-tests.yaml up -d`
1. run `yarn install`
1. create a `.env` file
1. run `yarn test` which only shows failing tests or `yarn jest` for a verbose output

Alternatively, run `test.sh` that does all of it in a single step.
