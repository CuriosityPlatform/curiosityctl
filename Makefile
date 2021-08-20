# enable buildkit
export DOCKER_BUILDKIT=1

all: build test check

.PHONY: build
build: modules
	@docker build . --target ctl \
	--progress tty \
	--output ./bin

.PHONY: modules
modules:
	@docker build . --target go-mod-tidy \
	--progress tty \
	--output .

.PHONY: test
test:
	go test ./...

.PHONY: check
check:
	@docker build . --target lint

.PHONY: cache-clear
cache-clear: ## Clear the builder cache
	@docker builder prune --force --filter type=exec.cachemount --filter=unused-for=24h