BINARY = csi-snapshot-utils
GOARCH = amd64

COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

BUILD=0
VERSION=0.0.0
IMAGE=ame/csi-snapshot-utils

ifneq (${BRANCH}, release)
	BRANCH := -${BRANCH}
else
	BRANCH :=
endif

PKG_LIST := $(shell go list ./... | grep -v /vendor/)
LDFLAGS = -ldflags "-X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.Branch=${BRANCH}"

all: link clean linux darwin

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o bin/${BINARY}-linux-${GOARCH} ./cmd/${BINARY} ;

darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o bin/${BINARY}-darwin-${GOARCH} ./cmd/${BINARY} ;

windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=${GOARCH} go build ${LDFLAGS} -o bin/${BINARY}-windows-${GOARCH}.exe ./cmd/${BINARY} ;

build:
	CGO_ENABLED=0 go build ${LDFLAGS} -o bin/${BINARY} ./cmd/${BINARY} ;
	chmod +x bin/${BINARY};

lint: ## Lint the files
	@golint -set_exit_status ${PKG_LIST}

test: ## Run unittests
	@go test -short ${PKG_LIST}

race: ## Run data race detector
	@go test -race -short ${PKG_LIST}

msan: ## Run memory sanitizer
	@go test -msan -short ${PKG_LIST}

fmt:
	@go fmt ${PKG_LIST};

install:
	go mod download;

update-deps:
	go mod vendor;

docker:
	docker build -t registry.avisi.cloud/library/${BINARY}:dev .

clean:
	-rm -f bin/${BINARY}-* bin/${BINARY}

.PHONY: link linux darwin windows test vet fmt clean
