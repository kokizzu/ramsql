# -*- makefile -*-

export GO111MODULE=on

check: lint vet test

generate:
	go generate ./...

lint:
	golint ./...

vet:
	go vet ./...

test:
	go test -vet=off -coverprofile=coverage.out ./...

test-coverage:
	@go tool cover -func=coverage.out \
		| perl -p -e 's/([^\s])\s+([^\s])/$$1 $$2/g' \
		| column -t

build:
	go build -v -o ramsql main.go
