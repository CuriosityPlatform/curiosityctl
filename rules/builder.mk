export APP_CMD_NAME = curiosity
APP_EXECUTABLE_OUT?=bin

.PHONY: build
build: modules
	bin/go-build.sh "cmd/$(APP_CMD_NAME)" "$(APP_EXECUTABLE_OUT)/$(APP_CMD_NAME)" $(APP_CMD_NAME)

.PHONY: modules
modules:
	go mod download

.PHONY: check
check:
	golangci-lint run

