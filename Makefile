GO_PACKAGES ?= $(shell go list ./... | grep -v 'mock' | grep -v 'cmd')
export RELEASE_VERSION    ?= $(shell git show -q --format=%h)

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

install-dependencies:
	go get -u ./...
	go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@v4

db-migrate:
	@. ./.env && \
	migrate -path internal/db/migrations/ -database "$$DB_ADAPTER://$$DB_USER:$$DB_PASSWORD@tcp($$DB_HOST:$$DB_PORT)/$$DB_NAME" up

db-rollback:
	@. ./.env && \
	migrate -path internal/db/migrations/ -database "$$DB_ADAPTER://$$DB_USER:$$DB_PASSWORD@tcp($$DB_HOST:$$DB_PORT)/$$DB_NAME" down 1

docker-build:
	docker build -t github.com/raihannurr/simple-auth-api:$(RELEASE_VERSION) -f ./deployments/rest-api/Dockerfile .