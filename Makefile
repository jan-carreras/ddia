

.PHONY: test
test:
	go test -cover -race ./src/...

.PHONY: fmt
fmt:
	go fmt ./...