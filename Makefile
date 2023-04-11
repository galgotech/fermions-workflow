include .bingo/Variables.mk

.PHONY: build build-server build-worker build-standalone

GO = go

build: build-server build-worker build-standalone

build-server: gen-go
	@echo "build server"
	$(GO) build -o ./bin/fermions-workflow-server ./cmd/workflow-runtime/server

build-worker: gen-go
	$(GO) build -o ./bin/fermions-workflow-worker ./cmd/workflow-runtime/worker

build-standalone: gen-go
	$(GO) build -o ./bin/fermions-workflow-standalone ./cmd/workflow-runtime/standalone

gen-go: $(WIRE)
	@echo "generate go files"
	$(WIRE) gen ./pkg/worker
	$(WIRE) gen ./pkg/server
	$(WIRE) gen ./pkg/standalone

test-go:
	$(GO) test -short -covermode=atomic -timeout=30m ./pkg/...

test-coverage-go:
	$(GO) test -coverprofile=coverage.out ./pkg/...
	$(GO) tool cover -html=coverage.out
	rm coverage.out
