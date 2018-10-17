.PHONY: setup
setup: ## Install all the build and lint dependencies
	go mod download

.PHONY: verify
verify: ## Verify module
	go mod tidy
	go mod verify

.PHONY: test
test: ## Run all the tests
	echo 'mode: atomic' > coverage.out && go test  ./... -timeout=5s

.PHONY: cover
cover: ## Run all the tests with race detection and opens the coverage report
	echo 'mode: atomic' > coverage.out && go test  ./... -coverprofile=coverage.out -race -timeout=5s
	go tool cover -html=coverage.out

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