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

all: hw helloworld argonone douglas dnsregister rpcping mediakit
	@echo Use "make debian" to release to packaging

clean: 
	rm -fr $(BUILDDIR)
	$(GO) clean

# Darwin anticipates additional libraries installed via homebrew
darwin:
ifeq ($(shell test -d /usr/local/lib/pkgconfig; echo $$?),0)
	@echo "Targetting darwin"
	$(eval PKG_CONFIG_PATH += /usr/local/lib/pkgconfig)
endif

# Raspberry Pi anticipates additional libraries in /opt/vc
rpi:
ifeq ($(shell test -d /opt/vc/lib/pkgconfig; echo $$?),0)
	@echo "Targetting rpi"
	$(eval TAGS += rpi)
	$(eval PKG_CONFIG_PATH += /opt/vc/lib/pkgconfig)
endif

# MMAL package
mmal: rpi
	$(eval MMAL = $(shell PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" pkg-config --silence-errors --modversion mmal))
ifneq ($strip $(MMAL)),)
	@echo "Targetting mmal"
	$(eval TAGS += mmal)
endif


# EGL package
egl: rpi
	$(eval EGL = $(shell PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" pkg-config --silence-errors --modversion brcmegl))
ifneq ($strip $(EGL)),)
	@echo "Targetting egl"
	$(eval TAGS += egl)
endif

# Freetype package
freetype: darwin rpi
	$(eval FT = $(shell PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" pkg-config --silence-errors --modversion freetype2))
ifneq ($strip $(FT)),)
	@echo "Targetting freetype2"
	$(eval TAGS += freetype)
endif

# FFmpeg package
ffmpeg: darwin
	$(eval FT = $(shell PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" pkg-config --silence-errors --modversion libavcodec))
ifneq ($strip $(FT)),)
	@echo "Targetting ffmpeg"
	$(eval TAGS += ffmpeg)
endif

# Chromaprint package
chromaprint: darwin
	$(eval FT = $(shell PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" pkg-config --silence-errors --modversion libchromaprint))
ifneq ($strip $(FT)),)
	@echo "Targetting chromaprint"
	$(eval TAGS += chromaprint)
endif

# Create build
builddir:
	install -d $(BUILDDIR)

# Make debian packages
debian: builddir argonone dnsregister douglas nfpm
	$(eval VERSION = $(shell git describe --tags))
	$(eval ARCH = $(shell $(GO) env GOARCH))
	$(eval PLATFORM = $(shell $(GO) env GOOS))

	@sed \
		-e 's/^version:.*$$/version: $(VERSION)/'  \
		-e 's/^arch:.*$$/arch: $(ARCH)/' \
		-e 's/^platform:.*$$/platform: $(PLATFORM)/' \
		etc/nfpm/argonone.yaml > $(BUILDDIR)/argonone.yaml
	@nfpm pkg -f $(BUILDDIR)/argonone.yaml --packager deb --target $(BUILDDIR)

	@sed \
		-e 's/^version:.*$$/version: $(VERSION)/'  \
		-e 's/^arch:.*$$/arch: $(ARCH)/' \
		-e 's/^platform:.*$$/platform: $(PLATFORM)/' \
		etc/nfpm/dnsregister.yaml > $(BUILDDIR)/dnsregister.yaml
	@nfpm pkg -f $(BUILDDIR)/dnsregister.yaml --packager deb --target $(BUILDDIR)

	@sed \
		-e 's/^version:.*$$/version: $(VERSION)/'  \
		-e 's/^arch:.*$$/arch: $(ARCH)/' \
		-e 's/^platform:.*$$/platform: $(PLATFORM)/' \
		etc/nfpm/douglas.yaml > $(BUILDDIR)/douglas.yaml
	@nfpm pkg -f $(BUILDDIR)/douglas.yaml --packager deb --target $(BUILDDIR)

	@echo
	@ls -1 $(BUILDDIR)/*.deb
	@echo
	@echo "Use sudo dpkg -i <package> to install"
	@echo

# Commands
helloworld: builddir
	PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" $(GO) build -o ${BUILDDIR}/helloworld -tags "$(TAGS)" ${GOFLAGS} ./cmd/helloworld

hw: rpi darwin freetype
	PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" $(GO) build -o ${BUILDDIR}/hw -tags "$(TAGS)" ${GOFLAGS} ./cmd/hw

argonone: builddir protogen rpi
	PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" $(GO) build -o ${BUILDDIR}/argonone -tags "$(TAGS)" ${GOFLAGS} ./cmd/argonone

douglas: builddir rpi
	PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" $(GO) build -o ${BUILDDIR}/douglas -tags "$(TAGS)" ${GOFLAGS} ./cmd/douglas

dnsregister: builddir
	PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" $(GO) build -o ${BUILDDIR}/dnsregister -tags "$(TAGS)" ${GOFLAGS} ./cmd/dnsregister

rpcping: builddir protogen
	PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" $(GO) build -o ${BUILDDIR}/rpcping -tags "$(TAGS)" ${GOFLAGS} ./cmd/rpcping

mediakit: builddir ffmpeg chromaprint
	PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" $(GO) build -o ${BUILDDIR}/mediakit -tags "$(TAGS)" ${GOFLAGS} ./cmd/mediakit

# Build rules - dependencies
nfpm:
	$(GO) get github.com/goreleaser/nfpm/cmd/nfpm

protogen: protoc
	$(GO) get google.golang.org/protobuf/cmd/protoc-gen-go
	$(GO) get google.golang.org/grpc/cmd/protoc-gen-go-grpc
	$(GO) generate -x ./rpc

protoc:
ifeq ($(shell which protoc),)
	@echo apt install protobuf-compiler
	$(error protoc is not installed)
endif
