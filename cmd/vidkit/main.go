package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tekenstam/vidkit/internal/pkg/config"
	"github.com/tekenstam/vidkit/internal/pkg/media"
	"github.com/tekenstam/vidkit/internal/pkg/metadata"
	"github.com/tekenstam/vidkit/pkg/resolution"
)

// Version information set by goreleaser
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func processFile(path string, cfg *config.Config) error {
	info, err := media.GetVideoInfo(path)
	if err != nil {
		return fmt.Errorf("error analyzing video: %v", err)
	}

	fmt.Printf("\n=== Processing: %s ===\n", path)

	// Print file information
	fmt.Println("\n=== File Information ===")
	fmt.Printf("Filename: %s\n", info.Format.Filename)
	fmt.Printf("Container Format: %s\n", info.Format.FormatName)
	fmt.Printf("Duration: %s seconds\n", info.Format.Duration)
	fmt.Printf("File Size: %s\n", media.FormatSize(info.Format.Size))
	fmt.Printf("Overall Bitrate: %s\n", media.FormatBitRate(info.Format.BitRate))

	// Print stream information
	for _, stream := range info.Streams {
		if stream.CodecType == "video" {
			fmt.Println("\n=== Video Stream ===")
			fmt.Printf("Codec: %s\n", stream.CodecName)
			fmt.Printf("Resolution: %s\n", resolution.GetStandardResolution(stream.Width, stream.Height))
			fmt.Printf("Bitrate: %s\n", media.FormatBitRate(stream.BitRate))
			fmt.Printf("Frame Rate: %s\n", media.FormatFrameRate(stream.FrameRate))
		} else if stream.CodecType == "audio" {
			fmt.Println("\n=== Audio Stream ===")
			fmt.Printf("Codec: %s\n", stream.CodecName)
			if stream.SampleRate != "" {
				fmt.Printf("Sample Rate: %s Hz\n", stream.SampleRate)
			}
			fmt.Printf("Channels: %d\n", stream.Channels)
			fmt.Printf("Channel Layout: %s\n", stream.ChannelLayout)
			fmt.Printf("Bitrate: %s\n", media.FormatBitRate(stream.BitRate))
		}
	}

	// Skip metadata lookup if requested
	if cfg.NoMetadata || !cfg.EnableMetadata {
		return nil
	}

	// Check if this is a TV show
	tvShowInfo := metadata.ExtractTVShowInfo(path)
	if tvShowInfo.Season > 0 && tvShowInfo.Episode > 0 {
		// This is a TV show, process it accordingly
		return processTVShow(path, info, tvShowInfo, cfg)
	}

	// If not a TV show, treat as movie
	movieInfo := metadata.ExtractMovieInfo(path)
	if movieInfo.Title != "" {
		// This appears to be a movie
		return processMovie(path, info, movieInfo, cfg)
	}

	return nil
}

