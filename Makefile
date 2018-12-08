SUBPACKAGES := $(shell go list ./...)
APP_MAIN    := cmd/btcli/btcli.go
VERSION     := $(shell git describe --tags --abbrev=0)
REVISION    := $(shell git rev-parse --short HEAD)
LDFLAGS     := -X 'main.Version=$(VERSION)' \
               -X 'main.Revision=$(REVISION)'

.DEFAULT_GOAL := help

##### Operation

.PHONY: build run

build: $(APP_MAIN) ## Build application
	go build -a -o bin/btcli -ldflags "$(LDFLAGS)" $(APP_MAIN)

run: $(APP_MAIN) ## Run application
	go run $(APP_MAIN)

##### Development

.PHONY: deps test vet lint generate testdata

test: ## Run go test
	go test -v $(SUBPACKAGES)

vet: ## Check go vet
	go vet $(SUBPACKAGES)

lint: ## Check golint
	golint $(SUBPACKAGES)

generate: ## Run go generate
	go generate $(SUBPACKAGES)

testdata: ## Initialize test data
	./_tools/setup_bt.sh test-project test-instance dummy

##### Utilities

.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
