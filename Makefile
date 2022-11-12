.PHONY: tidy fmt lint tools

GO = go
GOFILES = $(shell find . -name "*.go")

tidy:
	$(GO) mod tidy

fmt:
	gofmt -s -l -w $(GOFILES)

lint:
	golangci-lint run .

tools:
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest