PROJECT_NAME := "go-ipfs-stat-top"
PKG := "gitlab.com/coderlane/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

# Makefile based on Pantomath example:
# https://about.gitlab.com/2017/11/27/go-tools-and-gitlab-how-to-do-continuous-integration-like-a-boss/

.PHONY: all dep build clean test coverage coverhtml lint

all: build

lint: ## Lint the files
	@golint -set_exit_status ${PKG_LIST}

test: ## Run unittests
	@go test -short ${PKG_LIST}

coverage: ## Generate global code coverage report
	go test ./... -v -coverprofile .testCoverage.txt

dep: ## Get the dependencies
	@go get -v -d ./...
	@go get -u github.com/golang/lint/golint

build: dep ## Build the binary file
	@go build -i -v $(PKG)

clean: ## Remove previous build
	@rm -f $(PROJECT_NAME)

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
