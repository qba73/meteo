.PHONY: check cover test tidy vet

ROOT := $(PWD)
GO_HTML_COV := ./coverage.html
GO_TEST_OUTFILE := ./c.out


test: ## Run unit tests locally
	go test -shuffle=on -race -v ./...

vet: ## Run go vet and shadow
	go vet ./...

check: ## Run static check analyzer
	staticcheck ./...

cover: ## Run unit tests and generate test coverage report
	go test -shuffle=on -race -v ./... -count=1 -cover -covermode=atomic -coverprofile=coverage.out
	go tool cover -html coverage.out
	staticcheck ./...

# MODULES
tidy: ## Run go mod tidy and vendor
	go mod tidy
	go mod vendor
