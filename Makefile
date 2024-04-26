# Many mickle makes a Makefile ;) 

# variable definitions

APP := flibgolite

VERSION := $(shell git describe --tags --always --dirty)
GOVERSION := $(shell go env GOVERSION)
BUILDTIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

HOSTOS := $(shell go env GOHOSTOS)
HOSTARCH := $(shell go env GOHOSTARCH)

CMD_MAIN := $(shell find cmd/$(APP) -name main.go)
OUTPUT := $(patsubst cmd/%/main.go,%,$(CMD_MAIN))

build_cmd = \
	GOOS=$(1) \
	GOARCH=$(2) \
	$(if $(3),GOARM=$(3)) \
	go build -ldflags " \
	-X 'main.version=$(VERSION)' \
    -X 'main.buildTime=$(BUILDTIME)' \
    -X 'main.goversion=$(GOVERSION)' \
	-X 'main.target=$(1)-$(2)$(if $(3),-$(3))'" \
	-o $(OUTPUT)-$(1)-$(2)$(if $(3),-$(3))$(if $(findstring windows,$(1)),.exe) \
	$(abspath $(dir $(CMD_MAIN)))

all: build

# Current host build
build:
	$(call build_cmd,$(HOSTOS),$(HOSTARCH),)


# Cross builds
xbuild: linux darwin windows

# Linux builds ========================
linux: build_linux_arm build_linux_arm64 build_linux_amd64

build_linux_amd64:
	$(call build_cmd,linux,amd64,)

build_linux_arm:
	$(call build_cmd,linux,arm,6)

build_linux_arm64:
	$(call build_cmd,linux,arm64,)

# OSX (Darwin) builds =================
darwin: build_darwin_amd64

build_darwin_amd64:
	$(call build_cmd,darwin,amd64,)

build_darwin_arm64:
	$(call build_cmd,darwin,arm64,)

# Windows builds ======================
windows:  build_windows_amd64

build_windows_amd64:
	$(call build_cmd,windows,amd64,)

.PHONY: all check build xbuild linux darwin windows build_linux_arm build_linux_arm64 build_linux_amd64 build_darwin_amd64 build_darwin_arm64 build_windows_amd64