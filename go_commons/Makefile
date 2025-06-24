GO ?= go
GOBIN ?= $$($(GO) env GOPATH)/bin
GOLANGCI_LINT ?= $(GOBIN)/golangci-lint
GOLANGCI_LINT_VERSION ?= v1.60.3

# Code tidy
.PHONY: tidy
tidy:
	go mod tidy
	go fmt ./...

.PHONY: get-golangcilint
get-golangcilint:
	test -f $(GOLANGCI_LINT) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$($(GO) env GOPATH)/bin $(GOLANGCI_LINT_VERSION)

# Runs lint on entire repo
.PHONY: lint
lint: get-golangcilint
	$(GOLANGCI_LINT) run ./...

.PHONY: install-go-test-coverage
install-go-test-coverage:
	go install github.com/vladopajic/go-test-coverage/v2@latest

# Runs tests on entire repo
.PHONY: test
test: 
	go test -timeout=50s -race -count=1 -failfast  -shuffle=on ./... -coverprofile=./coverage.out -covermode=atomic -coverpkg=./...

# Runs tests on entire repo with verbose output
.PHONY: check-coverage
check-coverage: install-go-test-coverage
	go test ./... -coverprofile=./coverage.out -covermode=atomic -coverpkg=./...
	${GOBIN}/go-test-coverage --config=./.github/.testcoverage.yml

# View coverage profile
.PHONY: view-coverage
view-coverage:
	go test ./... -coverprofile=./cover.all.profile -covermode=atomic -coverpkg=./...
	go tool cover -html=cover.all.profile -o=cover.html
	xdg-open cover.html
