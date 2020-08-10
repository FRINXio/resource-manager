#! /bin/sh

set +x

go build -gcflags \"all=-N\" -o ./resourceManager

sleep 15

# With debugging
go get github.com/go-delve/delve/cmd/dlv
dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec ./resourceManager \
 -- --mysql.dsn="$RM_DB_CONNECTION_STRING"  --tenancy.db_max_conn=10 --web.listen-address=0.0.0.0:$RM_API_PORT
