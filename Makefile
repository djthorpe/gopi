# Go parameters
GO=go

# App parameters
GOPI=github.com/djthorpe/gopi/v3/pkg/config
GOLDFLAGS += -X $(GOPI).GitTag=$(shell git describe --tags)
GOLDFLAGS += -X $(GOPI).GitBranch=$(shell git name-rev HEAD --name-only --always)
GOLDFLAGS += -X $(GOPI).GitHash=$(shell git rev-parse HEAD)
GOLDFLAGS += -X $(GOPI).GoBuildTime=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
GOFLAGS = -ldflags "-s -w $(GOLDFLAGS)" 

all:
	@echo "Synax: make linux|darwin|rpi|test|clean"

# Build for different platforms
linux: TAGS = -tags linux
linux: test install

darwin: TAGS = -tags darwin
#darwin: testrace install
darwin: install

rpi: TAGS = -tags "rpi egl freetype"
rpi: PKG_CONFIG_PATH = /opt/vc/lib/pkgconfig
rpi: test install

# Build rules
protogen:
	$(GO) generate -x ./pkg/rpc

testrace: protogen
	$(GO) clean -testcache
	PKG_CONFIG_PATH="${PKG_CONFIG_PATH}" $(GO) test $(TAGS) -race ./pkg/...

test: protogen
	$(GO) clean -testcache
	PKG_CONFIG_PATH="${PKG_CONFIG_PATH}" $(GO) test -count 5 $(TAGS) ./pkg/...

install: protogen
	PKG_CONFIG_PATH="${PKG_CONFIG_PATH}" $(GO) install $(TAGS) ${GOFLAGS} ./cmd/...

clean: 
	$(GO) clean
