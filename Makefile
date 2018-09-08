# Go parameters
GOCMD=go
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
    
all: test install

install:
	$(GOINSTALL) -x ./cmd/...

test: 
	$(GOTEST) -v ./...

clean: 
	$(GOCLEAN)

