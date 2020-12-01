#! /bin/sh

set -xe
./build.sh
go test -v -short ./pools/...
