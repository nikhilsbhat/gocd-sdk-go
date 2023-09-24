GOFMT_FILES?=$(shell find . -not -path "./vendor/*" -type f -name '*.go')
BUILD_ENVIRONMENT?=${ENVIRONMENT}
GOVERSION?=$(shell go version | awk '{printf $$3}')
APP_DIR?=$(shell git rev-parse --show-toplevel)
SOURCE_PACKAGES?=$(shell go list -mod=vendor ./... | grep -v "vendor" | grep -v "mocks")
TEST_FILES?=$(shell go list ./... | grep -v /vendor/ | grep -v examples)

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

.PHONY: help
help: ## Prints help (only for targets with comments)
	@grep -E '^[a-zA-Z0-9._-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

local.fmt: ## Lints all the go code in the application.
	@gofmt -w $(GOFMT_FILES)
	$(GOBIN)/gofumpt -l -w $(GOFMT_FILES)
	$(GOBIN)/goimports -w $(GOFMT_FILES)
	$(GOBIN)/gci write $(GOFMT_FILES) --skip-generated

local.check: local.fmt ## Loads all the dependencies to vendor directory
	@go mod vendor
	@go mod tidy

local.build: local.check ## Generates the artifact with the help of 'go build'
	GOVERSION=${GOVERSION} BUILD_ENVIRONMENT=${BUILD_ENVIRONMENT} goreleaser build --rm-dist

publish: local.check ## Builds and publishes the app
	GOVERSION=${GOVERSION} BUILD_ENVIRONMENT=${BUILD_ENVIRONMENT} PLUGIN_PATH=${APP_DIR} goreleaser release --rm-dist

mock.publish: local.check ## Builds and mocks app release
	GOVERSION=${GOVERSION} BUILD_ENVIRONMENT=${BUILD_ENVIRONMENT} PLUGIN_PATH=${APP_DIR} goreleaser release --skip-publish --rm-dist

lint: ## Lint's application for errors, it is a linters aggregator (https://github.com/golangci/golangci-lint).
	if [ -z "${DEV}" ]; then golangci-lint run --color always ; else docker run --rm -v $(APP_DIR):/app -w /app golangci/golangci-lint:v1.46.2-alpine golangci-lint run --color always ; fi

report: ## Publishes the go-report of the appliction (uses go-reportcard)
	if [ -z "${DEV}" ]; then goreportcard -v ; else docker run --rm -v $(APP_DIR):/app -w /app basnik/goreportcard-cli:latest goreportcard-cli -v ; fi

test: ## runs test cases
	@time go test $(TEST_FILES) -mod=vendor -coverprofile cover.out && go tool cover -html=cover.out -o cover.html && open cover.html

generate.mocks:  ## Generates mocks to those methods that has comments //go:generate
	@go generate ${SOURCE_PACKAGES}