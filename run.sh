#! /bin/sh

set -x

go build -o ./resourceManager

# to make sure mysql is ready
sleep 15

RM_DB_CONNECTION_STRING_DEFAULT="root:root@tcp(localhost:3306)/?charset=utf8&parseTime=true&interpolateParams=true"
RM_API_PORT_DEFAULT=8884

./resourceManager --mysql.dsn="${RM_DB_CONNECTION_STRING:-$RM_DB_CONNECTION_STRING_DEFAULT}"  --tenancy.db_max_conn=10 \
--web.listen-address=0.0.0.0:${RM_API_PORT:-$RM_API_PORT_DEFAULT}
