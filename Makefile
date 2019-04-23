# Go parameters
GOCMD=go
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# App parameters
GOLDFLAGS += -X main.GitTag=$(shell git describe --tags)
GOLDFLAGS += -X main.GitBranch=$(shell git name-rev HEAD --name-only --always)
GOLDFLAGS += -X main.GitHash=$(shell git rev-parse HEAD)
GOLDFLAGS += -X main.GoBuildTime=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
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

