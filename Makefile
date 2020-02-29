# -*- makefile -*-

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
	echo "Coverage Analysis:"
	go tool cover -func=coverage.out | perl -p -e '$$_=join(",",split())."\n"' | column -t -s,

build:
	go build -v -o ramsql main.go
