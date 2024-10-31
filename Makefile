SHELL := /bin/bash

.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Compile packages and dependencies
	go build ./...
	
test: ## Run the project’s tests.
	go test -v ./...

fmt: ## Format the project’s source code.
	go fmt ./...