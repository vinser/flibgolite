# Makefile for Go cross-builds with Windows 7 legacy support
# Many mickle makes a Makefile ;)

# ==============================================================================
# Go versions
# ==============================================================================
GO_STD_VER := 1.23.2
GO_LEGACY_VER := 1.25.3-1

# ==============================================================================
# Variables
# ==============================================================================

CMD_MAIN := $(shell find cmd/ -name main.go)
APP := $(patsubst cmd/%/main.go,%,$(CMD_MAIN))

VERSION := $(shell git describe --tags --always --dirty)
BUILDTIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Control verbosity: make v=1 ...
ifdef v
  Q :=
else
  Q := @
endif

# === Argument handling (allows "make <target> ./cmd/myapp") ===
ifeq (,$(filter help detect docker_xbuild docker_push run_container,$(MAKECMDGOALS)))
  TARGET := $(word 2,$(MAKECMDGOALS))
  ifneq ($(TARGET),)
    # Prevent make from treating TARGET as a separate goal
    $(eval $(TARGET):;@:)
    TARGET_NAME := $(notdir $(TARGET))
  else
    TARGET := ./cmd/$(APP)
    TARGET_NAME := $(APP)
  endif
endif

# --- Check standard Go ---
GO_STD := $(shell command -v go 2>/dev/null)
# --- Check go-legacy-win7 ---
GO_LEGACY := $(HOME)/opt/go-legacy-win7/bin/go

