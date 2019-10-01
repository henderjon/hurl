export CC_TEST_REPORTER_ID = 22b9feb0e15027c6928559bf45bbd8fbd2489be2687a008b326613e321d30c05

BIN=hurl
HEAD=$(shell ([ -n "$${CI_TAG}" ] && echo "$$CI_TAG" || exit 1) || git describe --dirty --long --tags 2> /dev/null || git rev-parse --short HEAD)
TIMESTAMP=$(shell date '+%Y-%m-%dT%H:%M:%S %z %Z')
COVERAGEOUTFILE=c.out
LDFLAGS="-X 'main.buildVersion=$(HEAD)' -X 'main.buildTimestamp=$(TIMESTAMP)' -X 'main.compiledBy=$(shell go version)'" # `-s -w` removes some debugging info that might not be necessary in production (smaller binaries)

all: local

.PHONY: clean
clean:
	rm -f $(BIN) $(BIN)-* $(COVERAGEOUTFILE)
	go clean -i ./...

.PHONY: dep
dep:
	go mod vendor

.PHONY: install
install:
	go clean -i
	go install -ldflags $(LDFLAGS)

.PHONY: build
build: $(BIN)-darwin64-$(HEAD) $(BIN)-linux64-$(HEAD)

.PHONY: local
local:
	go build -ldflags $(LDFLAGS) -o $(BIN)

$(BIN)-darwin64-$(HEAD): clean
	env GOOS=darwin GOARCH=amd64 go build -ldflags $(LDFLAGS) -o $(BIN)-darwin64-$(HEAD)

$(BIN)-linux64-$(HEAD): clean
	env GOOS=linux GOARCH=amd64 go build -ldflags $(LDFLAGS) -o $(BIN)-linux64-$(HEAD)

.PHONY: test-vendor
test-vendor:
	go test -mod=vendor -coverprofile=coverage.out -covermode=count

.PHONY: test
test:
	go test -coverprofile=coverage.out -covermode=count

.PHONY: race
race:
	go test -race

.PHONY: test-report
test-report: test
	go tool cover -html=coverage.out

.PHONY: travis
travis:
	go test -coverprofile $(COVERAGEOUTFILE) ./...

.PHONY: cclimate-linux
cclimate-linux:
	curl -L https://codeclimate.com/downloads/test-reporter/test-reporter-latest-linux-amd64 > ./cc-test-reporter
	# curl -L https://codeclimate.com/downloads/test-reporter/test-reporter-latest-darwin-amd64 > ./cc-test-reporter
	chmod +x ./cc-test-reporter
	./cc-test-reporter before-build
	go test -coverprofile $(COVERAGEOUTFILE) ./...
	./cc-test-reporter after-build --exit-code $(TRAVIS_TEST_RESULT)
