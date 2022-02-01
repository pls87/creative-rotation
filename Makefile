BIN := "./bin/cr"
REPO=github.com/pls87/creative-rotation

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X '${REPO}/cmd/calendar/cmd.Release=develop' -X '${REPO}/cmd/calendar/cmd.BuildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S)' -X '${REPO}/cmd/calendar/cmd.GitHash=$(GIT_HASH)'

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd

lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: lint-deps
	golangci-lint run ./...

test:
	go test -race ./internal/...

build-img-migrations:
	docker build migrations/.

build-img-api:
	docker build -f build/api/Dockerfile .

build-img-stats:
	docker build -f build/stats/Dockerfile .

build-img-integration:
	docker build -f build/integration/Dockerfile .

run-docker-api-with-tool: build-img-api build-img-stats build-img-migrations
	./scripts/run-api-with-tool.sh

run-docker-integration-test: build-img-api build-img-stats build-img-migrations build-img-integration
	./scripts/run-integration-test.sh

.PHONY: build run version test lint