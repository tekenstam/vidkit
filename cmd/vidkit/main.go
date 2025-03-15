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
	if cfg.NoMetadata {
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
	ext := filepath.Ext(originalPath)
	baseDir := filepath.Dir(originalPath)
	
	// Find first video stream
	videoStreamIndex := -1
	for i, stream := range info.Streams {
		if stream.CodecType == "video" {
			videoStreamIndex = i
			break
		}
	}

	// Extract resolution
	resolutionString := "unknown"
	codec := "unknown"
	if videoStreamIndex >= 0 {
		videoStream := info.Streams[videoStreamIndex]
		resolutionString = resolution.GetStandardResolution(videoStream.Width, videoStream.Height)
		codec = videoStream.CodecName
	}

	// Apply movie filename template from configuration
	template := cfg.MovieFilenameTemplate
	if template == "" {
		template = "{title} ({year}) [{resolution} {codec}]"
	}

	// Replace template variables
	filename := template
	filename = strings.ReplaceAll(filename, "{title}", metadata.Title)
	filename = strings.ReplaceAll(filename, "{year}", fmt.Sprintf("%d", metadata.Year))
	filename = strings.ReplaceAll(filename, "{resolution}", resolutionString)
	filename = strings.ReplaceAll(filename, "{codec}", codec)
	
	// Replace genres if needed
	if strings.Contains(filename, "{genre}") && len(metadata.Genres) > 0 {
		filename = strings.ReplaceAll(filename, "{genre}", metadata.Genres[0])
	} else {
		filename = strings.ReplaceAll(filename, "{genre}", "Unknown")
	}

	// Apply word separator
	if cfg.SceneStyle {
		filename = strings.ReplaceAll(filename, " ", ".")
	}

	// Apply lowercase if requested
	if cfg.Lowercase {
		filename = strings.ToLower(filename)
	}

	// Create organized output directory structure using templates if configured
	if metadata.Title != "" && cfg.OrganizeFiles && cfg.MovieDirectoryTemplate != "" {
		// Apply movie directory template
		dirTemplate := cfg.MovieDirectoryTemplate
		if dirTemplate == "" {
			dirTemplate = "Movies/{title} ({year})"
		}
		directory := dirTemplate
		directory = strings.ReplaceAll(directory, "{title}", metadata.Title)
		directory = strings.ReplaceAll(directory, "{year}", fmt.Sprintf("%d", metadata.Year))
		if strings.Contains(directory, "{genre}") && len(metadata.Genres) > 0 {
			directory = strings.ReplaceAll(directory, "{genre}", metadata.Genres[0])
		} else {
			directory = strings.ReplaceAll(directory, "{genre}", "Unknown")
		}
		if cfg.SceneStyle {
			directory = strings.ReplaceAll(directory, " ", ".")
		}
		if cfg.Lowercase {
			directory = strings.ToLower(directory)
		}
		fullPath := filepath.Join(baseDir, directory, filename+ext)
		return fullPath
	}
	
	// If no directory organization, just return the filename
	return filepath.Join(baseDir, filename+ext)
}

