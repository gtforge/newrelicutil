GOLANGCI_LINT_VERSION = 1.52.2


help:
	@echo "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' | column -c2 -t -s :)"
.PHONY: help

dev-deps: lint-install ## Install essential dev tools locally
	@echo "+ $@"
	go install -v gotest.tools/gotestsum@v1.8.2
	go install -v github.com/daixiang0/gci@v0.9.0
	go install -v github.com/golang/mock/mockgen@v1.7.0-rc.1
	go install -v golang.org/x/tools/cmd/goimports@v0.3.0
	go install mvdan.cc/gofumpt@v0.4.0
.PHONY: dev-deps

lint-install: ## Install linter locally
	@echo "+ $@"
	go install "github.com/golangci/golangci-lint/cmd/golangci-lint@v$(GOLANGCI_LINT_VERSION)"
.PHONY: lint-install

lint: ## Run local linter
	@echo "+ $@"
	golangci-lint run
.PHONY: lint

format: ## Format code
	@echo "+ $@"
	find . -type f -name '*.go' -not -path "./vendor/*" -not -path "*_mock.go" \
		-exec goimports -local "github.com/gtforge" -w {} \; \
		-exec gci write -s standard -s default -s "Prefix(github.com/gtforge)" -s "Prefix(github.com/gtforge/newrelicutil)" {} \;
	gofumpt -l -w .
.PHONY: format

test: ## Run tests
	go test -race ./...
.PHONY: test
