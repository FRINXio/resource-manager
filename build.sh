#! /bin/sh

GITHUB_TOKEN=$(grep GITHUB_TOKEN .env | cut -d '=' -f2)

# symphony is a private repo, you need a github access token to access it
export GOPRIVATE="github.com/facebookincubator"
git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/facebookincubator".insteadOf "https://github.com/facebookincubator"

go generate ./pools/...
go generate ./ent
go generate ./graph/graphql

echo ""
echo "------> Building"
go build -o ./resourceManager

echo ""
echo "------> Testing"
go test ./pools/...
