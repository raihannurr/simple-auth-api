
GO_PACKAGES ?= $(shell go list ./... | grep -v 'mock' | grep -v 'cmd')

test:
	@go test -v ${GO_PACKAGES}

coverage:
	@go test -cover -coverprofile=coverage.out ${GO_PACKAGES}
	@go tool cover -func=coverage.out

run:
	go run cmd/rest-api/server.go

build:
	go build -o bin/rest-api cmd/rest-api/server.go

start:build
	@bin/rest-api
