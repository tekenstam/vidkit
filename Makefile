.PHONY: build test test-verbose test-race test-cover lint lint-vet lint-fmt clean run generate-test-videos integration-test

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOVET=$(GOCMD) vet
GOFMT=$(GOCMD) fmt
BINARY_NAME=vidkit
MAIN_PATH=./cmd/vidkit

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)

test:
	$(GOTEST) ./...

test-verbose:
	$(GOTEST) -v ./...

test-race:
	$(GOTEST) -race ./...

test-cover:
	$(GOTEST) -cover ./...

# Basic lint with Go's built-in tools
lint: lint-vet lint-fmt

# Run go vet
lint-vet:
	@echo "Running go vet..."
	$(GOVET) ./...

# Run go fmt
lint-fmt:
	@echo "Running go fmt..."
	$(GOFMT) ./...
	@echo "Checking for formatting issues..."
	@gofmt -l . | grep ".*\.go" > /dev/null && echo "Some files need formatting. Run 'go fmt ./...'" && exit 1 || echo "All files properly formatted."

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf test_videos/*.mp4

run:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)
	./$(BINARY_NAME) $(ARGS)

deps:
	$(GOMOD) download

# Generate test videos
generate-test-videos:
	./tools/generate_test_videos.sh

# Run integration tests with test videos
integration-test: build generate-test-videos
	@echo "=== Running integration tests ==="
	./$(BINARY_NAME) --preview --batch test_videos/*.mp4
	@echo "Integration tests completed successfully"

# CI-specific integration tests (more stringent, fails on errors)
ci-integration-test: build
	@echo "=== Running CI integration tests ==="
	@echo "Generating test videos with tools/generate_test_videos.sh..."
	./tools/generate_test_videos.sh
	@echo "Running format detection tests..."
	./$(BINARY_NAME) --preview test_videos/Breaking.Bad.S01E01.Pilot.mp4
	./$(BINARY_NAME) --preview test_videos/The.Shawshank.Redemption.1994.1080p.BluRay.x264.mp4
	./$(BINARY_NAME) --preview test_videos/Game.of.Thrones.1x01.mp4
	./$(BINARY_NAME) --preview "test_videos/The Good Place (2016) - Season 1 Episode 1 - Everything Is Fine.mp4"
	@echo "Running organization tests..."
	./tools/test_directory_organization.sh
	@echo "CI integration tests completed successfully"

# Test TV show functionality
test-tv:
	$(GOTEST) ./internal/pkg/metadata -run TestExtractTVShowInfo
	$(GOTEST) ./internal/pkg/metadata -run TestTvMazeProvider

# Test different TV formats (requires generating test videos first)
test-tv-formats: generate-test-videos
	./$(BINARY_NAME) --preview test_videos/Breaking.Bad.S01E01.Pilot.mp4
	./$(BINARY_NAME) --preview test_videos/Game.of.Thrones.1x01.mp4
	./$(BINARY_NAME) --preview "test_videos/The Good Place (2016) - Season 1 Episode 1 - Everything Is Fine.mp4"

# Test custom formatting
test-custom-format: generate-test-videos
	./$(BINARY_NAME) --preview --tv-format "{title}.S{season:02d}E{episode:02d}.{episode_title}" test_videos/Breaking.Bad.S01E01.Pilot.mp4
	./$(BINARY_NAME) --preview --scene test_videos/Breaking.Bad.S01E01.Pilot.mp4

# Test directory organization
test-organization: build generate-test-videos
	@echo "Testing directory organization features..."
	./tools/test_directory_organization.sh

# Test resolution detection
test-resolution: build
	@echo "=== Testing resolution and video quality detection ==="
	./tools/test_resolution.sh

# Test error handling
test-error-handling: build
	@echo "=== Testing error handling ==="
	./tools/test_error_handling.sh

# Test command-line interface
test-cli: build
	@echo "=== Testing command-line interface ==="
	./tools/test_cli.sh

# Run all tests
test-all: test test-race generate-test-videos test-tv-formats test-custom-format test-organization test-resolution test-error-handling test-cli

# Run code quality checks - both tests and linting
quality: test lint

help:
	@echo "Available commands:"
	@echo "  make build               - Build the vidkit binary"
	@echo "  make test                - Run all tests"
	@echo "  make test-verbose        - Run tests with verbose output"
	@echo "  make test-race           - Run tests with race detection"
	@echo "  make test-cover          - Run tests with coverage report"
	@echo "  make lint                - Run basic linting (vet + fmt)"
	@echo "  make lint-vet            - Run go vet static analysis"
	@echo "  make lint-fmt            - Check code formatting with go fmt"
	@echo "  make quality             - Run both tests and linting"
	@echo "  make run ARGS=\"file.mp4\" - Build and run with arguments"
	@echo "  make clean               - Clean up build artifacts"
	@echo "  make deps                - Download dependencies"
	@echo "  make generate-test-videos - Generate test video files"
	@echo "  make integration-test    - Run integration tests with test videos"
	@echo "  make ci-integration-test - Run CI-specific integration tests (stricter)"
	@echo "  make test-tv             - Run TV show-specific tests"
	@echo "  make test-tv-formats     - Test different TV show naming formats"
	@echo "  make test-custom-format  - Test custom naming formats"
	@echo "  make test-organization   - Test directory organization features"
	@echo "  make test-resolution     - Test resolution and video quality detection"
	@echo "  make test-error-handling - Test error handling capabilities"
	@echo "  make test-cli            - Test command-line interface functionality"
	@echo "  make test-all            - Run all tests including integration tests"
	@echo "  make help                - Show this help message"
