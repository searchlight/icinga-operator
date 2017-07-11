#!/usr/bin/env bash

pushd $GOPATH/src/github.com/appscode/searchlight/hack/gendocs
go run main.go

cd $GOPATH/src/github.com/appscode/searchlight/docs/hyperalert
sed -i 's/######\ Auto\ generated\ by.*//g' *
popd
