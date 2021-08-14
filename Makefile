export APP_CMD_NAME = curiosity

all: build test check

.PHONY: build
build: modules
	bin/go-build.sh "cmd/$(APP_CMD_NAME)" "bin/$(APP_CMD_NAME)" $(APP_CMD_NAME) .env

.PHONY: modules
modules:
	go mod tidy

.PHONY: test
test:
	go test ./...

.PHONY: check
check:
	golangci-lint run