#!/usr/bin/env make

.PHONY: test
test:
	go test -cover -race -count=1 ./src/...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: fmt
redis-server:
	redis-server --loglevel verbose
