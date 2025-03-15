# Contributing to VidKit

Thank you for your interest in contributing to VidKit! This document provides guidelines and instructions for contributing to the project.

## Development Setup

1. Fork and clone the repository
2. Install dependencies:
   - Go 1.21+
   - FFmpeg (for ffprobe)
3. Set up your environment:
   ```bash
   go mod download
   ```

## Code Standards

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Run `golangci-lint run` before submitting code
- Keep functions small and focused
- Write comprehensive documentation for public functions
- Follow the existing code structure

## Testing Guidelines

VidKit uses Go's built-in testing framework. All new features should include appropriate tests.

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with race condition detection
go test -race ./...

# Run tests with coverage report
go test -cover ./...

# Run tests for a specific package
go test ./internal/pkg/metadata
```

### Types of Tests

1. **Unit Tests**
   - Located in files named `*_test.go` alongside the code they test
   - Focus on testing individual functions or components in isolation
   - Use mocks for external dependencies (like API clients)

2. **Integration Tests**
   - Test interaction between components
   - May require external dependencies (e.g., FFmpeg)

3. **End-to-End Tests**
   - Test complete workflows with real files
   - Use the test data in `test_videos/` directory

### Testing TV Show Metadata

The TV show metadata functionality is particularly complex and requires careful testing:

1. **Setup Test Data**
   ```bash
   # Generate test video files with different naming patterns
   ./tools/generate_test_videos.sh
   ```

2. **Test Different Filename Formats**
   - Test S01E01 format: `Breaking.Bad.S01E01.Pilot.mp4`
   - Test 1x01 format: `Game.of.Thrones.1x01.mp4`
   - Test descriptive format: `Stranger.Things.Season.1.Episode.1.mp4`
   - Test with year: `The Office (2005) S01E01.mp4`
   - Test with quality info: `Friends.S01E01.DVDRip.x264.mp4`

3. **Test API Integration**
   - The tests use mock HTTP servers to test API interaction
   - Ensure real-world testing with the actual TvMaze API
   - Test error handling for network issues or missing data

4. **Test Output Formats**
   ```bash
   # Test with custom TV show format
   go run cmd/vidkit/main.go --preview --tv-format "{title}.S{season:02d}E{episode:02d}.{episode_title}" test_videos/*.mp4
   
   # Test with scene-style naming (dots instead of spaces)
   go run cmd/vidkit/main.go --preview --scene test_videos/*.mp4
   ```

### Testing Year Extraction Rules

VidKit now has specific rules for extracting years from filenames. When testing or adding features, make sure to test:

1. **Year Extraction Test Cases**
   - Years in parentheses: `Movie Title (2023).mp4` (should extract 2023)
   - Years in square brackets: `Movie Title [2023].mp4` (should extract 2023)
   - Years without delimiters: `Movie.Title.2023.mp4` (should NOT extract year)
   - Multiple potential years: `Movie (2023) [2024].mp4` (should extract first valid year)

2. **Title Preservation**
   - When no delimited year is found, the original format should be preserved
   - Example: `The.Matrix.1999.mp4` should maintain dot format

3. **Metadata Search Impact**
   - Test how year extraction affects metadata search accuracy
   - Compare search results with and without proper year delimiting
   - Verify the handling of titles containing numeric sequences

These tests are implemented in the `metadata` package:
- `TestExtractMovieInfo` in `tmdb_test.go`
- `TestExtractTVShowInfo` in `tvmaze_test.go`

### Test Mocks

The test files use mock implementations for external dependencies:

- `tmdb_test.go` includes a mock TMDb client
- `tvmaze_test.go` includes a mock HTTP server for TvMaze responses

When adding new features, follow this pattern to create testable code.

### Integration Tests

VidKit includes several integration tests to verify proper behavior with real-world examples:

- **Basic Integration Test**: Tests file detection and basic operations
  ```bash
  make integration-test
  ```

- **TV Show Format Tests**: Tests various TV show filename formats
  ```bash
  make test-tv-formats
  ```

- **Custom Format Tests**: Tests custom formatting options
  ```bash
  make test-custom-format
  ```

- **Directory Organization Tests**: Tests metadata-based directory organization
  ```bash
  make test-organization
  ```

- **CI Integration Tests**: A comprehensive test suite designed for CI environments
  ```bash
  make ci-integration-test
  ```

### GitHub Actions

The project includes multiple GitHub Actions workflows:

1. **Test and Lint**: Runs unit tests and code linting
2. **Integration Tests**: Runs integration tests to verify functionality with real-world examples
3. **Release**: Builds and publishes new releases

When submitting a pull request, make sure all tests pass with:

```bash
make test-all
```

### Test Videos

The test videos are automatically generated using FFmpeg. You can regenerate them with:

```bash
make generate-test-videos
```

No actual video content is included in the repository - the test files are small, valid MP4 containers with test patterns.

## Pull Request Process

1. Create a new branch for your feature or fix
2. Add tests for your changes
3. Ensure all tests pass
4. Update documentation if needed
5. Submit a pull request with a clear description of the changes

## Release Process

VidKit uses [GoReleaser](https://goreleaser.com/) for building and publishing releases:

1. Ensure all tests pass
2. Update version in relevant files
3. Create and push a git tag with the version
4. GoReleaser will automatically build and publish the release

## Getting Help

If you need help with contributing, please open an issue or contact the maintainers directly.

Thank you for contributing to VidKit!