func generateTVFilename(originalPath string, info *media.VideoInfo, metadata *metadata.TVShowMetadata, cfg *config.Config) string {
	ext := filepath.Ext(originalPath)
	baseDir := filepath.Dir(originalPath)
	
	// Find first video stream
	videoStreamIndex := -1
	for i, stream := range info.Streams {
		if stream.CodecType == "video" {
			videoStreamIndex = i
			break
		}
	}

	// Extract resolution
	resolutionString := "unknown"
	codec := "unknown"
	if videoStreamIndex >= 0 {
		videoStream := info.Streams[videoStreamIndex]
		resolutionString = resolution.GetStandardResolution(videoStream.Width, videoStream.Height)
		codec = videoStream.CodecName
	}

	// Apply TV show filename template from configuration
	template := cfg.TVFilenameTemplate
	if template == "" {
		template = "{title} - S{season:02d}E{episode:02d} - {episode_title}"
	}

	// Replace template variables
	filename := template
	filename = strings.ReplaceAll(filename, "{title}", metadata.Title)
	filename = strings.ReplaceAll(filename, "{year}", fmt.Sprintf("%d", metadata.Year))
	filename = strings.ReplaceAll(filename, "{resolution}", resolutionString)
	filename = strings.ReplaceAll(filename, "{codec}", codec)
	
	// Handle any special formatting for season/episode
	if strings.Contains(filename, "{season:02d}") {
		filename = strings.ReplaceAll(filename, "{season:02d}", fmt.Sprintf("%02d", metadata.Season))
	} else {
		filename = strings.ReplaceAll(filename, "{season}", fmt.Sprintf("%d", metadata.Season))
	}
	
	if strings.Contains(filename, "{episode:02d}") {
		filename = strings.ReplaceAll(filename, "{episode:02d}", fmt.Sprintf("%02d", metadata.Episode))
	} else {
		filename = strings.ReplaceAll(filename, "{episode}", fmt.Sprintf("%d", metadata.Episode))
	}
	
	filename = strings.ReplaceAll(filename, "{episode_title}", metadata.EpisodeTitle)
	
	// Replace network if needed
	if strings.Contains(filename, "{network}") {
		filename = strings.ReplaceAll(filename, "{network}", metadata.Network)
	}
	
	// Replace genres if needed
	if strings.Contains(filename, "{genre}") && len(metadata.Genres) > 0 {
		filename = strings.ReplaceAll(filename, "{genre}", metadata.Genres[0])
	} else {
		filename = strings.ReplaceAll(filename, "{genre}", "Unknown")
	}

	// Apply word separator
	if cfg.SceneStyle {
		filename = strings.ReplaceAll(filename, " ", ".")
	}

	// Apply lowercase if requested
	if cfg.Lowercase {
		filename = strings.ToLower(filename)
	}

	// Create organized output directory structure using templates if configured
	if cfg.OrganizeFiles && cfg.TVDirectoryTemplate != "" {
		// Apply TV show directory template
		dirTemplate := cfg.TVDirectoryTemplate
		if dirTemplate == "" {
			dirTemplate = "TV/{title}/Season {season:02d}"
		}
		directory := dirTemplate
		directory = strings.ReplaceAll(directory, "{title}", metadata.Title)
		directory = strings.ReplaceAll(directory, "{year}", fmt.Sprintf("%d", metadata.Year))
		if strings.Contains(directory, "{season:02d}") {
			directory = strings.ReplaceAll(directory, "{season:02d}", fmt.Sprintf("%02d", metadata.Season))
		} else {
			directory = strings.ReplaceAll(directory, "{season}", fmt.Sprintf("%d", metadata.Season))
		}
		if strings.Contains(directory, "{network}") {
			directory = strings.ReplaceAll(directory, "{network}", metadata.Network)
		}
		if strings.Contains(directory, "{genre}") && len(metadata.Genres) > 0 {
			directory = strings.ReplaceAll(directory, "{genre}", metadata.Genres[0])
		} else {
			directory = strings.ReplaceAll(directory, "{genre}", "Unknown")
		}
		if cfg.SceneStyle {
			directory = strings.ReplaceAll(directory, " ", ".")
		}
		if cfg.Lowercase {
			directory = strings.ToLower(directory)
		}
		fullPath := filepath.Join(baseDir, directory, filename+ext)
		return fullPath
	}
	
	// If no directory organization, just return the filename
	return filepath.Join(baseDir, filename+ext)
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
	batchMode := flag.Bool("batch", false, "Process files without prompting")
	recursive := flag.Bool("recursive", false, "Process directories recursively")
	lowercase := flag.Bool("lowercase", false, "Convert filenames to lowercase")
	sceneStyle := flag.Bool("scene-style", false, "Use dots instead of spaces (scene style)")
	organize := flag.Bool("organize", false, "Organize files into directories")
	noOverwrite := flag.Bool("no-overwrite", false, "Don't overwrite existing files")
	noMetadata := flag.Bool("no-metadata", false, "Skip metadata lookup")
	previewMode := flag.Bool("preview", false, "Preview mode (don't modify files)")
	showVersion := flag.Bool("version", false, "Show version information")
	
	// Language and filename template options
	lang := flag.String("lang", "en", "Metadata language (ISO 639-1 code)")
	movieFilenameTemplate := flag.String("movie-filename-template", "", "Template for movie filenames (e.g., '{title} ({year}) [{resolution}]')")
	tvFilenameTemplate := flag.String("tv-filename-template", "", "Template for TV show filenames (e.g., '{title} S{season:02d}E{episode:02d} {episode_title}')")
	separator := flag.String("separator", "", "Character to use as separator in filenames")
	movieProvider := flag.String("movie-provider", "", "Select movie metadata provider (tmdb, omdb)")
	tvProvider := flag.String("tv-provider", "", "Select TV show metadata provider (tvmaze, tvdb)")

	// Directory organization templates
	movieDirectoryTemplate := flag.String("movie-directory-template", "", "Template for movie directory organization (e.g., 'Movies/{genre}/{title} ({year})')")
	tvDirectoryTemplate := flag.String("tv-directory-template", "", "Template for TV show directory organization (e.g., 'TV/{genre}/{title}/Season {season:02d}')")
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
	if *batchMode {
		cfg.BatchMode = true
	}

	if *recursive {
		cfg.Recursive = true
	}

	if *lowercase {
		cfg.Lowercase = true
	}

	if *sceneStyle {
		cfg.SceneStyle = true
	}

	if *separator != "" {
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

	if *lang != "en" {
		cfg.Language = *lang
	}

	if *movieFilenameTemplate != "" {
		cfg.MovieFilenameTemplate = *movieFilenameTemplate
	}

	if *tvFilenameTemplate != "" {
		cfg.TVFilenameTemplate = *tvFilenameTemplate
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

	// Set organize flags
	cfg.OrganizeFiles = *organize
	cfg.MovieDirectoryTemplate = *movieDirectoryTemplate
	cfg.TVDirectoryTemplate = *tvDirectoryTemplate

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