func processTVShow(path string, info *media.VideoInfo, tvShowInfo metadata.TVShowSearch, cfg *config.Config) error {
	fmt.Println("\n=== Looking up TV show metadata... ===")

	// Create a formatted search string
	searchString := fmt.Sprintf("'%s'", tvShowInfo.Title)
	if tvShowInfo.Year > 0 {
		searchString = fmt.Sprintf("%s (year: %d)", searchString, tvShowInfo.Year)
	}
	searchString = fmt.Sprintf("%s - S%02dE%02d", searchString, tvShowInfo.Season, tvShowInfo.Episode)
	fmt.Printf("Searching for %s\n", searchString)

	// Create the appropriate provider using factory
	provider, err := metadata.CreateTVShowProvider(cfg)
	if err != nil {
		return fmt.Errorf("failed to create TV show provider: %v", err)
	}

	// Search for the TV show
	tvShowMetadata, err := provider.SearchTVShow(tvShowInfo, cfg.Language)
	if err != nil {
		// Just log the error and continue without metadata
		fmt.Printf("Warning: Failed to look up TV show: %v\n", err)
		return nil
	}

	// Print TV show metadata
	fmt.Println("\n=== TV Show Metadata ===")
	fmt.Printf("Title: %s\n", tvShowMetadata.Title)
	if tvShowMetadata.Year > 0 {
		fmt.Printf("Year: %d\n", tvShowMetadata.Year)
	}
	if tvShowMetadata.Network != "" {
		fmt.Printf("Network: %s\n", tvShowMetadata.Network)
	}
	if tvShowMetadata.Status != "" {
		fmt.Printf("Status: %s\n", tvShowMetadata.Status)
	}
	if len(tvShowMetadata.Genres) > 0 {
		fmt.Printf("Genres: %s\n", strings.Join(tvShowMetadata.Genres, ", "))
	}
	if tvShowMetadata.SeasonCount > 0 {
		fmt.Printf("Season Count: %d\n", tvShowMetadata.SeasonCount)
	}

	// Print episode information
	fmt.Println("\n=== Episode Information ===")
	fmt.Printf("Season: %d\n", tvShowMetadata.Season)
	fmt.Printf("Episode: %d\n", tvShowMetadata.Episode)
	if tvShowMetadata.EpisodeTitle != "" {
		fmt.Printf("Title: %s\n", tvShowMetadata.EpisodeTitle)
	}
	if tvShowMetadata.AirDate != "" {
		fmt.Printf("Air Date: %s\n", tvShowMetadata.AirDate)
	}
	if tvShowMetadata.Overview != "" {
		fmt.Printf("Overview: %s\n", tvShowMetadata.Overview)
	}

	// Generate a new filename using the metadata
	newFileName := generateTVFilename(path, info, tvShowMetadata, cfg)

	// Show rename preview
	fmt.Println("\n=== File Renaming ===")
	fmt.Printf("Original: %s\n", path)
	fmt.Printf("New name: %s\n", newFileName)

	// Skip renaming if preview mode
	if cfg.PreviewMode {
		fmt.Println("\n[PREVIEW MODE] File would be renamed as shown above")
		return nil
	}

	// Skip if the target file already exists
	if cfg.NoOverwrite && fileExists(newFileName) {
		fmt.Println("\nSkipping rename: Target file already exists")
		return nil
	}

	// In batch mode, rename without confirmation
	if cfg.BatchMode || confirmRename() {
		err := os.Rename(path, newFileName)
		if err != nil {
			return fmt.Errorf("error renaming file: %v", err)
		}
		fmt.Println("File renamed successfully!")
	}

	return nil
}

func processMovie(path string, info *media.VideoInfo, movieInfo metadata.MovieSearch, cfg *config.Config) error {
	fmt.Println("\n=== Looking up movie metadata... ===")

	// Create a formatted search string
	searchString := fmt.Sprintf("'%s'", movieInfo.Title)
	if movieInfo.Year > 0 {
		searchString = fmt.Sprintf("%s (year: %d)", searchString, movieInfo.Year)
	}
	fmt.Printf("Searching for %s...\n", searchString)

	// Create the appropriate provider using factory
	provider, err := metadata.CreateMovieProvider(cfg)
	if err != nil {
		return fmt.Errorf("failed to create movie provider: %v", err)
	}

	// Search for the movie
	movieMetadata, err := provider.SearchMovie(movieInfo, cfg.Language)
	if err != nil {
		// Just log the error and continue without metadata
		fmt.Printf("Warning: Failed to look up movie: %v\n", err)
		return nil
	}

	// Print movie metadata
	fmt.Println("\n=== Movie Metadata ===")
	fmt.Printf("Title: %s\n", movieMetadata.Title)
	if movieMetadata.Year > 0 {
		fmt.Printf("Year: %d\n", movieMetadata.Year)
	}
	if movieMetadata.Overview != "" {
		fmt.Printf("Overview: %s\n", movieMetadata.Overview)
	}

	// Generate a new filename using the metadata
	newFileName := generateFilename(path, info, movieMetadata, cfg)

	// Show rename preview
	fmt.Println("\n=== File Renaming ===")
	fmt.Printf("Original: %s\n", path)
	fmt.Printf("New name: %s\n", newFileName)

	// Skip renaming if preview mode
	if cfg.PreviewMode {
		fmt.Println("\n[PREVIEW MODE] File would be renamed as shown above")
		return nil
	}

	// Skip if the target file already exists
	if cfg.NoOverwrite && fileExists(newFileName) {
		fmt.Println("\nSkipping rename: Target file already exists")
		return nil
	}

	// In batch mode, rename without confirmation
	if cfg.BatchMode || confirmRename() {
		err := os.Rename(path, newFileName)
		if err != nil {
			return fmt.Errorf("error renaming file: %v", err)
		}
		fmt.Println("File renamed successfully!")
	}

	return nil
}

