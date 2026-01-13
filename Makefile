# GitSync Makefile
BINARY_NAME=gitsync
VERSION=1.0.0
BUILD_DIR=build

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-s -w -X github.com/princetheprogrammerbtw/gitsynq/cmd.version=$(VERSION)"

.PHONY: all build clean test install

all: clean build

build:
	@echo "🔨 Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

build-all: clean
	@echo "🔨 Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	
	# Linux AMD64
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	
	# Linux ARM64
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .
	
	# macOS AMD64
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	
	# macOS ARM64 (M1/M2)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	
	# Windows
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	
	@echo "✅ All builds complete!"

install: build
	@echo "📦 Installing to ~/bin..."
	@mkdir -p ~/bin
	@cp $(BUILD_DIR)/$(BINARY_NAME) ~/bin/
	@echo "✅ Installed! Make sure ~/bin is in your PATH"

clean:
	@echo "🧹 Cleaning..."
	@rm -rf $(BUILD_DIR)
	$(GOCLEAN)

test:
	@echo "🧪 Running tests..."
	$(GOTEST) -v ./...

deps:
	@echo "📦 Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

run:
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) . && ./$(BUILD_DIR)/$(BINARY_NAME)
