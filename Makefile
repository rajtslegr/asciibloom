.PHONY: all build install-tools lint fmt vet test check clean

BINARY_NAME=asciibloom
GO_FILES=$(shell find . -name '*.go' -type f)
GOPATH=$(shell go env GOPATH)
GOFUMPT=$(GOPATH)/bin/gofumpt
GOLANGCI_LINT=$(GOPATH)/bin/golangci-lint

all: check build

build:
	go build -o $(BINARY_NAME) .

install-tools:
	@echo "Installing linting tools..."
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	go install mvdan.cc/gofumpt@latest
	@echo "Tools installed to $(GOPATH)/bin"
	@echo "Add $(GOPATH)/bin to your PATH or use: export PATH=\$$PATH:$(GOPATH)/bin"

fmt:
	@echo "Running gofumpt..."
	$(GOFUMPT) -w .

fmt-check:
	@echo "Checking formatting..."
	@test -z "$(shell gofmt -l .)" || (echo "Please run 'make fmt' to fix formatting issues:" && gofmt -l . && exit 1)
	@echo "Formatting OK"

vet:
	@echo "Running go vet..."
	go vet ./...

test:
	@echo "Running tests with race detector..."
	go test -race ./...

lint: fmt-check vet
	@echo "Running golangci-lint..."
	$(GOLANGCI_LINT) run ./...

check: fmt-check vet test
	@echo "All checks passed!"

clean:
	rm -f $(BINARY_NAME)
	rm -rf .golangci-lint-cache/

help:
	@echo "Available targets:"
	@echo "  all            - Run checks and build (default)"
	@echo "  build          - Build the binary"
	@echo "  check          - Run all checks (fmt, vet, test)"
	@echo "  install-tools  - Install linting tools (golangci-lint, gofumpt)"
	@echo "  fmt            - Format code with gofumpt"
	@echo "  fmt-check      - Check if code is formatted"
	@echo "  vet            - Run go vet"
	@echo "  test           - Run tests with race detector"
	@echo "  lint           - Run fmt-check, vet, and golangci-lint"
	@echo "  clean          - Clean build artifacts"
	@echo "  help           - Show this help"
