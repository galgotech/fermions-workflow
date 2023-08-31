## This is a self-documented Makefile. For usage information, run `make help`:
##
## For more information, refer to https://suva.sh/posts/well-documented-makefiles/

include .bingo/Variables.mk

.PHONY: build build-server build-worker build-standalone build-standalone-wasm

GO = go

##@ Building
build: build-server build-worker build-standalone build-standalone-wasm ## Build Fermions Workflow Server, Worker, Standalone, and Standalone WASM.

build-server: ## Build Fermions Workflow Server.
	@echo "build server"
	$(WIRE) gen ./pkg/server
	$(GO) build -o ./bin/fermions-workflow-server ./cmd/workflow-runtime/server

build-worker: ## Build Fermions Workflow Worker.
	$(WIRE) gen ./pkg/worker
	$(GO) build -o ./bin/fermions-workflow-worker ./cmd/workflow-runtime/worker

build-standalone: ## Build Fermions Workflow Standalone.
	$(WIRE) gen ./pkg/standalone
	$(GO) build -o ./bin/fermions-workflow-standalone ./cmd/workflow-runtime/standalone

build-standalone-wasm: ## Build Fermions Workflow Standalone Web assembly.
	$(WIRE) gen ./pkg/standalonewasm
	GOOS=js GOARCH=wasm $(GO)  build -o ./bin/fermions-workflow-standalone.wasm ./cmd/workflow-runtime/standalonewasm


##@ Tests
test: test-go ## Run all tests.

test-go: ## Run go tests.
	$(GO) test -short -covermode=atomic -timeout=30m ./pkg/...

test-coverage-go: ## Run go coverage.
	$(GO) test -coverprofile=coverage.out ./pkg/...
	$(GO) tool cover -html=coverage.out
	rm coverage.out

##@ Helpers
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
