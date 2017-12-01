BIN=hurl
HEAD=$(shell git describe --tags 2> /dev/null  || git rev-parse --short HEAD)
TIMESTAMP=$(shell date '+%Y-%m-%dT%H:%M:%S')
DEPLOYMENT_PATH=s3://myon-deployment/pre-release/$(BIN)/$(BIN)-$(HEAD)

LDFLAGS="-X main.buildVersion=$(HEAD) -X main.buildTimestamp=$(TIMESTAMP)"

all: build

clean:
	-rm -f $(BIN) $(BIN)-*

install:
	go clean -i
	go install -ldflags $(LDFLAGS)

build:
	go build -ldflags $(LDFLAGS)


