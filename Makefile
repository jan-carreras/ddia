#!/usr/bin/env make

help: ## Show help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: ## Run all tests
	go test -cover -race -count=1 ./src/...

.PHONY: fmt
fmt: ## Formats project
	go fmt ./...

.PHONY: lint
lint: ## Run the linter
	golangci-lint run

.PHONY: fmt
redis-server: ## Starts a real redis server
	redis-server --loglevel verbose

generate: ## Run go generate
	go generate ./...

commands.md: ## Generates the Commands available in Redis server
	python3 scripts/redis-commands.py 1.0.0 > commands.md 

