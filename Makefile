# -*- makefile -*-

export GO111MODULE=on

check: lint vet test

generate:
	go generate ./...

lint: generate
	golint ./...

vet: generate
	go vet ./...

test: generate
	go test -vet=off ./...
