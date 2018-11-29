#!/bin/bash

go get -u github.com/kardianos/govendor
$GOPATH/bin/govendor sync
#CGO_ENABLED=0 go build -a -v -ldflags '-s'
go build -gcflags '-N -l'
# Run unit tests and generate code coverage reports -- an html one
# for local viewing and a cobertura one for jenkins builds.
go get -u github.com/t-yuki/gocover-cobertura
kubectl proxy --port=8080&
go test -coverprofile coverage.txt
result=$?
go tool cover -html=coverage.txt -o coverage.html
$GOPATH/bin/gocover-cobertura < coverage.txt > coverage.xml
exit $result

