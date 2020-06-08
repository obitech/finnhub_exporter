GO 			?= go
GOFMT 		?= $(go)fmt
BIN_DIR 	?= $(shell pwd)/bin/
BIN_NAME 	?= finnhub_exporter
GOOS		?= darwin
GOARCH		?= amd64

.PHONY: help
## help: prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: build
## build: builds finnhub_exporter
build:
	GO111MODULE=on CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -o $(BIN_DIR)/$(BIN_NAME) main.go log.go

.PHONY: run
## build: builds finnhub_exporter
run:
	GO111MODULE=on CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) run main.go log.go