#! /bin/sh

set +x

go build -o ./resourceManager

# to make sure mysql is ready
sleep 15

./resourceManager --mysql.dsn="$RM_DB_CONNECTION_STRING"  --tenancy.db_max_conn=10 --web.listen-address=0.0.0.0:$RM_API_PORT
