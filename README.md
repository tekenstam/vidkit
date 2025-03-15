# vidkit

VidKit is an intelligent video file management system for media enthusiasts and collectors. This powerful Go-based CLI tool organizes your media library with metadata-driven intelligence - analyzing video files, fetching metadata from trusted sources, and intelligently renaming and organizing your content following configurable patterns.

[![Release](https://img.shields.io/github/v-release/tekenstam/vidkit)](https://github.com/tekenstam/vidkit/releases/latest)
[![Go Version](https://img.shields.io/github/go-mod/go-version/tekenstam/vidkit)](https://golang.org/doc/devel/release.html)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Test and Lint](https://github.com/tekenstam/vidkit/actions/workflows/test.yml/badge.svg)](https://github.com/tekenstam/vidkit/actions/workflows/test.yml)
[![Integration Tests](https://github.com/tekenstam/vidkit/actions/workflows/integration.yml/badge.svg)](https://github.com/tekenstam/vidkit/actions/workflows/integration.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/tekenstam/vidkit)](https://goreportcard.com/report/github.com/tekenstam/vidkit)
[![Maintained](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://github.com/tekenstam/vidkit/commits/main)
[![Dependabot](https://img.shields.io/badge/Dependabot-enabled-brightgreen.svg)](https://github.com/tekenstam/vidkit/blob/main/.github/dependabot.yml)

## Overview

VidKit solves the chaos of unorganized media collections by combining technical video analysis with rich metadata to create a perfectly structured library. Whether you're organizing a personal collection or preparing content for a media server, VidKit handles the heavy lifting - from video codec detection to fetching episode titles and organizing by genre.

Built with Go 1.21+ and powered by FFmpeg, VidKit brings professional media management capabilities to your command line.

## Features

- Displays comprehensive video information:
  - File metadata (format, size, duration)
  - Video stream details (codec, resolution, bitrate, frame rate)
  - Audio stream details (codec, sample rate, channels, bitrate)
- Multiple metadata provider options:
  - Movies: TMDb (default) or OMDb
  - TV Shows: TvMaze (default) or TVDb
  - Command-line provider selection
  - Configurable API keys
- Online movie metadata lookup:
  - Automatic movie identification using filename
  - Smart title and year extraction from filenames
  - Movie overview and details
  - Configurable API key
- TV show metadata lookup:
  - Automatic TV show identification using filenames
  - Smart extraction of series name, season and episode from filenames
  - TV show overview and episode details
  - Support for various episode naming conventions
- Intelligent media organization:
  - Customizable directory structure templates
  - Organize by genre, title, year, and more
  - First-letter categorization for large libraries
  - Separate templates for movies and TV shows
- Batch processing:
  - Process single files or entire directories
  - Recursive directory scanning
  - Preview mode to see changes without applying them

## Prerequisites

- FFmpeg (specifically ffprobe) installed on your system
- TMDb API key (optional, required only for metadata lookup)

## Installation

### Option 1: Download Pre-built Binary (Recommended)

1. Download the latest binary for your platform from the [Releases page](https://github.com/tekenstam/vidkit/releases/latest)

**Linux/macOS**:
```bash
# Download the appropriate binary for your system (examples below)
# For macOS (Apple Silicon):
curl -L https://github.com/tekenstam/vidkit/releases/latest/download/vidkit_darwin_arm64.tar.gz -o vidkit.tar.gz

# For macOS (Intel):
curl -L https://github.com/tekenstam/vidkit/releases/latest/download/vidkit_darwin_amd64.tar.gz -o vidkit.tar.gz

# For Linux (x86_64):
curl -L https://github.com/tekenstam/vidkit/releases/latest/download/vidkit_linux_amd64.tar.gz -o vidkit.tar.gz

# Extract and install
tar -xzf vidkit.tar.gz
chmod +x vidkit
sudo mv vidkit /usr/local/bin/  # Optional: move to a directory in your PATH
```

**Windows**:
1. Download the Windows zip file (`vidkit_windows_amd64.zip`)
2. Extract the zip file
3. Optionally, add the extracted directory to your PATH

### Option 2: Build from Source

For developers working on VidKit or users who prefer to build from the latest source:

```bash
git clone https://github.com/tekenstam/vidkit.git
cd vidkit
go build -o vidkit ./cmd/vidkit
```

See the [Installation](#installation) section for simpler options using pre-built binaries.

### API Key Setup (for metadata features)

Set up a TMDb API key (required only if you want to use metadata features):
- Get a free API key from [TMDb](https://www.themoviedb.org/settings/api)
- Add it to `~/.config/vidkit/config.json`:
```json
{
  "tmdb_api_key": "your_api_key_here"
}
```

The tool will create a template config file if none exists and guide you through the setup process.

## Usage

Process a single video file:
```bash
vidkit <video_file>
```

Process all videos in a directory (recursively):
```bash
vidkit --recursive <directory>
```

Preview changes without modifying files:
```bash
vidkit --preview <file_or_directory>
```

Skip online metadata lookup:
```bash
vidkit --no-metadata <file_or_directory>
```

Available options:
```
-b             Batch mode: process automatically without interactive prompts
-r             Search recursively in directories
-l             Use lowercase characters in filenames
-s             Use dots in place of spaces (scene style)
--preview       Preview mode: show what would be done without making changes
--no-metadata   Skip online metadata lookup
--no-overwrite  Prevent renaming if it would overwrite a file
--lang <code>   Metadata language (ISO 639-1 code, default: en)
--movie-filename-template  Custom template for movie filenames (e.g., "{title} ({year}) [{resolution}]")
--tv-filename-template     Custom template for TV show filenames (e.g., "{title} S{season:02d}E{episode:02d} {episode_title}")
--separator     Character to use as separator in filenames
--movie-provider Select movie metadata provider (tmdb, omdb)
--tv-provider    Select TV show metadata provider (tvmaze, tvdb)
--organize      Enable/disable organizing files into directory structures (default: true)
--movie-directory-template     Directory template for movies (e.g., "Movies/{title[0]}/{title} ({year})")
--tv-directory-template        Directory template for TV shows (e.g., "TV/{title}/Season {season:02d}")
--version       Show version information
```

### Supported Video Formats

The tool automatically detects and processes these video formats:
- .mp4  - MPEG-4 Part 14
- .mkv  - Matroska Video
- .avi  - Audio Video Interleave
- .mov  - QuickTime Movie
- .wmv  - Windows Media Video
- .flv  - Flash Video
- .webm - WebM
- .m4v  - MPEG-4 Video
- .mpg  - MPEG-1 Video
- .mpeg - MPEG-1 Video
- .3gp  - 3GPP Multimedia

### Processing Steps

For each video file, the tool will:
1. Display comprehensive video information
2. Look up movie metadata online (if API key is configured)
   - Uses filename year to improve search accuracy
   - Falls back to title-only search if needed
3. Show the proposed new filename
4. In preview mode: Show what would happen without making changes
5. In normal mode: Ask if you want to rename the file

Example output:
```
=== File Information ===
Filename: Big Buck Bunny (2008).mp4
Container Format: mov,mp4,m4a,3gp,3g2,mj2
Duration: 596.458333 seconds
File Size: 61.66 MB
Overall Bitrate: 867.21 Kbps

=== Video Stream ===
Codec: h264
Resolution: 180p
Bitrate: 702.65 Kbps
Frame Rate: 24.00 fps

=== Audio Stream ===
Codec: aac
Sample Rate: 48000 Hz
Channels: 2
Channel Layout: stereo
Bitrate: 160.00 Kbps

=== Looking up metadata... ===
Searching for 'Big Buck Bunny' (year: 2008)...

=== Movie Metadata ===
Title: Big Buck Bunny
Year: 2008
Overview: Big Buck Bunny tells the story of a giant rabbit...

=== File Renaming ===
Original: Big Buck Bunny (2008).mp4
New name: Big Buck Bunny (2008) [180p h264].mp4

Do you want to rename the file? (y/N):
```

## Supported Resolution Standards

- 360p (640x360)
- 480p (640x480)
- 720p (1280x720)
- 1080p (1920x1080)
- 1440p (2560x1440)
- 2K (2048x1080)
- 4K (3840x2160)
- 8K (7680x4320)

## Development

VidKit is built with Go and follows standard Go project practices.

### Code Quality

VidKit uses Go's native linting tools to maintain code quality:

- `go vet` - Checks for common code issues and bugs
- `go fmt` - Ensures consistent code formatting

To run these checks locally:

```bash
make lint      # Runs all linting checks
make test      # Runs all unit tests
make quality   # Runs both tests and linting
```

### Dependencies

- Go 1.21 or higher
- FFmpeg (for video analysis)

## Testing

VidKit includes various tests for ensuring proper functionality.

### Running Tests

Run the full test suite:
```bash
go test ./...
```

Run tests with verbose output:
```bash
go test -v ./...
```

Run specific package tests:
```bash
go test ./internal/pkg/metadata
```

### Testing TV Show Functionality

To test the TV show metadata functionality with sample files:

1. Generate test videos (requires FFmpeg):
```bash
./tools/generate_test_videos.sh
```

2. Test TV show detection in preview mode:
```bash
vidkit --preview test_videos/Breaking.Bad.S01E01.Pilot.mp4
```

3. Test batch processing of multiple formats:
```bash
vidkit --preview --batch test_videos/*.mp4
```

### Testing Different Naming Formats

Test with custom TV show format:
```bash
vidkit --preview --tv-filename-template "{title}.S{season:02d}E{episode:02d}.{episode_title}" test_videos/Breaking.Bad.S01E05.Gray.Matter.mp4
```

Test with scene-style naming (dots instead of spaces):
```bash
vidkit --preview --scene-style test_videos/Breaking.Bad.S01E05.Gray.Matter.mp4
```

For more detailed testing information, see [CONTRIBUTING.md](CONTRIBUTING.md).

## Filename Format Guidelines

For optimal metadata extraction, VidKit has specific rules for recognizing different components in filenames:

### Year Detection

Years must be enclosed in parentheses or square brackets to be recognized:
- `The Matrix (1999).mp4` - Year will be detected as 1999
- `The Matrix [1999].mp4` - Year will be detected as 1999
- `The.Matrix.1999.mp4` - Year will NOT be detected (will be treated as part of title)

This explicit delimiter requirement helps avoid false positives when numbers appear in titles.
