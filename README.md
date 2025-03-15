# vidkit

A powerful command-line toolkit for video file analysis, organization, and metadata management.

## Features

- Displays comprehensive video information:
  - File metadata (format, size, duration)
  - Video stream details (codec, resolution, bitrate, frame rate)
  - Audio stream details (codec, sample rate, channels, bitrate)
- Online movie metadata lookup:
  - Automatic movie identification using TMDb
  - Smart title and year extraction from filenames
  - Movie overview and details
  - Configurable API key
- Batch processing:
  - Process single files or entire directories
  - Recursive directory scanning
  - Supports common video formats (mp4, mkv, avi, etc.)
- Human-readable resolution standards (e.g., "1080p" instead of just dimensions)
- Supports standard resolutions from 360p to 8K
- Smart resolution detection with 10-pixel tolerance
- Formatted output with proper units (KB/MB/GB, Kbps, fps)
- Smart file renaming with metadata-based format:
  - Pattern: `{title} ({year}) [{resolution} {codec}].{ext}`
  - Example: `Big Buck Bunny (2008) [360p h264].mp4`
- Preview mode to check changes without modifying files

## Prerequisites

- Go 1.21 or later
- FFmpeg (specifically ffprobe) installed on your system
- TMDb API key (required for metadata lookup)

## Installation

1. Build the tool:
```bash
go build -o vidkit
```

2. Set up TMDb API key (required for metadata lookup):
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
./vidkit <video_file>
```

Process all videos in a directory (recursively):
```bash
./vidkit <directory>
```

Preview changes without modifying files:
```bash
./vidkit -preview <file_or_directory>
```

Skip online metadata lookup:
```bash
./vidkit -no-metadata <file_or_directory>
```

Available options:
```
-preview      Preview mode: show what would be done without making changes
-no-metadata  Skip online metadata lookup
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
