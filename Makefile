# This is a self-documenting Makefile.
# See https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

.PHONY: test
test: ## Execute all tests and show a coverage summary
	go test -coverprofile=coverage.out ./...

race: ## Execute all tests with race detector enabled
	go test -race ./...

.PHONY: coverageHTML
coverageHTML: test ## Create HTML coverage report
	go tool cover -html=coverage.out

.PHONY: bench
bench: ## Execute all tests and include benchmarks
	go test -bench=. ./...

.PHONY: help
help: ## Display this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
