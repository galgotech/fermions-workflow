include .bingo/Variables.mk

.PHONY: build build-server build-worker

GO = go

build: build-server build-worker-standalone

build-server: gen-go ## Build server
	@echo "build server"
	$(GO) build -o ./bin/workflow-server ./cmd/workflow-runtime/server

build-worker-standalone: gen-go ## Build server
	$(GO) build -o ./bin/workflow-standalone ./cmd/workflow-runtime/worker

gen-go: $(WIRE)
	@echo "generate go files"
	$(WIRE) gen ./pkg/worker
	$(WIRE) gen ./pkg/server

test-go:
	$(GO) test -short -covermode=atomic -timeout=30m ./pkg/...

test-coverage-go:
	$(GO) test -coverprofile=coverage.out ./pkg/...
	$(GO) tool cover -html=coverage.out
	rm coverage.out
