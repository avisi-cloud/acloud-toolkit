BINARY = acloud-toolkit
GOARCH = amd64

COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

BUILD=0
VERSION=0.0.0
IMAGE=ame/acloud-toolkit

ifneq (${BRANCH}, release)
	BRANCH := -${BRANCH}
else
	BRANCH :=
endif

PKG_LIST := $(shell go list ./... | grep -v /vendor/)
LDFLAGS = -ldflags "-X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.Branch=${BRANCH}"

.PHONY: all
all: clean build

.PHONY: linux
linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o bin/${BINARY}-linux-${GOARCH} . ;

.PHONY: darwin
darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o bin/${BINARY}-darwin-${GOARCH} . ;

.PHONY: windows
windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=${GOARCH} go build ${LDFLAGS} -o bin/${BINARY}-windows-${GOARCH}.exe . ;

.PHONY: build
build:
	CGO_ENABLED=0 go build ${LDFLAGS} -o bin/${BINARY} . ;
	chmod +x bin/${BINARY};

.PHONY: lint
lint: ## Lint the files
	@golint -set_exit_status ${PKG_LIST}

.PHONY: test
test: ## Run unittests
	@go test -short ${PKG_LIST}

.PHONY: race
race: ## Run data race detector
	@go test -race -short ${PKG_LIST}

.PHONY: msan
msan: ## Run memory sanitizer
	@go test -msan -short ${PKG_LIST}

.PHONY: fmt
fmt:
	@go fmt ${PKG_LIST};

.PHONY: review
review:
	reviewdog -diff="git diff FETCH_HEAD" -tee

.PHONY: docs
docs:
	go run tools/docs.go

.PHONY: docker
docker:
	docker build -t registry.avisi.cloud/library/${BINARY}:dev .

.PHONY: clean
clean:
	-rm -f bin/${BINARY}-* bin/${BINARY}

.PHONY: goreleaser-release-snapshot
goreleaser-release-snapshot: ## Build and run release in snapshot mode
	goreleaser release --snapshot --clean

.PHONY: goreleaser-build
goreleaser-build: ## Build and run release in snapshot mode
	goreleaser build --snapshot --single-target --clean