GOLANGCI_VERSION := v2.4.0
GO_TOOLCHAIN     := go1.25.1
GOBIN            := $(shell go env GOPATH)/bin
LINTER           := $(GOBIN)/golangci-lint

lint-tools:
	GOBIN=$(GOBIN) GOTOOLCHAIN=$(GO_TOOLCHAIN) \
		go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_VERSION)

lint: lint-tools
	$(LINTER) version -v
	go fmt ./...
	$(LINTER) run


test:
	go test ./...

test-short:
	go test -short ./...

generate:
	go generate ./...