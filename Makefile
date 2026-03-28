# HolyC Compiler - Makefile
# Build, test, and package for various Linux distributions

NAME := holyc
VERSION := 0.1.0
GO := go
GIT := git

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[1;33m
NC := \033[0m # No Color

.PHONY: all build clean test install uninstall help
.PHONY: pkg-arch pkg-deb pkg-rpm pkg-all
.PHONY: release release-all

# Default target
all: build

## Build
build:
	@echo "$(GREEN)Building $(NAME)...$(NC)"
	$(GO) build -o $(NAME) -v .
	@echo "$(GREEN)Build complete: $(NAME)$(NC)"

## Build with debug symbols
build-debug:
	@echo "$(YELLOW)Building with debug symbols...$(NC)"
	$(GO) build -gcflags="all=-N -l" -o $(NAME) -v .

## Build static binary
build-static:
	@echo "$(YELLOW)Building static binary...$(NC)"
	CGO_ENABLED=0 $(GO) build -a -installsuffix cgo -o $(NAME) -v .

## Clean build artifacts
clean:
	@echo "$(YELLOW)Cleaning...$(NC)"
	rm -f $(NAME) $(NAME).exe
	rm -rf dist/
	rm -rf pkg/
	rm -rf *.bin
	rm -rf debian/holyc
	rm -rf rpm-build/
	@echo "$(GREEN)Clean complete$(NC)"

## Run tests
test:
	@echo "$(GREEN)Running tests...$(NC)"
	$(GO) test -v ./...

## Run tests with coverage
test-coverage:
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report: coverage.html$(NC)"

## Install to /usr/local/bin
install: build
	@echo "$(GREEN)Installing $(NAME) to /usr/local/bin...$(NC)"
	sudo cp $(NAME) /usr/local/bin/$(NAME)
	sudo chmod +x /usr/local/bin/$(NAME)
	@echo "$(GREEN)Installation complete$(NC)"

## Uninstall from /usr/local/bin
uninstall:
	@echo "$(YELLOW)Uninstalling $(NAME)...$(NC)"
	sudo rm -f /usr/local/bin/$(NAME)
	@echo "$(GREEN)Uninstallation complete$(NC)"

## Build Arch Linux package (PKGBUILD)
pkg-arch: clean
	@echo "$(GREEN)Building Arch Linux package...$(NC)"
	mkdir -p pkg
	cd $(CURDIR) && makepkg -f --cleanbuild
	@echo "$(GREEN)Arch package complete$(NC)"

## Build Debian/Ubuntu package
pkg-deb: clean
	@echo "$(GREEN)Building Debian package...$(NC)"
	mkdir -p pkg
	cp -r debian/ debian-build/
	cd $(CURDIR) && debuild -us -uc
	rm -rf debian-build/
	@echo "$(GREEN)Debian package complete$(NC)"

## Build RPM package (Fedora/RHEL/openSUSE)
pkg-rpm: clean
	@echo "$(GREEN)Building RPM package...$(NC)"
	mkdir -p rpm-build/{BUILD,RPMS,SOURCES,SPECS,SRPMS}
	tar -czf rpm-build/SOURCES/$(NAME)-$(VERSION).tar.gz \
		--transform 's,^,$(NAME)-$(VERSION)/,' \
		--exclude='.git' \
		--exclude='pkg' \
		--exclude='dist' \
		--exclude='*.bin' \
		--exclude='debian' \
		--exclude='rpm' \
		--exclude='Makefile' \
		.
	rpmbuild --define "_topdir $(CURDIR)/rpm-build" \
		-bb rpm/$(NAME).spec
	@echo "$(GREEN)RPM package complete$(NC)"

## Build all packages
pkg-all: pkg-arch pkg-deb pkg-rpm
	@echo "$(GREEN)All packages built successfully$(NC)"

## Create a release (git tag + build)
release: build
	@echo "$(GREEN)Creating release $(VERSION)...$(NC)"
	$(GIT) tag -a v$(VERSION) -m "Release version $(VERSION)"
	@echo "$(YELLOW)Push tag with: git push origin v$(VERSION)$(NC)"

## Build release binaries for multiple platforms
release-all: clean
	@echo "$(GREEN)Building multi-platform release...$(NC)"
	mkdir -p dist
	
	# Linux amd64
	GOOS=linux GOARCH=amd64 $(GO) build -o dist/$(NAME)-linux-amd64 -v .
	
	# Linux arm64
	GOOS=linux GOARCH=arm64 $(GO) build -o dist/$(NAME)-linux-arm64 -v .
	
	# macOS amd64
	GOOS=darwin GOARCH=amd64 $(GO) build -o dist/$(NAME)-darwin-amd64 -v .
	
	# macOS arm64
	GOOS=darwin GOARCH=arm64 $(GO) build -o dist/$(NAME)-darwin-arm64 -v .
	
	# Windows amd64
	GOOS=windows GOARCH=amd64 $(GO) build -o dist/$(NAME)-windows-amd64.exe -v .
	
	@echo "$(GREEN)Release builds complete in dist/$(NC)"
	ls -la dist/

## Show help
help:
	@echo "$(NAME) $(VERSION) - Build System"
	@echo ""
	@echo "$(YELLOW)Usage:$(NC) make [target]"
	@echo ""
	@echo "$(GREEN)Build Targets:$(NC)"
	@echo "  all          Build the compiler (default)"
	@echo "  build        Build the compiler"
	@echo "  build-debug  Build with debug symbols"
	@echo "  build-static Build static binary"
	@echo "  clean        Clean build artifacts"
	@echo ""
	@echo "$(GREEN)Test Targets:$(NC)"
	@echo "  test         Run tests"
	@echo "  test-coverage Run tests with coverage report"
	@echo ""
	@echo "$(GREEN)Install Targets:$(NC)"
	@echo "  install      Install to /usr/local/bin"
	@echo "  uninstall    Remove from /usr/local/bin"
	@echo ""
	@echo "$(GREEN)Package Targets:$(NC)"
	@echo "  pkg-arch     Build Arch Linux package (PKGBUILD)"
	@echo "  pkg-deb      Build Debian/Ubuntu package"
	@echo "  pkg-rpm      Build RPM package (Fedora/RHEL)"
	@echo "  pkg-all      Build all packages"
	@echo ""
	@echo "$(GREEN)Release Targets:$(NC)"
	@echo "  release      Create git release tag"
	@echo "  release-all  Build multi-platform binaries"
	@echo ""
	@echo "$(YELLOW)Examples:$(NC)"
	@echo "  make build           # Build compiler"
	@echo "  make install         # Build and install"
	@echo "  make pkg-arch        # Build Arch package"
	@echo "  make release-all     # Build all platforms"
