.PHONY: setup
setup: ## Install all the build and lint dependencies
	go mod download

.PHONY: verify
verify: ## Verify module
	go mod tidy
	go mod verify

.PHONY: test
test: ## Run all the tests
	echo 'mode: atomic' > coverage.txt && go test -covermode=atomic -coverprofile=coverage.txt -v -race -timeout=30s ./...

.PHONY: cover
cover: test ## Run all the tests and opens the coverage report
	go tool cover -html=coverage.txt

.PHONY: ci
ci: ## Run all the tests and code checks 
	verify
	lint
	test

.PHONY: build
build: ## Build a version
	go build -v ./...

.PHONY: clean
clean: ## Remove temporary files
	go clean

.DEFAULT_GOAL := build