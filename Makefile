# CleanPDF Makefile
# Cross-platform build targets for Go/Fyne application

APP_NAME = CleanPDF
BINARY_NAME = cleanpdf
VERSION ?= 2.0.0
BUILD_DIR = build
MODULE = github.com/james-see/cleanpdfapp

# Go build flags
LDFLAGS = -s -w -X main.Version=$(VERSION)

.PHONY: all clean build test run install deps

all: build

# Install dependencies
deps:
	go mod download
	go install fyne.io/fyne/v2/cmd/fyne@latest

# Build for current platform
build:
	go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) .

# Run the application
run:
	go run .

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME)

# === Cross-Platform Builds ===

# Build for macOS (Universal Binary - arm64 + amd64)
build-macos: $(BUILD_DIR)
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	lipo -create -output $(BUILD_DIR)/$(BINARY_NAME)-darwin-universal $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64

# Build for Windows
build-windows: $(BUILD_DIR)
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .

# Build for Linux
build-linux: $(BUILD_DIR)
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .

# === macOS App Bundle ===

# Package macOS app bundle using fyne
package-macos: deps $(BUILD_DIR)
	fyne package -os darwin -icon Icon.png -name "$(APP_NAME)" -appID us.jamescampbell.cleanpdf -release
	mv "$(APP_NAME).app" $(BUILD_DIR)/

# Package macOS app for release (creates zip)
package-macos-release: package-macos
	cd $(BUILD_DIR) && zip -r $(APP_NAME)-macos-v$(VERSION).zip "$(APP_NAME).app"

# === Windows Package ===

package-windows: deps $(BUILD_DIR)
	fyne package -os windows -icon Icon.png -name "$(APP_NAME)" -appID us.jamescampbell.cleanpdf
	mv "$(APP_NAME).exe" $(BUILD_DIR)/

package-windows-release: package-windows
	cd $(BUILD_DIR) && zip -r $(APP_NAME)-windows-v$(VERSION).zip "$(APP_NAME).exe"

# === Linux Package ===

package-linux: deps $(BUILD_DIR)
	fyne package -os linux -icon Icon.png -name "$(APP_NAME)" -appID us.jamescampbell.cleanpdf
	mv $(APP_NAME).tar.xz $(BUILD_DIR)/ 2>/dev/null || mv $(BINARY_NAME) $(BUILD_DIR)/$(BINARY_NAME)-linux

package-linux-release: build-linux
	cd $(BUILD_DIR) && tar -czvf $(APP_NAME)-linux-v$(VERSION).tar.gz $(BINARY_NAME)-linux-amd64

# === Full Release ===

# Build all platforms for release
release: clean $(BUILD_DIR) package-macos-release
	@echo "Release artifacts created in $(BUILD_DIR)/"
	@ls -la $(BUILD_DIR)/

# Create build directory
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# Install locally (macOS)
install: package-macos
	cp -r $(BUILD_DIR)/"$(APP_NAME).app" /Applications/

# Uninstall (macOS)
uninstall:
	rm -rf /Applications/"$(APP_NAME).app"

# Development helpers
fmt:
	go fmt ./...

lint:
	golangci-lint run

.PHONY: build-macos build-windows build-linux package-macos package-windows package-linux
.PHONY: package-macos-release package-windows-release package-linux-release
.PHONY: release install uninstall fmt lint

