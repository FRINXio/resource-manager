#! /bin/sh

set -xe
./build_strategies.sh
./build.sh
go test -v -short ./pools/...
