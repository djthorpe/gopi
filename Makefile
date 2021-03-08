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
PACKAGECLOUD_REPO = djthorpe/gopi/raspbian/buster

all: hw httpserver helloworld argonone douglas dnsregister rpc googlecast mediakit
	@echo Use "make debian" to release to packaging
	@echo Use "make clean" to clear build cache
	@echo Use "make test" to run tests

clean: 
	rm -fr $(BUILDDIR)
	$(GO) clean
	$(GO) mod tidy

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

# Older Broadcom Graphics (dispmanx)
dispmanx: rpi
	$(eval DX = $(shell PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" pkg-config --silence-errors --modversion bcm_host))
ifneq ($strip $(DX)),)
	@echo "Targetting dispmanx"
	$(eval TAGS += dispmanx)
endif

# MMAL bindings
mmal: rpi
	$(eval MMAL = $(shell PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" pkg-config --silence-errors --modversion mmal))
ifneq ($strip $(MMAL)),)
	@echo "Targetting mmal"
	$(eval TAGS += mmal)
endif

# OpenVG bindings
openvg: rpi dispmanx
	$(eval OPENVG = $(shell PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" pkg-config --silence-errors --modversion brcmvg))
ifneq ($strip $(OPENVG)),)
	@echo "Targetting openvg"
	$(eval TAGS += openvg)
endif

# EGL bindings
egl:
	$(eval EGL = $(shell PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" pkg-config --silence-errors --modversion egl))
ifneq ($strip $(EGL)),)
	@echo "Targetting egl"
	$(eval TAGS += egl)
endif

# GBM bindings
gbm:
	$(eval GBM = $(shell PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" pkg-config --silence-errors --modversion gbm))
ifneq ($strip $(GBM)),)
	@echo "Targetting gbm"
	$(eval TAGS += gbm)
endif

# DRM bindings
drm:
	$(eval DRM = $(shell PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" pkg-config --silence-errors --modversion drm))
ifneq ($strip $(DRM)),)
	@echo "Targetting drm"
	$(eval TAGS += drm)
endif

# Freetype bindings
freetype: darwin rpi
	$(eval FT = $(shell PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" pkg-config --silence-errors --modversion freetype2))
ifneq ($strip $(FT)),)
	@echo "Targetting freetype2"
	$(eval TAGS += freetype)
endif

# FFmpeg bindings
ffmpeg: darwin
	$(eval FT = $(shell PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" pkg-config --silence-errors --modversion libavcodec))
ifneq ($strip $(FT)),)
	@echo "Targetting ffmpeg"
	$(eval TAGS += ffmpeg)
endif

# DVB bindings
dvb:
	@echo "Targetting dvb"
	$(eval TAGS += dvb)

# Chromaprint bindings
chromaprint: darwin
	$(eval FT = $(shell PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" pkg-config --silence-errors --modversion libchromaprint))
ifneq ($strip $(FT)),)
	@echo "Targetting chromaprint"
	$(eval TAGS += chromaprint)
endif

# Create build folder
builddir:
	install -d $(BUILDDIR)

# Make debian packages
debian: clean builddir argonone dnsregister douglas httpserver rpc googlecast nfpm
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

	@sed \
		-e 's/^version:.*$$/version: $(VERSION)/'  \
		-e 's/^arch:.*$$/arch: $(ARCH)/' \
		-e 's/^platform:.*$$/platform: $(PLATFORM)/' \
		etc/nfpm/httpserver.yaml > $(BUILDDIR)/httpserver.yaml
	@nfpm pkg -f $(BUILDDIR)/httpserver.yaml --packager deb --target $(BUILDDIR)

	@sed \
		-e 's/^version:.*$$/version: $(VERSION)/'  \
		-e 's/^arch:.*$$/arch: $(ARCH)/' \
		-e 's/^platform:.*$$/platform: $(PLATFORM)/' \
		etc/nfpm/rpc.yaml > $(BUILDDIR)/rpc.yaml
	@nfpm pkg -f $(BUILDDIR)/rpc.yaml --packager deb --target $(BUILDDIR)

	@sed \
		-e 's/^version:.*$$/version: $(VERSION)/'  \
		-e 's/^arch:.*$$/arch: $(ARCH)/' \
		-e 's/^platform:.*$$/platform: $(PLATFORM)/' \
		etc/nfpm/googlecast.yaml > $(BUILDDIR)/googlecast.yaml
	@nfpm pkg -f $(BUILDDIR)/googlecast.yaml --packager deb --target $(BUILDDIR)

	@echo
	@ls -1 $(BUILDDIR)/*.deb
	@echo
	@echo "Use sudo dpkg -i <package> to install"
	@echo

release: debian pkgcloud
	@$(foreach file, $(wildcard $(BUILDDIR)/*.deb), pkgcloud-push $(PACKAGECLOUD_REPO) $(file);)

	@echo
	@echo "Use the following command in order to add the gopi repository to your list of"
	@echo "repos:"
	@echo
	@echo "  curl -s https://packagecloud.io/install/repositories/djthorpe/gopi/script.deb.sh | sudo bash"
	@echo
	@echo "Then use sudo apt install <package> to install"
	@echo
	
# Commands
helloworld: builddir
	PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" $(GO) build -o ${BUILDDIR}/helloworld -tags "$(TAGS)" ${GOFLAGS} ./cmd/helloworld

httpserver: builddir
	PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" $(GO) build -o ${BUILDDIR}/httpserver -tags "$(TAGS)" ${GOFLAGS} ./cmd/httpserver

hw: rpi darwin freetype
	PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" $(GO) build -o ${BUILDDIR}/hw -tags "$(TAGS)" ${GOFLAGS} ./cmd/hw

argonone: builddir protogen rpi
	PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" $(GO) build -o ${BUILDDIR}/argonone -tags "$(TAGS)" ${GOFLAGS} ./cmd/argonone

douglas: builddir rpi
	PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" $(GO) build -o ${BUILDDIR}/douglas -tags "$(TAGS)" ${GOFLAGS} ./cmd/douglas

dnsregister: builddir
	PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" $(GO) build -o ${BUILDDIR}/dnsregister -tags "$(TAGS)" ${GOFLAGS} ./cmd/dnsregister

rpc: builddir protogen
	PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" $(GO) build -o ${BUILDDIR}/rpc -tags "$(TAGS)" ${GOFLAGS} ./cmd/rpc

googlecast: builddir protogen
	PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" $(GO) build -o ${BUILDDIR}/googlecast -tags "$(TAGS)" ${GOFLAGS} ./cmd/googlecast

# In testing
mediakit: builddir ffmpeg chromaprint
	PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" $(GO) build -o ${BUILDDIR}/mediakit -tags "$(TAGS)" ${GOFLAGS} ./cmd/mediakit

dvbkit: builddir dvb
	PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" $(GO) build -o ${BUILDDIR}/dvbkit -tags "$(TAGS)" ${GOFLAGS} ./cmd/dvbkit

gx: builddir rpi egl drm gbm
	PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" $(GO) build -o ${BUILDDIR}/gx -tags "$(TAGS)" ${GOFLAGS} ./cmd/gx

dx: builddir rpi egl dispmanx
	PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" $(GO) build -o ${BUILDDIR}/dx -tags "$(TAGS)" ${GOFLAGS} ./cmd/dx

mmaltest: rpi mmal
	PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" $(GO) test -v -tags "$(TAGS)" ./pkg/sys/mmal

mmaldecoder: rpi mmal
	PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" $(GO) build -o ${BUILDDIR}/mmalreader -tags "$(TAGS)" ${GOFLAGS} ./cmd/mmalreader

mmalplayer: rpi mmal
	PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" $(GO) build -o ${BUILDDIR}/mmalplayer -tags "$(TAGS)" ${GOFLAGS} ./cmd/mmalplayer

openvgtest: rpi dispmanx openvg egl
	PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" $(GO) test -v -tags "$(TAGS)" ${GOFLAGS} ./pkg/sys/openvg

dxtest: dispmanx egl rpi
	PKG_CONFIG_PATH="$(PKG_CONFIG_PATH)" $(GO) test -v -tags "$(TAGS)" ${GOFLAGS} ./pkg/graphics/surface/dispmanx

# Build rules - dependencies
nfpm:
	$(GO) get github.com/goreleaser/nfpm/cmd/nfpm@v1.10.1

pkgcloud:
	$(GO) get github.com/mlafeldt/pkgcloud/cmd/pkgcloud-push

protogen: protoc
	$(GO) get google.golang.org/protobuf/cmd/protoc-gen-go
	$(GO) get google.golang.org/grpc/cmd/protoc-gen-go-grpc
	$(GO) generate -x ./rpc

protoc:
ifeq ($(shell which protoc),)
	@echo apt install protobuf-compiler
	$(error protoc is not installed)
endif
