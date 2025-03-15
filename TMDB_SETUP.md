# Setting up TMDb Integration for VidKit

VidKit uses The Movie Database (TMDb) API to fetch accurate metadata for your video files. Here's how to set it up:

## 1. Get a TMDb API Key

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

## 2. Configure VidKit

1. Copy the example configuration:
   ```bash
   cp .config.example.json ~/.config/vidkit/config.json
   ```

2. Edit the configuration file:
   ```bash
   nano ~/.config/vidkit/config.json
   ```

3. Replace `"YOUR_API_KEY_HERE"` with your actual TMDb API key

## 3. Configuration Options

- `tmdb_api_key`: Your TMDb API key
- `language`: Preferred language for metadata (e.g., "en", "es", "fr")
- `separator`: Character to use between words in filenames (default: " ")
  - Use " " (space) for standard naming: `Big Buck Bunny (2008) [1080p h264].mp4`
  - Use "." for scene style: `Big.Buck.Bunny.(2008).[1080p.h264].mp4`
  - Use "_" for underscore style: `Big_Buck_Bunny_(2008)_[1080p_h264].mp4`
  - Use "-" for dash style: `Big-Buck-Bunny-(2008)-[1080p-h264].mp4`
- `movie_format`: Template for movie filenames
  - `{title}`: Movie title
  - `{year}`: Release year
  - `{resolution}`: Video resolution (e.g., "1080p")
  - `{codec}`: Video codec (e.g., "h264")
- `tv_format`: Template for TV episode filenames (future use)
- `batch_mode`: Process files without prompting (default: false)
- `recursive`: Search subdirectories (default: false)
- `lower_case`: Use lowercase in filenames (default: false)
- `scene_style`: Use dots instead of spaces (shortcut for separator: ".") (default: false)
- `no_overwrite`: Prevent overwriting existing files (default: true)
- `file_extensions`: List of video file extensions to process
- `ignore_patterns`: List of patterns to ignore (e.g., "sample", "trailer")

## 4. Command Line Options

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

## 5. Examples

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

4. Lowercase with custom separator:
   ```bash
   vidkit -l --separator "-" "movie.mp4"
   # Result: "movie-title-(2023)-[1080p-h264].mp4"
   ```

## Troubleshooting

1. If you get "TMDb API key is required":
   - Check that config.json exists in ~/.config/vidkit/
   - Verify your API key is correctly set
   - Ensure config.json has correct permissions

2. If metadata lookup fails:
   - Check your internet connection
   - Verify your API key is valid
   - Try with --no-metadata flag to test other features

3. If filename separators aren't working:
   - Check your config.json for the `separator` setting
   - Command-line --separator flag overrides config file
   - -s (scene style) overrides both config and --separator

## Security Note

The config.json file contains your API key. For security:
- Keep it readable only by you: `chmod 600 ~/.config/vidkit/config.json`
- Never share your API key
- Don't commit config.json to version control
