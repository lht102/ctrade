SHELL := /bin/sh

-include .env.$(ENV)
export

.PHONY: build
build:
	@rm -rf bin/
	@mkdir bin/
	CGO_ENABLED=0 go build -ldflags "-w -s" -v -o ./bin ./cmd/ctraded

.PHONY: dep
dep:
	go mod tidy

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -v -race -cover -coverprofile=c.out -covermode=atomic ./...

.PHONY: test-all
test-all: lint test
