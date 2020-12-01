#! /bin/sh

set -xe

GITHUB_TOKEN_FILE=$(grep GITHUB_TOKEN .env | cut -d '=' -f2)
GITHUB_TOKEN=${GITHUB_TOKEN_EXTERNAL:-${GITHUB_TOKEN_FILE}}

# symphony is a private repo, you need a github access token to access it
export GOPRIVATE="github.com/FRINXio"
git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/FRINXio".insteadOf "https://github.com/FRINXio"

go generate ./pools/...
go generate ./ent
go generate ./graph/graphql

echo ""
echo "------> Building"
go build -o ./resourceManager

echo "All OK"
