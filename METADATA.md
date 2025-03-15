# Metadata Integration for VidKit

VidKit now supports multiple metadata providers for both movies and TV shows, giving you flexibility in how you obtain information about your video files. This document explains how to set up and use each metadata provider.

## Supported Providers

| Provider | Type           | API Key Required | Default Status | Command Line Flag   |
|----------|----------------|------------------|----------------|---------------------|
| TMDb     | Movies         | Yes              | Default        | `--movie-provider tmdb` |
| OMDb     | Movies         | Yes              | Optional       | `--movie-provider omdb` |
| TvMaze   | TV Shows       | No               | Default        | `--tv-provider tvmaze` |
| TVDb     | TV Shows       | Yes              | Optional       | `--tv-provider tvdb` |

## 1. Movie Metadata Providers

### TMDb (The Movie Database)

TMDb is the default provider for movie metadata, offering comprehensive information about films.

#### TMDb Setup

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

#### TMDb Configuration

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
     "tmdb_api_key": "YOUR_API_KEY_HERE",
     "movie_provider": "tmdb"
   }
   ```

### OMDb (Open Movie Database)

OMDb provides access to a large collection of movie information and can be used as an alternative to TMDb.

#### OMDb Setup

1. Visit [OMDb API](https://www.omdbapi.com/apikey.aspx)
2. Select the Free tier (1,000 daily requests limit) or pay for a higher tier
3. Fill out the form and submit
4. You'll receive an email with a verification link
5. Click the verification link to activate your API key

#### OMDb Configuration

1. Edit your configuration file:
   ```bash
   nano ~/.config/vidkit/config.json
   ```

2. Add your OMDb API key:
   ```json
   {
     "omdb_api_key": "YOUR_OMDB_KEY_HERE",
     "movie_provider": "omdb"
   }
   ```

#### Using OMDb from Command Line

```bash
vidkit --movie-provider omdb movie.mp4
```

## 2. TV Show Metadata Providers

### TvMaze

TvMaze is the default provider for TV show metadata and doesn't require an API key for basic usage.

#### TvMaze Features

- TV show information (title, year, network, status, genres)
- Season and episode information
- Episode titles and air dates
- Show overviews and descriptions

#### TvMaze Configuration

No configuration is required for TvMaze as it's free to use without an API key.

### TVDb (The TV Database)

TVDb provides comprehensive TV show information but requires an API key.

#### TVDb Setup

1. Visit [TheTVDB.com](https://thetvdb.com/)
2. Create an account
3. Go to your account dashboard
4. Navigate to the API section
5. Register for an API key (follow their instructions)

#### TVDb Configuration

1. Edit your configuration file:
   ```bash
   nano ~/.config/vidkit/config.json
   ```

2. Add your TVDb API key:
   ```json
   {
     "tvdb_api_key": "YOUR_TVDB_KEY_HERE",
     "tv_provider": "tvdb"
   }
   ```

#### Using TVDb from Command Line

```bash
vidkit --tv-provider tvdb tvshow.mp4
```

## Configuration Options

You can set your preferred providers in the config file:

```json
{
  "tmdb_api_key": "your_tmdb_api_key",
  "omdb_api_key": "your_omdb_api_key",
  "tvdb_api_key": "your_tvdb_api_key",
  "movie_provider": "tmdb",  // "tmdb" or "omdb"
  "tv_provider": "tvmaze"    // "tvmaze" or "tvdb"
}
```

## Using Multiple Providers

VidKit allows you to switch between providers without changing your configuration:

```bash
# Use TMDb for movie lookup
vidkit --movie-provider tmdb movie.mp4

# Use OMDb for movie lookup
vidkit --movie-provider omdb movie.mp4

# Use TvMaze for TV show lookup
vidkit --tv-provider tvmaze tvshow.mp4

# Use TVDb for TV show lookup
vidkit --tv-provider tvdb tvshow.mp4
```

## Provider Comparison

### Movie Providers

| Feature          | TMDb                 | OMDb                |
|------------------|----------------------|---------------------|
| Movie Data       | Comprehensive        | Good                |
| Update Frequency | Very frequent        | Regular             |
| API Limitations  | 40 requests/10s      | 1,000/day (free)    |
| Language Support | Multiple languages   | Limited             |
| Data Richness    | Very detailed        | Good                |

### TV Show Providers

| Feature          | TvMaze               | TVDb                |
|------------------|----------------------|---------------------|
| API Key Required | No                   | Yes                 |
| Update Frequency | Regular              | Very frequent       |
| API Limitations  | 20 calls/10s per IP  | Varies by account   |
| Community Data   | Moderate             | Very active         |
| Data Richness    | Good                 | Very detailed       |

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
  --movie-provider select movie metadata provider (tmdb, omdb)
  --tv-provider    select TV show metadata provider (tvmaze, tvdb)
```

## Troubleshooting

1. **API Key Invalid**: Ensure your API keys are correctly entered in the config file.

2. **Not Finding TV Shows/Movies**: Try different search formats or use a different provider.

3. **Rate Limiting**: If you see errors about too many requests, slow down batch processing.

4. **Language Issues**: Check that the language code specified is supported by the metadata provider.

5. **Provider Selection**: If you get errors about missing API keys, ensure you have configured the appropriate key for your selected provider.
