# VidKit Testing Guide

This document describes the testing infrastructure for VidKit, a Go-based CLI tool for video file analysis and management with TMDb integration and smart file renaming capabilities.

## Testing Philosophy

VidKit follows a comprehensive testing approach with multiple test types:

1. **Unit Tests**: Test individual functions and components in isolation
2. **Integration Tests**: Test complete workflows and CLI functionality
3. **CI Tests**: Run in GitHub Actions without requiring API keys

## Requirements

To run tests locally, you'll need:

- Go 1.21+
- FFmpeg installed and available in your PATH
- API keys (for complete testing with metadata providers)

## Test Structure

VidKit's test structure is organized as follows:

```
/vidkit
├── Makefile              # Contains test targets
├── tools/                # Integration test scripts
│   ├── generate_test_videos.sh       # Creates test videos
│   ├── test_resolution.sh            # Tests resolution detection
│   ├── test_directory_organization.sh # Tests directory organization
│   ├── test_error_handling.sh        # Tests error handling
│   └── test_cli.sh                   # Tests CLI functionality
└── .github/workflows/    # CI configuration
    ├── integration.yml   # Integration test workflow
    └── test.yml          # Unit test workflow
```

## Running Tests

### All Tests

```bash
make test-all
```

This will run all tests, including unit tests, race condition tests, and integration tests.

### Unit Tests

```bash
make test
```

To run tests with race detection:

```bash
make test-race
```

### Integration Tests

Run specific integration tests:

```bash
make test-resolution       # Test resolution detection
make test-organization     # Test directory organization
make test-error-handling   # Test error handling
make test-cli              # Test CLI functionality
```

### CI Mode (No API Keys)

To run tests without API keys (as they run in CI environments):

```bash
make test-resolution-ci
make test-organization-ci
make test-error-handling-ci
make test-cli-ci
```

Or run all CI tests:

```bash
make ci-integration-test
```

## Test Video Generation

VidKit uses synthetic test videos to ensure consistent testing. Generate them with:

```bash
make generate-test-videos
```

This creates videos with different:
- Resolutions (480p, 720p, 1080p)
- Frame rates (24fps, 30fps, 60fps)
- Codecs (h264, mpeg4)
- Naming patterns (movies, TV shows)

## Writing New Tests

### Adding a New Integration Test

1. Create a new script in the `tools/` directory
2. Follow the pattern in existing test scripts:
   - Add `--no-metadata` flag support
   - Use the report_test_result function for consistent output
   - Add color formatting for readability
   - Implement proper exit code handling

2. Add corresponding targets to the Makefile:
   ```makefile
   test-your-feature: build
       @echo "=== Testing Your Feature ==="
       ./tools/test_your_feature.sh

   test-your-feature-ci: build
       @echo "=== Testing Your Feature in CI ==="
       ./tools/test_your_feature.sh --no-metadata
   ```

3. Add the test to the GitHub Actions workflow in `.github/workflows/integration.yml`

### Adding Unit Tests

Add Go test files following standard Go testing conventions:

```go
func TestYourFeature(t *testing.T) {
    // Test implementation
}
```

## Test Result Interpretation

All test scripts provide consistent output with:

- ✅ Green checkmarks for passing tests
- ❌ Red X for failing tests
- ⚠️ Yellow warnings for skipped or conditional tests

Each script will exit with code 0 if all tests pass, or non-zero if any test fails.

## CI Testing (GitHub Actions)

VidKit uses GitHub Actions for continuous integration testing. Two workflows are defined:

1. **Test**: Runs unit tests on each push and pull request
2. **Integration Tests**: Runs full integration test suite

The CI environment:
- Runs without API keys using `--no-metadata` flags
- Uses Ubuntu latest with FFmpeg installed
- Generates test videos for consistent testing
- Archives test artifacts on failure for debugging

## API Keys and Testing

For complete testing with metadata providers:

1. Create a `config.json` file with your API keys:
   ```json
   {
     "tmdb_api_key": "your-tmdb-api-key",
     "omdb_api_key": "your-omdb-api-key",
     "tvdb_api_key": "your-tvdb-api-key"
   }
   ```

2. Run tests without the `--no-metadata` flag to enable metadata lookups

## Troubleshooting Common Issues

### FFmpeg Not Found

If tests fail with "FFmpeg not found" or similar errors:
- Ensure FFmpeg is installed and available in your PATH
- Verify with `ffmpeg -version` in your terminal

### API Key Issues

If tests fail with "API key is required for metadata lookup":
- Using `--no-metadata` flag or `-ci` test targets will bypass this error
- Alternatively, set up your API keys in `config.json`

### CI Test Failures

If tests pass locally but fail in CI:
- Ensure tests handle `--no-metadata` mode correctly
- Check if tests depend on specific output formats that might differ in CI
- Verify test videos are being generated correctly in the CI environment

## Extending the Test Suite

When adding new features to VidKit, follow these guidelines for testing:

1. Add unit tests for new functions and methods
2. Create or extend integration tests for CLI features
3. Ensure tests work in both local and CI environments
4. Update documentation if test procedures change
