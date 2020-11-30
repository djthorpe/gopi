# Alias for go exec
GO=go

# Go parameters
GOPI=github.com/djthorpe/gopi/v3/pkg/config
GOLDFLAGS += -X $(GOPI).GitTag=$(shell git describe --tags)
GOLDFLAGS += -X $(GOPI).GitBranch=$(shell git name-rev HEAD --name-only --always)
GOLDFLAGS += -X $(GOPI).GitHash=$(shell git rev-parse HEAD)
GOLDFLAGS += -X $(GOPI).GoBuildTime=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
GOFLAGS = -ldflags "-s -w $(GOLDFLAGS)" 
BUILDDIR = build

all: checkdeps
	@echo "Synax: make linux|darwin|rpi|test|clean"

# Build for different platforms
linux: TAGS = -tags linux
linux: checkdeps install

darwin: TAGS = -tags "darwin ffmpeg"
darwin: PKG_CONFIG_PATH = /usr/local/lib/pkgconfig
darwin: checkdeps install

rpi: TAGS = -tags "rpi egl freetype"
rpi: checkdeps argonone

# Build rules - commands
argonone: PKG_CONFIG_PATH = /opt/vc/lib/pkgconfig
argonone: VERSION = $(shell git describe --tags)
argonone: nfpm
	install -d $(BUILDDIR)
	PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" $(GO) build -o ${BUILDDIR}/argonone $(TAGS) ${GOFLAGS} ./cmd/argonone
	sed -e 's/^version:.*$$/version: $(VERSION)/' etc/nfpm/argonone.yaml > $(BUILDDIR)/argonone.yaml
	nfpm pkg -f $(BUILDDIR)/argonone.yaml --packager deb --target $(BUILDDIR)
	@echo "Use sudo dpkg -i <package> to install"

# Build rules - dependencies
nfpm:
	$(GO) get github.com/goreleaser/nfpm/cmd/nfpm

protogen: protoc-gen-go
	$(GO) generate -x ./pkg/rpc

protoc-gen-go:
	$(GO) get github.com/golang/protobuf/protoc-gen-go

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

checkdeps:
ifndef GOBIN
	$(error GOBIN is undefined)
endif
ifeq (,$(shell which protoc))
	$(error protoc is not installed)
endif
