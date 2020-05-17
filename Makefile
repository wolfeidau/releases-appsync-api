STAGE ?= dev
BRANCH ?= master
APP_NAME ?= realworld-appsync-api
PACKAGE ?= $(shell go mod edit -json | jq -r .Module.Path)
GIT_HASH ?= $(shell git rev-parse --short HEAD)
RAW_EVENT_LOGGING ?= false

GOLANGCI_VERSION = 1.27.0

# https://tip.golang.org/cmd/go/#hdr-Module_configuration_for_non_public_modules
#
# for all modules in the jumacloud org avoid the GOPROXY/GOSUM requests and go direct
export GOPRIVATE=github.com/jumacloud

BUILD_OVERRIDES = \
	-X $(PACKAGE)/internal/app.Name=$(APP_NAME) \
	-X $(PACKAGE)/internal/app.BuildDate=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ') \
	-X $(PACKAGE)/internal/app.Commit=$(GIT_HASH)

LDFLAGS := -ldflags='-w -s $(BUILD_OVERRIDES)' -trimpath

default: build archive package deploy
.PHONY: default

ci: clean generate lint test build archive
.PHONY: ci

bin/golangci-lint: bin/golangci-lint-${GOLANGCI_VERSION}
	@ln -sf golangci-lint-${GOLANGCI_VERSION} bin/golangci-lint
bin/golangci-lint-${GOLANGCI_VERSION}:
	@curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | BINARY=golangci-lint bash -s -- v${GOLANGCI_VERSION}
	@mv bin/golangci-lint $@

bin/mockgen:
	@env GOBIN=$$PWD/bin GO111MODULE=on go install github.com/golang/mock/mockgen

bin/gcov2lcov:
	@env GOBIN=$$PWD/bin GO111MODULE=on go install github.com/jandelgado/gcov2lcov

clean:
	@echo "--- clean all the things"
	@rm -rf dist
.PHONY: clean

generate:
	@echo "--- generate all the things"
	@go generate ./...
.PHONY: generate

lint: bin/golangci-lint generate
	@echo "--- lint all the things"
	@bin/golangci-lint run
.PHONY: lint

test: bin/gcov2lcov
	@echo "--- test all the things"
	@go test -v -covermode=count -coverprofile=coverage.txt ./pkg/... ./internal/... ./cmd/...
	@bin/gcov2lcov -infile=coverage.txt -outfile=coverage.lcov
.PHONY: test

build:
	@echo "--- build all the things"
	@go build $(LDFLAGS) -o dist/lambda-api ./cmd/lambda-api
.PHONY: build

archive:
	@echo "--- build an archive"
	@cd dist && zip -X -9 -r ./handler.zip ./service-security-events
.PHONY: archive

package:
	@echo "--- package CFN assets"
	@echo aws cloudformation package \
		--template-file sam/api/template.yaml \
		--s3-bucket $(PACKAGE_BUCKET) \
		--s3-prefix $(APP_NAME) \
		--output-template-file dist/packaged-api-template.yaml
.PHONY: package

deploy:
	@echo "--- deploy $(APP_NAME)-$(STAGE)-$(BRANCH) stack to aws"
	@echo aws cloudformation deploy \
		--template-file dist/packaged-api-template.yaml \
		--capabilities CAPABILITY_NAMED_IAM CAPABILITY_AUTO_EXPAND \
		--tags "environment=$(STAGE)" "branch=$(BRANCH)" "service=$(APP_NAME)" "owner=$(USER)" \
		--stack-name $(APP_NAME)-$(STAGE)-$(BRANCH) \
		--parameter-overrides AppName=$(APP_NAME) Stage=$(STAGE) Branch=$(BRANCH) RawEventLogging=$(RAW_EVENT_LOGGING)
.PHONY: deploy
