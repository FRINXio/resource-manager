#! /bin/sh

set -x

if [ $DEBUG = "true" ]
then
  echo "Running in DEBUG mode"
  go build -gcflags \"all=-N\" -o ./resourceManager
  CMD="dlv --continue --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec ./resourceManager --"
else
  echo "Running in PROD mode, use DEBUG=true to start in debug mode"
  go build -o ./resourceManager
  CMD="./resourceManager"
fi

RM_DB_CONNECTION_STRING_DEFAULT="root:root@tcp(localhost:3306)/?charset=utf8&parseTime=true&interpolateParams=true"
RM_API_PORT_DEFAULT=8884

# Default no admin roles, all users have access
RM_ADMIN_ROLES_DEFAULT=""
# Default no admin groups, all users have access
RM_ADMIN_GROUPS_DEFAULT=""

$CMD \
--mysql.dsn="${RM_DB_CONNECTION_STRING:-$RM_DB_CONNECTION_STRING_DEFAULT}"  --tenancy.db_max_conn=10 \
--web.listen-address=0.0.0.0:${RM_API_PORT:-$RM_API_PORT_DEFAULT} \
--rbac.admin-roles=${RM_ADMIN_ROLES:-$RM_ADMIN_ROLES_DEFAULT} --rbac.admin-groups=${RM_ADMIN_GROUPS:-$RM_ADMIN_GROUPS_DEFAULT}