func processPath(path string, cfg *config.Config) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("error accessing path: %v", err)
	}

	if info.IsDir() {
		// Process directory
		entries, err := os.ReadDir(path)
		if err != nil {
			return fmt.Errorf("error reading directory: %v", err)
		}

		for _, entry := range entries {
			entryPath := filepath.Join(path, entry.Name())
			entryInfo, err := entry.Info()
			if err != nil {
				fmt.Printf("Warning: Could not get info for %s: %v\n", entryPath, err)
				continue
			}

			if entryInfo.IsDir() {
				if cfg.Recursive {
					// Process subdirectory recursively
					if err := processPath(entryPath, cfg); err != nil {
						fmt.Printf("Warning: Error processing %s: %v\n", entryPath, err)
					}
				}
			} else {
				// Check if it's a supported video file
				ext := strings.ToLower(filepath.Ext(entryPath))
				isSupported := false
				for _, supportedExt := range cfg.FileExtensions {
					if ext == supportedExt {
						isSupported = true
						break
					}
				}

				if isSupported {
					// Process video file
					if err := processFile(entryPath, cfg); err != nil {
						fmt.Printf("Warning: Error processing %s: %v\n", entryPath, err)
					}
				}
			}
		}
	} else {
		// Process single file
		ext := strings.ToLower(filepath.Ext(path))
		isSupported := false
		for _, supportedExt := range cfg.FileExtensions {
			if ext == supportedExt {
				isSupported = true
				break
			}
		}

		if isSupported {
			if err := processFile(path, cfg); err != nil {
				return fmt.Errorf("error processing file: %v", err)
			}
		} else {
			return fmt.Errorf("unsupported file type: %s", path)
		}
	}

	return nil
}

func generateFilename(originalPath string, info *media.VideoInfo, metadata *metadata.MovieMetadata, cfg *config.Config) string {
	dir := filepath.Dir(originalPath)
	ext := filepath.Ext(originalPath)

	// Find resolution and codec
	resolutionStr := ""
	codec := ""
	for _, stream := range info.Streams {
		if stream.CodecType == "video" {
			resolutionStr = resolution.GetStandardResolution(stream.Width, stream.Height)
			codec = stream.CodecName
			break
		}
	}

	// Default format
	format := cfg.MovieFormat

	// Replace variables in format string
	name := format

	// Replace title
	name = strings.ReplaceAll(name, "{title}", metadata.Title)

	// Replace year if available
	yearStr := ""
	if metadata.Year > 0 {
		yearStr = fmt.Sprintf("%d", metadata.Year)
	}
	name = strings.ReplaceAll(name, "{year}", yearStr)

	// Replace resolution
	name = strings.ReplaceAll(name, "{resolution}", resolutionStr)

	// Replace codec
	name = strings.ReplaceAll(name, "{codec}", codec)

	// Apply lowercase if configured
	if cfg.LowerCase {
		name = strings.ToLower(name)
	}

	// Replace spaces with separator
	if cfg.Separator != " " {
		name = strings.ReplaceAll(name, " ", cfg.Separator)
	}

	// Create new filename
	return filepath.Join(dir, name+ext)
}

func generateTVFilename(originalPath string, info *media.VideoInfo, metadata *metadata.TVShowMetadata, cfg *config.Config) string {
	dir := filepath.Dir(originalPath)
	ext := filepath.Ext(originalPath)

	// Find resolution and codec
	resolutionStr := ""
	codec := ""
	for _, stream := range info.Streams {
		if stream.CodecType == "video" {
			resolutionStr = resolution.GetStandardResolution(stream.Width, stream.Height)
			codec = stream.CodecName
			break
		}
	}

	// Default format
	format := cfg.TVFormat

	// Replace variables in format string
	name := format

	// Replace title
	name = strings.ReplaceAll(name, "{title}", metadata.Title)

	// Replace year if available
	yearStr := ""
	if metadata.Year > 0 {
		yearStr = fmt.Sprintf("%d", metadata.Year)
	}
	name = strings.ReplaceAll(name, "{year}", yearStr)

	// Replace season
	season := fmt.Sprintf("%d", metadata.Season)
	seasonWithZero := fmt.Sprintf("%02d", metadata.Season)
	name = strings.ReplaceAll(name, "{season:02d}", seasonWithZero)
	name = strings.ReplaceAll(name, "{season}", season)

	// Replace episode
	episode := fmt.Sprintf("%d", metadata.Episode)
	episodeWithZero := fmt.Sprintf("%02d", metadata.Episode)
	name = strings.ReplaceAll(name, "{episode:02d}", episodeWithZero)
	name = strings.ReplaceAll(name, "{episode}", episode)

	// Replace episode title
	episodeTitle := metadata.EpisodeTitle
	if episodeTitle == "" {
		episodeTitle = "Episode " + episode
	}
	name = strings.ReplaceAll(name, "{episode_title}", episodeTitle)

	// Replace resolution
	name = strings.ReplaceAll(name, "{resolution}", resolutionStr)

	// Replace codec
	name = strings.ReplaceAll(name, "{codec}", codec)

	// Apply lowercase if configured
	if cfg.LowerCase {
		name = strings.ToLower(name)
	}

	// Replace spaces with separator
	if cfg.Separator != " " {
		name = strings.ReplaceAll(name, " ", cfg.Separator)
	}

	// Create new filename
	return filepath.Join(dir, name+ext)
}