define select_go_bin
$(if $(findstring win7,$(1)), \
  $(if $(wildcard $(GO_LEGACY)), $(GO_LEGACY), \
    $(error go-legacy-win7 not found! Run "make install-legacy" to install it.)), \
  $(if $(GO_STD),$(GO_STD), \
    $(error go (standard) not found! Install from https://go.dev/dl/)))
endef

# ==============================================================================
# Build macro
# ==============================================================================

define build_cmd
	@echo "───────────────────────────────────────────────"
	@echo "🔧 Building for: $(1)-$(2)$(if $(3),-$(3))"
	$(eval GO_BIN := $(call select_go_bin,$(1)))
	$(eval GOVERSION := $(shell $(GO_BIN) env GOVERSION))
	$(Q)CGO_ENABLED=0 \
	GOOS=$(word 1, $(subst -, , $(1)))\
	GOARCH=$(2) \
	$(if $(3),GOARM=$(3)) \
	$(GO_BIN) build -ldflags "-s \
		-X 'main.version=$(VERSION)' \
		-X 'main.buildTime=$(BUILDTIME)' \
		-X 'main.goversion=$(GOVERSION)' \
		-X 'main.target=$(1)-$(2)$(if $(3),-$(3))'" \
		-o bin/$(TARGET_NAME)-$(1)-$(2)$(if $(3),-$(3))$(if $(findstring windows,$(1)),.exe) \
		$(TARGET)
	@echo "✅ Built with: $(GOVERSION)"
	@echo "📦 Output: bin/$(TARGET_NAME)-$(1)-$(2)$(if $(3),-$(3))$(if $(findstring windows,$(1)),.exe)"
endef

# ==============================================================================
# Detect tools
# ==============================================================================
# --- Detect host OS and architecture ---
GOHOSTOS  := $(shell go env GOHOSTOS 2>/dev/null || uname | tr A-Z a-z)
GOHOSTARCH := $(shell go env GOHOSTARCH 2>/dev/null || uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/;s/armv7l/arm/')


detect:
	@echo "Detecting Go toolchains..."
	@echo "-------------------------------------------------------------------"
	@if [ -n "$(GO_STD)" ]; then \
		echo "Standard Go:      ✅ $$($(GO_STD) version | awk '{print $$3, $$4}')"; \
	else \
		echo "Standard Go:      ❌ not found"; \
		echo "  → Install from https://go.dev/dl/"; \
	fi
	@if [ -x "$(GO_LEGACY)" ]; then \
		echo "Legacy Go (Win7): ✅ $$($(GO_LEGACY) version | awk '{print $$3, $$4}')"; \
	else \
		echo "Legacy Go (Win7): ❌ not found"; \
		echo "  → Run: make install-legacy"; \
	fi
	@echo "-------------------------------------------------------------------"
	@$(MAKE) --no-print-directory help


# --- Install go-legacy-win7 ---
install-legacy:
	@echo "Installing go-legacy-win7 $(GO_LEGACY_VER)..."
	@mkdir -p $(HOME)/opt
	@cd $(HOME)/opt && \
	curl -L -o go-legacy.tar.gz https://github.com/thongtech/go-legacy-win7/releases/download/v$(GO_LEGACY_VER)/go-legacy-win7-$(GO_LEGACY_VER).$(GOHOSTOS)_$(GOHOSTARCH).tar.gz && \
	tar -xzf go-legacy.tar.gz && rm go-legacy.tar.gz && \
	ln -sf $(HOME)/opt/go-legacy-win7/bin/go $(HOME)/go/bin/go-legacy-win7
	@echo "✅ Installed to $(HOME)/opt/go-legacy-win7"
	@echo "👉 To use it in Makefile or scripts run go-legacy-win7 with options"
	@echo "   Ex. go-legacy-win7 version"

# ==============================================================================
# Build targets
# ==============================================================================

all: build

# Build for current host
build:
	$(call build_cmd,$(GOHOSTOS),$(GOHOSTARCH),)

# Cross builds
xbuild: linux darwin windows freebsd

# ---------------- Linux ----------------
linux: linux_amd64 linux_armv6 linux_armv7 linux_arm64 linux_386
linux_amd64:      ; $(call build_cmd,linux,amd64,)
linux_armv6:      ; $(call build_cmd,linux,arm,6)
linux_armv7:      ; $(call build_cmd,linux,arm,7)
linux_arm64:      ; $(call build_cmd,linux,arm64,)
linux_386:        ; $(call build_cmd,linux,386,)

# ---------------- macOS ----------------
darwin: darwin_amd64 darwin_arm64
darwin_amd64:     ; $(call build_cmd,darwin,amd64,)
darwin_arm64:     ; $(call build_cmd,darwin,arm64,)

# ---------------- Windows ----------------
windows: windows_amd64 windows_arm64 windows_386 windows_amd64_win7 windows_386_win7
windows_amd64:    ; $(call build_cmd,windows,amd64,)
windows_arm64:    ; $(call build_cmd,windows,arm64,)
windows_386:      ; $(call build_cmd,windows,386,)
windows_amd64_win7: ; $(call build_cmd,windows-win7,amd64,)
windows_386_win7: ; $(call build_cmd,windows-win7,386,)

# ---------------- FreeBSD ----------------
freebsd: freebsd_amd64 freebsd_armv6 freebsd_armv7 freebsd_arm64 freebsd_386
freebsd_amd64:    ; $(call build_cmd,freebsd,amd64,)
freebsd_armv6:    ; $(call build_cmd,freebsd,arm,6)
freebsd_armv7:    ; $(call build_cmd,freebsd,arm,7)
freebsd_arm64:    ; $(call build_cmd,freebsd,arm64,)
freebsd_386:      ; $(call build_cmd,freebsd,386,)

# ==============================================================================
# Docker helper
# ==============================================================================
docker_xbuild:
	docker build --platform linux/amd64,linux/arm64,linux/arm/v7,linux/arm/v6 --tag vinser/$(APP):$(VERSION) .
	docker image tag vinser/$(APP):$(VERSION) vinser/$(APP):latest

docker_push:
	docker push vinser/$(APP):$(VERSION)

run_container:
	$(shell ./docker_run.sh)

# ==============================================================================
# Help
# ==============================================================================
help:$(1) 
	@echo "🛟  Help to make"
	@echo "-------------------------------------------------------------------"
	@echo "Available targets:"
	@echo "  make build [path]                Build for current host"
	@echo "  make linux [path]                Build all Linux targets"
	@echo "  make windows [path]              Build all Windows targets"
	@echo "  make windows_amd64_win7 [path]   Build with go-legacy-win7"
	@echo "  make detect                      Show detected toolchains"
	@echo "  make xbuild [path]               Build all OS/arch combinations"
	@echo ""
	@echo "Examples:"
	@echo "  make build ./cmd/myapp"
	@echo "  make windows_amd64 ./cmd/myapp"
	@echo "  make windows_386_win7 ./cmd/myapp"
	@echo ""

# ==============================================================================
# PHONY
# ==============================================================================
.PHONY: all build xbuild detect help \
        linux linux_amd64 linux_armv6 linux_armv7 linux_arm64 linux_386 \
        darwin darwin_amd64 darwin_arm64 \
        windows windows_amd64 windows_arm64 windows_386 windows_amd64_win7 windows_386_win7 \
        freebsd freebsd_amd64 freebsd_armv6 freebsd_armv7 freebsd_arm64 freebsd_386 \
        docker_xbuild docker_push run_container
