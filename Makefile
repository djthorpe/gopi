# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOBIN=$(GOPATH)/bin

# App parameters
GOPI=github.com/djthorpe/gopi
GOLDFLAGS += -X $(GOPI).GitTag=$(shell git describe --tags)
GOLDFLAGS += -X $(GOPI).GitBranch=$(shell git name-rev HEAD --name-only --always)
GOLDFLAGS += -X $(GOPI).GitHash=$(shell git rev-parse HEAD)
GOLDFLAGS += -X $(GOPI).GoBuildTime=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
GOFLAGS = -ldflags "-s -w $(GOLDFLAGS)" 

all: test install

install: helloworld hw

clean: 
	$(GOCLEAN)

test: 
	$(GOTEST) -v ./...

helloworld:
	$(GOBUILD) $(GOFLAGS) -o $(GOBIN)/helloworld ./cmd/helloworld/...

hw:
	$(GOBUILD) $(GOFLAGS) -o $(GOBIN)/hw ./cmd/hw/...

