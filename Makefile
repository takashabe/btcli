SUBPACKAGES       := $(shell go list ./...)
APP_MAIN          := cmd/btcli/btcli.go
BRANCH_NAME       := $(shell git symbolic-ref --short HEAD)
CURRENT_TIMESTAMP := $(shell date +%Y-%m-%d-%H%M%S)
VERSION           := $(CURRENT_TIMESTAMP).$(subst /,_,$(BRANCH_NAME))
LDFLAGS           := -ldflags='-s -w -X "main.Version=$(VERSION)"'

.DEFAULT_GOAL := help

##### Operation

build: $(APP_MAIN) ## Build application
	go build -a $(LDFLAGS) $(APP_MAIN)

run: $(APP_MAIN) ## Run application
	go run $(APP_MAIN)

##### Development

.PHONY: deps test vet lint

deps: ## Setup dependencies package
	dep ensure

test: ## Run go test
	go test -v $(SUBPACKAGES)

vet: ## Check go vet
	go vet $(SUBPACKAGES)

lint: ## Check golint
	golint $(SUBPACKAGES)

generate: ## Run go generate
	go generate $(SUBPACKAGES)

test-data: ## Initialize test data
	./_tools/setup_bt.sh test-project test-instance dummy

##### Utilities

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
