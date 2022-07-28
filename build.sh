#! /bin/bash
set -ex
cd `dirname $0`

GOOS=linux GOARCH=amd64 go build -v -ldflags '-w -extldflags "-static"' -o main
mkdir -p output/
mv main output/
cp run.sh output/