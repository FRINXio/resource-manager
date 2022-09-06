#! /bin/sh

set -xe

GITHUB_TOKEN_FILE=$(grep GITHUB_TOKEN .env | cut -d '=' -f2)
GITHUB_TOKEN=${GITHUB_TOKEN_EXTERNAL:-${GITHUB_TOKEN_FILE}}

# symphony is a private repo, you need a github access token to access it
export GOPRIVATE="github.com/FRINXio"

go generate ./pools/...
go get entgo.io/ent/cmd/ent@v0.11.3-0.20220830071904-3b1b75b9d7a9
go get entgo.io/ent/cmd/internal/printer@v0.11.3-0.20220830071904-3b1b75b9d7a9
go get -d github.com/99designs/gqlgen@v0.17.16
go generate ./ent
go generate ./graph/graphql

echo ""
echo "------> Building"
go build -o ./resourceManager

echo "All OK"
