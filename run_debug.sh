#! /bin/sh

set -x

go build -gcflags \"all=-N\" -o ./resourceManager

RM_DB_CONNECTION_STRING_DEFAULT="root:root@tcp(localhost:3306)/?charset=utf8&parseTime=true&interpolateParams=true"
RM_API_PORT_DEFAULT=8884

# With debugging
dlv --continue --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec ./resourceManager \
 -- --mysql.dsn="${RM_DB_CONNECTION_STRING:-$RM_DB_CONNECTION_STRING_DEFAULT}"  --tenancy.db_max_conn=10 \
 --web.listen-address=0.0.0.0:${RM_API_PORT:-$RM_API_PORT_DEFAULT}
