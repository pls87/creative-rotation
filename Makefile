BIN := "./bin/cr"
REPO=github.com/pls87/creative-rotation

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X '${REPO}/cmd/commands.Release=develop' -X '${REPO}/cmd/commands.BuildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S)' -X '${REPO}/cmd/commands.GitHash=$(GIT_HASH)'

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd

run:
	$(BIN) server --config ./configs/sample.toml & $(BIN) update_stats --config ./configs/sample.toml

lint: lint-deps
	golangci-lint run ./...

test:
	go test -v -race -count 100 ./...

lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

run-database-rabbit: build-img-api build-img-migrations
	./scripts/run-database-rabbit.sh

run-integration-test: build-img-api build-img-stats build-img-migrations build-img-integration
	./scripts/run-integration-test.sh

run-docker-api-with-tool: build-img-api build-img-stats build-img-migrations
	./scripts/run-api-with-tool.sh

run-api: build
	$(BIN) server --config ./configs/sample.toml

run-stats-updater: build
	$(BIN) update_stats --config ./configs/sample.toml

build-img-migrations:
	docker build --no-cache -t cr:migrations migrations/.

build-img-api:
	docker build -t cr:api -f build/api/Dockerfile .

build-img-stats:
	docker build -t cr:stats -f build/stats/Dockerfile .

build-img-integration:
	docker build --no-cache -t cr:integration-tests -f build/integration/Dockerfile .

.PHONY: build run test lint run-api run-stats-updater