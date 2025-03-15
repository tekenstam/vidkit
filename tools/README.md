# VidKit Testing Tools

This directory contains tools to assist in testing VidKit functionality.

## Test Video Generator

The `generate_test_videos.sh` script creates a suite of test video files with various TV show naming conventions. These files are used for testing VidKit's filename parsing and metadata detection.

### Usage

```bash
# Generate all test videos
./generate_test_videos.sh

# This creates files in the test_videos/ directory
```

### Generated Files

The script generates 38 small video files with various naming patterns:

1. **Breaking Bad Files**
   - Standard format: `Breaking.Bad.S01E01.Pilot.mp4`
   - Alt format: `Breaking.Bad.1x01.mp4`
   - Descriptive: `Breaking.Bad.Season.1.Episode.1.mp4`
   - With year: `Breaking Bad (2008) S01E01.mp4`
   - With quality: `Breaking Bad S01E01 720p HEVC.mp4`
   - Scene style: `Breaking.Bad.S01E01.1080p.WEB-DL.x264.mp4`

2. **Game of Thrones Files**
   - Standard format: `Game.of.Thrones.S01E01.Winter.Is.Coming.mp4`
   - Alt format: `Game.of.Thrones.1x01.mp4`
   - With year: `Game of Thrones (2011) S01E01.mp4`
   - With quality: `Game.of.Thrones.S01E01.1080p.BluRay.x264.mp4`

3. **Stranger Things Files**
   - Standard format: `Stranger.Things.S01E01.The.Vanishing.of.Will.Byers.mp4`
   - Alt format: `Stranger.Things.1x01.mp4`
   - With year: `Stranger Things (2016) S01E01.mp4`
   - With quality: `Stranger.Things.S01E01.2160p.NF.WEB-DL.mp4`

4. **The Office Files**
   - Standard format: `The.Office.S01E01.Pilot.mp4`
   - With year: `The Office (2005) S01E01.mp4`
   - With quality: `The.Office.S01E01.720p.HDTV.x264.mp4`

5. **Friends Files**
   - Standard format: `Friends.S01E01.The.One.Where.Monica.Gets.a.Roommate.mp4`
   - With year: `Friends (1994) S01E01.mp4`
   - With quality: `Friends.S01E01.DVDRip.x264.mp4`

6. **Special Cases**
   - Release group: `[GROUP] Breaking Bad S01E05.mp4`
   - Underscores: `Breaking_Bad_-_s01e01_-_Pilot.mp4`
   - Very descriptive: `The Good Place (2016) - Season 1 Episode 1 - Everything Is Fine.mp4`
   - Scene release: `Loki.S01E01.Glorious.Purpose.1080p.DSNP.WEB-DL.DDP5.1.Atmos.H.264-CMRG.mp4`
   - Quality in title: `Severance 1x01 Good News About Hell 2160p.mp4`

### Testing Guidelines

1. Use preview mode when testing to avoid modifying files:
   ```bash
   go run cmd/vidkit/main.go --preview test_videos/*.mp4
   ```

2. Test various naming formats:
   ```bash
   # Test specific format
   go run cmd/vidkit/main.go --preview test_videos/Breaking.Bad.1x01.mp4
   
   # Test with custom format string
   go run cmd/vidkit/main.go --preview --tv-format "{title} - {season}x{episode} - {episode_title}" test_videos/*.mp4
   ```

3. Batch processing:
   ```bash
   go run cmd/vidkit/main.go --preview --batch test_videos/*.mp4
   ```

For more detailed testing information, see the [CONTRIBUTING.md](../CONTRIBUTING.md) file.
