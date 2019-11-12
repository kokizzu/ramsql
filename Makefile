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
	go test -vet=off -coverprofile=coverage.out ./...

test-coverage:
	@go tool cover -func=coverage.out \
		| perl -p -e 's/([^\s])\s+([^\s])/$$1 $$2/g' \
		| column -t
