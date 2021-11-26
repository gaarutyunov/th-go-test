BUILD = $(CURDIR)/build
LINT_FILE = $(CURDIR)/lint.toml
CLIENT_NAME = client
SERVER_NAME = server
GOBIN=$(shell pwd)/bin

.PHONY: help
help: ## Show help dialog
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

.PHONY: install
install: ## Install dependencies
	go get -d github.com/mgechev/revive;

.PHONY: build-client
build-client: ## Build client project
	go build -o $(BUILD)/$(CLIENT_NAME) $(CURDIR)/cmd/$(CLIENT_NAME)/main.go

.PHONY: build-server
build-server: ## Build server project
	go build -o $(BUILD)/$(SERVER_NAME) $(CURDIR)/cmd/$(SERVER_NAME)/main.go

.PHONY: clean
clean: ## Clean project
	go clean; \
	rm -rf $(BUILD);

.PHONY: fmt
fmt: ## Format project
	go fmt $(CURDIR)/...

.PHONY: lint
lint: ## Lint project
	revive -config $(LINT_FILE) -formatter friendly $(CURDIR)/...

.PHONY: check
check: ## Format and lint project
check: fmt lint

.PHONY: tidy
tidy: ## Tidy up go modules
	go mod tidy

.PHONY: setup
setup: ## Setup project
setup: clean install tidy


