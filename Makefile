PACKAGE = $(shell go list -m)
GIT_COMMIT_HASH = $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
GIT_VERSION = $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME = $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
BINARY_NAME = mcp-prometheus
LD_FLAGS = -s -w \
	-X '$(PACKAGE)/pkg/version.CommitHash=$(GIT_COMMIT_HASH)' \
	-X '$(PACKAGE)/pkg/version.Version=$(GIT_VERSION)' \
	-X '$(PACKAGE)/pkg/version.BuildTime=$(BUILD_TIME)' \
	-X '$(PACKAGE)/pkg/version.BinaryName=$(BINARY_NAME)'

# NPM version should not append the -dirty flag
GIT_TAG_VERSION ?= $(shell echo $(shell git describe --tags --always 2>/dev/null || echo "0.0.0") | sed 's/^v//')
OSES = darwin linux windows
ARCHS = amd64 arm64

CLEAN_TARGETS :=

.DEFAULT_GOAL := help

.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9\/-]+:.*?##/ { printf "  \033[36m%-21s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: build
build: ## Build the binary
	CGO_ENABLED=0 go build -ldflags "$(LD_FLAGS)" -o $(BINARY_NAME) ./cmd/mcp-prometheus

.PHONY: build-all-platforms
build-all-platforms: ## Build for all platforms
	$(foreach os,$(OSES),$(foreach arch,$(ARCHS), \
		CGO_ENABLED=0 GOOS=$(os) GOARCH=$(arch) go build -ldflags "$(LD_FLAGS)" \
			-o dist/$(BINARY_NAME)-$(os)-$(arch)$(if $(findstring windows,$(os)),.exe,) \
			./cmd/mcp-prometheus; \
	))

.PHONY: clean
clean: ## Clean build artifacts
	rm -rf $(BINARY_NAME) dist/ $(CLEAN_TARGETS)

.PHONY: tidy
tidy: ## Tidy go modules
	go mod tidy

.PHONY: fmt
fmt: ## Format code
	go fmt ./...

.PHONY: vet
vet: ## Vet code
	go vet ./...

.PHONY: test
test: ## Run tests
	go test -v ./...

# Include build configuration files
-include build/*.mk
