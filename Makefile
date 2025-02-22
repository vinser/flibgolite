# Many mickle makes a Makefile ;) 

# variable definitions

CMD_MAIN := $(shell find cmd/ -name main.go)
APP := $(patsubst cmd/%/main.go,%,$(CMD_MAIN))

VERSION := $(shell git describe --tags --always --dirty)
GOVERSION := $(shell go env GOVERSION)
BUILDTIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

HOSTOS := $(shell go env GOHOSTOS)
HOSTARCH := $(shell go env GOHOSTARCH)

build_cmd = \
	CGO_ENABLED=0 \
	GOOS=$(1) \
	GOARCH=$(2) \
	$(if $(3),GOARM=$(3)) \
	go build -ldflags "-s \
	-X 'main.version=$(VERSION)' \
    -X 'main.buildTime=$(BUILDTIME)' \
    -X 'main.goversion=$(GOVERSION)' \
	-X 'main.target=$(1)-$(2)$(if $(3),-$(3))'" \
	-o bin/$(APP)-$(1)-$(2)$(if $(3),-$(3))$(if $(findstring windows,$(1)),.exe) \
	$(abspath $(dir $(CMD_MAIN)))

all: build

# Current host build
build:
	$(call build_cmd,$(HOSTOS),$(HOSTARCH),)


# Cross builds
xbuild: linux darwin windows freebsd

# Linux builds ========================
linux: build_linux_amd64 build_linux_armV6 build_linux_armV7 build_linux_arm64

build_linux_amd64:
	$(call build_cmd,linux,amd64,)

build_linux_armV6:
	$(call build_cmd,linux,arm,6)

build_linux_armV7:
	$(call build_cmd,linux,arm,7)

build_linux_arm64:
	$(call build_cmd,linux,arm64,)

# OSX (Darwin) builds =================
darwin: build_darwin_amd64 build_darwin_arm64

build_darwin_amd64:
	$(call build_cmd,darwin,amd64,)

build_darwin_arm64:
	$(call build_cmd,darwin,arm64,)

# Windows builds ======================
windows: build_windows_386 build_windows_amd64

build_windows_386:
	$(call build_cmd,windows,386,)

build_windows_amd64:
	$(call build_cmd,windows,amd64,)

# FreeBSD builds ========================
freebsd: build_freebsd_amd64 build_freebsd_armV6 build_freebsd_armV7 build_freebsd_arm64

build_freebsd_amd64:
	$(call build_cmd,freebsd,amd64,)

build_freebsd_armV6:
	$(call build_cmd,freebsd,arm,6)

build_freebsd_armV7:
	$(call build_cmd,freebsd,arm,7)

build_freebsd_arm64:
	$(call build_cmd,freebsd,arm64,)

docker_xbuild:
	docker build --platform linux/amd64,linux/arm64,linux/arm/v7,linux/arm/v6 --tag vinser/$(APP):$(VERSION) .
	docker image tag vinser/$(APP):$(VERSION) vinser/$(APP):latest
	
docker_push:
	docker push vinser/$(APP):$(VERSION)

run_container:
	$(shell ./docker_run.sh)

.PHONY: all build xbuild 
.PHONY: linux darwin windows 
.PHONY: build_linux_arm build_linux_arm64 build_linux_amd64 
.PHONY: build_darwin_amd64 build_darwin_arm64 
.PHONY: build_windows_386 build_windows_amd64
.PHONY: build_freebsd_amd64 build_freebsd_armV6 build_freebsd_armV7 build_freebsd_arm64
.PHONY: docker_xbuild docker_push run_container