func confirmRename() bool {
	fmt.Print("\nDo you want to rename the file? (y/N): ")
	var response string
	fmt.Scanln(&response)
	return strings.ToLower(response) == "y"
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func main() {
	// Define command-line flags
	batchMode := flag.Bool("b", false, "Process automatically without prompts")
	batchModeLong := flag.Bool("batch", false, "Process automatically without prompts")

	recursiveMode := flag.Bool("r", false, "Search subdirectories")
	recursiveModeLong := flag.Bool("recursive", false, "Search subdirectories")

	lowerCase := flag.Bool("l", false, "Use lowercase for filenames")
	lowerCaseLong := flag.Bool("lower", false, "Use lowercase for filenames")

	sceneStyle := flag.Bool("s", false, "Use dots instead of spaces in filenames")
	sceneStyleLong := flag.Bool("scene", false, "Use dots instead of spaces in filenames")

	separator := flag.String("separator", " ", "Character to use as separator in filenames")

	noOverwrite := flag.Bool("no-overwrite", true, "Prevent renaming if it would overwrite a file")

	previewMode := flag.Bool("preview", false, "Show what would be done without making changes")

	noMetadata := flag.Bool("no-metadata", false, "Skip online metadata lookup")

	language := flag.String("lang", "en", "Metadata language (ISO 639-1 code)")

	movieFormat := flag.String("movie-format", "", "Custom format for movie files")
	tvFormat := flag.String("tv-format", "", "Custom format for TV show files")

	showVersion := flag.Bool("version", false, "Show version information")

	// Add new flags for provider selection
	movieProvider := flag.String("movie-provider", "", "Movie metadata provider (tmdb, omdb)")
	tvProvider := flag.String("tv-provider", "", "TV show metadata provider (tvmaze, tvdb)")

	// Parse command-line flags
	flag.Parse()

	// Show version and exit if requested
	if *showVersion {
		fmt.Printf("VidKit %s (commit %s, built on %s)\n", version, commit, date)
		return
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Warning: Failed to load config: %v\n", err)
		// Continue with default config
		cfg = config.DefaultConfig()
	}

	// Override config with command-line flags
	if *batchMode || *batchModeLong {
		cfg.BatchMode = true
	}

	if *recursiveMode || *recursiveModeLong {
		cfg.Recursive = true
	}

	if *lowerCase || *lowerCaseLong {
		cfg.LowerCase = true
	}

	if *sceneStyle || *sceneStyleLong {
		cfg.SceneStyle = true
	}

	if *separator != " " {
		cfg.Separator = *separator
	}

	if !*noOverwrite {
		cfg.NoOverwrite = false
	}

	if *previewMode {
		cfg.PreviewMode = true
	}

	if *noMetadata {
		cfg.NoMetadata = true
	}

	if *language != "en" {
		cfg.Language = *language
	}

	if *movieFormat != "" {
		cfg.MovieFormat = *movieFormat
	}

	if *tvFormat != "" {
		cfg.TVFormat = *tvFormat
	}

	// Apply provider selection from flags
	if *movieProvider != "" {
		switch *movieProvider {
		case "tmdb":
			cfg.MovieProvider = config.ProviderTMDb
		case "omdb":
			cfg.MovieProvider = config.ProviderOMDb
		default:
			fmt.Printf("Warning: Unknown movie provider '%s', using default\n", *movieProvider)
		}
	}

	if *tvProvider != "" {
		switch *tvProvider {
		case "tvmaze":
			cfg.TVProvider = config.ProviderTVMaze
		case "tvdb":
			cfg.TVProvider = config.ProviderTVDb
		default:
			fmt.Printf("Warning: Unknown TV provider '%s', using default\n", *tvProvider)
		}
	}

	// Set config back so it's available for next time
	if err := config.ValidateConfig(cfg); err != nil {
		fmt.Printf("Error in configuration: %v\n", err)
		os.Exit(1)
	}

	// Check if we have any paths to process
	if flag.NArg() == 0 {
		fmt.Println("Usage: vidkit [options] <file_or_directory>")
		flag.PrintDefaults()
		return
	}

	// Process each path
	for _, path := range flag.Args() {
		if err := processPath(path, cfg); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}
