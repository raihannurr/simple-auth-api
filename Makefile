GO_PACKAGES ?= $(shell go list ./... | grep -v 'mock' | grep -v 'cmd')
export RELEASE_VERSION    ?= $(shell git show -q --format=%h)

test:
	@go test -v ${GO_PACKAGES}

coverage:
	@go test -cover -coverprofile=coverage.out ${GO_PACKAGES}
	@go tool cover -func=coverage.out

run:
	go run cmd/rest-api/server.go

build: lint
	go build -o bin/rest-api cmd/rest-api/server.go

start:build
	@bin/rest-api

install-linter:
	@if ! [ -x "$$(command -v golangci-lint)" ]; then \
		echo "golangci-lint not found, installing..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.62.2; \
	else \
		echo "golangci-lint is already installed."; \
	fi

install-dependencies:install-linter
	go get -u ./...
	go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@v4

lint:
	@golangci-lint run

db-migrate:
	@. ./.env && \
	migrate -path internal/db/migrations/ -database "$$DB_ADAPTER://$$DB_USER:$$DB_PASSWORD@tcp($$DB_HOST:$$DB_PORT)/$$DB_NAME" up

db-rollback:
	@. ./.env && \
	migrate -path internal/db/migrations/ -database "$$DB_ADAPTER://$$DB_USER:$$DB_PASSWORD@tcp($$DB_HOST:$$DB_PORT)/$$DB_NAME" down 1

docker-build:
	docker build -t github.com/raihannurr/simple-auth-api:$(RELEASE_VERSION) -f ./deployments/rest-api/Dockerfile .

docker-compose-up:
	RELEASE_VERSION=$(RELEASE_VERSION) docker compose -f docker-compose.yml up

docker-compose-down:
	docker-compose down
