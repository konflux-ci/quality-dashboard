GO := GOFLAGS="-mod=vendor" go
SERVER_RUNTIME := $(addprefix bin/, server-runtime)
PKG := github.com/konflux-ci/quality-dashboard
GIT_COMMIT := $(or $(SOURCE_GIT_COMMIT),$(shell git rev-parse --short HEAD))
BUILD_DATE := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
TAGS := -tags "json1"

# -race is only supported on linux/amd64, linux/ppc64le, linux/arm64, freebsd/amd64, netbsd/amd64, darwin/amd64 and windows/amd64
ifeq ($(shell go env GOARCH),s390x)
TEST_RACE :=
else
TEST_RACE := -race
endif
$(SERVER_RUNTIME): server-runtime_version_flags=-ldflags "-X '$(PKG)/api/server/router/version.GitCommit=$(GIT_COMMIT)' -X '$(PKG)/api/server/router/version.BuildDate=$(BUILD_DATE)'"
$(SERVER_RUNTIME):
	go mod vendor
	$(GO) build $(server-runtime_version_flags) $(extra_flags) $(TAGS) -o $@ cmd/server/main.go
.PHONY: build
build: clean $(SERVER_RUNTIME)
.PHONY: cross
cross: server-runtime_version_flags=-ldflags "-X '$(PKG)/api/server/router/version.GitCommit=$(GIT_COMMIT)' -X '$(PKG)/api/server/router/version.BuildDate=$(BUILD_DATE)'"
cross:
ifeq ($(shell go env GOARCH),amd64)
	go mod vendor
	GOOS=darwin CC=o64-clang CXX=o64-clang++ CGO_ENABLED=1 $(GO) -mod=readonly build $(server-runtime_version_flags) -o "bin/darwin-amd64-server-runtime" --ldflags "-extld=o64-clang" cmd/server/main.go
endif
.PHONY: clean
clean:
	@rm -rf ./bin
.PHONY: static
static: extra_flags=-ldflags '-w -extldflags "-static"' -tags "json1"

.PHONY: swagger
swagger:
	go install github.com/swaggo/swag/cmd/swag@latest
	go get github.com/swaggo/swag/gen@latest
	go get github.com/swaggo/swag/cmd/swag@latest
	$$(go env GOPATH)/bin/swag init --output docs/swagger --parseVendor -g api/server/server.go
