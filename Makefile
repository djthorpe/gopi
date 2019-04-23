# Go parameters
GOCMD=go
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# App parameters
GOPI=github.com/djthorpe/gopi
GOLDFLAGS += -X $(GOPI).GitTag=$(shell git describe --tags)
GOLDFLAGS += -X $(GOPI).GitBranch=$(shell git name-rev HEAD --name-only --always)
GOLDFLAGS += -X $(GOPI).GitHash=$(shell git rev-parse HEAD)
GOLDFLAGS += -X $(GOPI).GoBuildTime=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
GOFLAGS = -ldflags "-s -w $(GOLDFLAGS)" 

all: test install

install: helloworld tasks timer

helloworld:
	$(GOINSTALL) $(GOFLAGS) ./cmd/helloworld/...

tasks:
	$(GOINSTALL) $(GOFLAGS) ./cmd/tasks/...

timer:
	$(GOINSTALL) $(GOFLAGS) ./cmd/timer/...

test: 
	$(GOTEST) -v ./...

clean: 
	$(GOCLEAN)

