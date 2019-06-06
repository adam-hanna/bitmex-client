##########
default: build
##########

GOX := $(shell which gox 2> /dev/null)
LINTER := $(shell which golangci-lint 2> /dev/null)

.PHONY: GOX/exists
GOX/exists:
ifndef GOX
	$(error "No gox in PATH, consider doing 'GO111MODULE=on go get github.com/mitchellh/gox'")
endif

.PHONY: LINTER/exists
LINTER/exists:
ifndef LINTER
	$(error "No golangci-lint in PATH, consider doing 'GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@latest'")
endif

.PHONY: deps
deps:
	@GO111MODULE=on go mod vendor

.PHONY: build
build: GOX/exists fmt lint test build/force

.PHONY: build/force
build/force: GOX/exists clean
	@gox -osarch="linux/amd64 darwin/amd64 windows/amd64" -output=".bin/{{.OS}}-{{.Arch}}/liquidity-bot" -verbose

.PHONY: clean
clean:
	@rm -rf ./.bin

.PHONY: fmt
fmt:
	@gofmt -e -s -w .

.PHONY: lint
lint: LINTER/exists
	@golangci-lint run --skip-dirs=swagger ./...

.PHONY: test
test: test/unit test/integration test/e2e
	@echo "passed all tests"

.PHONY: test/unit
test/unit:
	@go test ./... -tags=unit -v && echo "passed unit tests"

.PHONY: test/integration
test/integration:
	@GOCACHE=off go test ./... -tags=integration -v && echo "passed integration tests"

.PHONY: test/e2e
test/e2e:
	@GOCACHE=off go test ./... -tags=e2e -v && echo "passed e2e tests"

.PHONY: run
run:
	@GO111MODULE=on go run main.go
