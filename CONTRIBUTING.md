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

For comprehensive testing documentation, please refer to [TESTING.md](TESTING.md).

### Running Basic Tests

```bash
# Run all tests
go test ./...

# Run tests with race condition detection
go test -race ./...
```

When submitting a pull request, make sure all tests pass with:

```bash
make test-all
```

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
