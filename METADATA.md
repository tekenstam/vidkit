# Metadata Integration for VidKit

VidKit supports multiple metadata sources to provide rich information for your video files. This document explains how to set up and use each metadata provider.

## TMDb (The Movie Database)

TMDb is used for retrieving movie metadata, including title, year, and overview.

### 1. Get a TMDb API Key

1. Visit [The Movie Database](https://www.themoviedb.org/)
2. Create an account or sign in
3. Go to your account settings
4. Click on "API" in the left sidebar
5. Request an API key for a "Developer" account
6. Fill in the application form:
   - Application Name: VidKit
   - Application URL: (leave blank for personal use)
   - Application Summary: Personal media file organization tool
   - Use case: Personal use for organizing media files

### 2. Configure TMDb in VidKit

1. Copy the example configuration:
   ```bash
   cp .config.example.json ~/.config/vidkit/config.json
   ```

2. Edit the configuration file:
   ```bash
   nano ~/.config/vidkit/config.json
   ```

3. Add your API key to the configuration:
   ```json
   {
     "tmdb_api_key": "YOUR_API_KEY_HERE"
   }
   ```

### 3. TMDb-specific Options

- `tmdb_api_key`: Your TMDb API key
- `movie_format`: Template for movie filenames
  - `{title}`: Movie title
  - `{year}`: Release year
  - `{resolution}`: Video resolution (e.g., "1080p")
  - `{codec}`: Video codec (e.g., "h264")

## TvMaze

TvMaze is used for retrieving TV show metadata, including show information and episode details.

### 1. TvMaze API

The great news about TvMaze is that it provides a free API that doesn't require an API key for basic usage. VidKit is configured to use the public TvMaze API by default, which provides:

- TV show information (title, year, network, status, genres)
- Season and episode information
- Episode titles and air dates
- Show overviews and descriptions

### 2. TvMaze-specific Options

- `tv_format`: Template for TV episode filenames
  - `{title}`: TV show title
  - `{year}`: Show premiere year
  - `{season}`: Season number
  - `{season:02d}`: Season number with leading zero (01 instead of 1)
  - `{episode}`: Episode number
  - `{episode:02d}`: Episode number with leading zero (01 instead of 1)
  - `{episode_title}`: Episode title
  - `{resolution}`: Video resolution (e.g., "1080p")
  - `{codec}`: Video codec (e.g., "h264")
  
### 3. Rate Limiting

TvMaze has rate limits on its API:
- 20 calls every 10 seconds per IP address
- Please be respectful of these limits when processing large batches

## Filename Format Guidelines

For successful metadata extraction, VidKit requires specific formatting in your video filenames.

### Year Detection

VidKit only recognizes years when they are explicitly enclosed in delimiters:

✅ **Recognized Year Formats:**
- `Movie Title (2023).mp4` - Year will be detected as 2023
- `Movie Title [2023].mp4` - Year will be detected as 2023
- `TV Show (2020) S01E01.mp4` - Year will be detected as 2020

❌ **Unrecognized Year Formats:**
- `Movie.Title.2023.mp4` - Year will NOT be detected (will be treated as part of title)
- `TV.Show.2020.S01E01.mp4` - Year will NOT be detected

This requirement ensures that random 4-digit numbers in titles aren't mistakenly identified as years and improves search accuracy.

### TV Show Detection

VidKit recognizes several TV show episode naming patterns:

- `ShowName S01E02.mp4` - Standard format
- `ShowName.S01E02.mp4` - Scene-style format
- `ShowName 1x02.mp4` - Alternate format
- `ShowName Season 1 Episode 2.mp4` - Full word format

For most accurate metadata:
1. Use the `S01E02` format for season and episode numbers
2. Place years in parentheses, like `Show Name (2020) S01E01`
3. Keep episode titles after the episode number: `Show S01E02 Episode Title`

## General Configuration Options

These options apply to all metadata providers:

- `language`: Preferred language for metadata (e.g., "en", "es", "fr")
- `separator`: Character to use between words in filenames (default: " ")
  - Use " " (space) for standard naming: `Big Buck Bunny (2008) [1080p h264].mp4`
  - Use "." for scene style: `Big.Buck.Bunny.(2008).[1080p.h264].mp4`
  - Use "_" for underscore style: `Big_Buck_Bunny_(2008)_[1080p_h264].mp4`
  - Use "-" for dash style: `Big-Buck-Bunny-(2008)-[1080p-h264].mp4`
- `batch_mode`: Process files without prompting (default: false)
- `recursive`: Search subdirectories (default: false)
- `lower_case`: Use lowercase in filenames (default: false)
- `scene_style`: Use dots instead of spaces (shortcut for separator: ".") (default: false)
- `no_overwrite`: Prevent overwriting existing files (default: true)
- `file_extensions`: List of video file extensions to process
- `no_metadata`: Skip online metadata lookup entirely

## Command Line Options

```bash
vidkit [options] <file_or_directory>

Options:
  -b, --batch      process automatically without interactive prompts
  -r, --recursive  search for files within nested directories
  -l, --lower      rename files using lowercase characters
  -s, --scene      use dots in place of spaces (shortcut for --separator '.')
  --separator      character to use as separator in filenames
  --no-overwrite   prevent relocation if it would overwrite a file
  --lang <code>    metadata language (ISO 639-1 code)
  --movie-format   movie filename format template
  --tv-format      TV episode filename format template
  --preview        show what would be done without making changes
  --no-metadata    skip online metadata lookup
```

## Examples

### Movie Examples

1. Standard naming (with spaces):
   ```bash
   vidkit "movie.mp4"
   # Result: "Movie Title (2023) [1080p h264].mp4"
   ```

2. Scene-style naming (with dots):
   ```bash
   vidkit -s "movie.mp4"
   # Result: "Movie.Title.(2023).[1080p.h264].mp4"
   ```

3. Custom separator (underscore):
   ```bash
   vidkit --separator "_" "movie.mp4"
   # Result: "Movie_Title_(2023)_[1080p_h264].mp4"
   ```

### TV Show Examples

1. Standard TV show format:
   ```bash
   vidkit "Show.S01E01.mp4"
   # Result: "Show S01E01 Episode Title [1080p h264].mp4"
   ```

2. Custom TV format:
   ```bash
   vidkit --tv-format "{title} - S{season:02d}E{episode:02d} - {episode_title}" "Show.S01E01.mp4"
   # Result: "Show - S01E01 - Episode Title.mp4"
   ```

3. Scene-style TV show format:
   ```bash
   vidkit -s "Show.S01E01.mp4"
   # Result: "Show.S01E01.Episode.Title.[1080p.h264].mp4"
   ```

## Troubleshooting

1. **API Key Invalid**: Ensure your TMDb API key is correctly entered in the config file.

2. **Not Finding TV Shows**: Try different search formats. TvMaze works best with clean show titles.

3. **Rate Limiting**: If you see errors about too many requests, slow down batch processing.

4. **Language Issues**: Check that the language code specified is supported by the metadata provider.

5. **Metadata Not Found**: For obscure titles, you might need to manually specify more information.

## Metadata Provider Status

| Provider | Type           | API Key Required | Default Status |
|----------|----------------|------------------|----------------|
| TMDb     | Movies         | Yes              | Enabled        |
| TvMaze   | TV Shows       | No               | Enabled        |
