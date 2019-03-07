.PHONY: setup
setup: ## Install all the build and lint dependencies
	go get -u github.com/mattn/goveralls
	go get -u golang.org/x/tools/cmd/cover
	go get -t -v ./...

.PHONY: cover
cover: ## Run all the tests with race detection and opens the coverage report
	go test  ./... -coverprofile=coverage.out -race -timeout=5s
	go tool cover -html=coverage.out

