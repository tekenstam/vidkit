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

func processFile(path string, cfg *config.Config, metadataProvider *metadata.TMDbProvider) error {
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
			if stream.Channels > 0 {
				fmt.Printf("Channels: %d\n", stream.Channels)
			}
			if stream.ChannelLayout != "" {
				fmt.Printf("Channel Layout: %s\n", stream.ChannelLayout)
			}
			fmt.Printf("Bitrate: %s\n", media.FormatBitRate(stream.BitRate))
		}
	}

	// Look up movie metadata if enabled
	var movieData *metadata.MovieMetadata
	if !cfg.NoMetadata {
		fmt.Println("\n=== Looking up metadata... ===")
		searchInfo := metadata.ExtractMovieInfo(path)
		if searchInfo.Year > 0 {
			fmt.Printf("Searching for '%s' (year: %d)...\n", searchInfo.Title, searchInfo.Year)
		} else {
			fmt.Printf("Searching for '%s'...\n", searchInfo.Title)
		}

		movieData, err = metadataProvider.SearchMovie(searchInfo, cfg.Language)
		if err != nil {
			fmt.Printf("Failed to lookup metadata: %v\n", err)
		} else {
			fmt.Println("\n=== Movie Metadata ===")
			fmt.Printf("Title: %s\n", movieData.Title)
			fmt.Printf("Year: %d\n", movieData.Year)
			fmt.Printf("Overview: %s\n", movieData.Overview)
		}
	}

	// Generate new filename
	var newName string
	if movieData != nil {
		newName = generateFilename(path, info, movieData, cfg)
	} else {
		newName = path
	}

	fmt.Println("\n=== File Renaming ===")
	fmt.Printf("Original: %s\n", path)
	fmt.Printf("New name: %s\n", newName)

	if cfg.PreviewMode {
		fmt.Println("\n[PREVIEW MODE] File would be renamed as shown above")
		return nil
	}

	// In batch mode or if user confirms
	if cfg.BatchMode || confirmRename() {
		if cfg.NoOverwrite && fileExists(newName) {
			return fmt.Errorf("target file already exists: %s", newName)
		}
		if err := os.Rename(path, newName); err != nil {
			return fmt.Errorf("error renaming file: %v", err)
		}
		fmt.Println("File renamed successfully!")
	}

	return nil
}

func processPath(path string, cfg *config.Config, metadataProvider *metadata.TMDbProvider) error {
	// Get file info
	fileInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("error accessing path: %v", err)
	}

	// If it's a file, process it directly
	if !fileInfo.IsDir() {
		if !media.IsVideoFile(path, cfg.FileExtensions) {
			return fmt.Errorf("not a video file: %s", path)
		}
		return processFile(path, cfg, metadataProvider)
	}

	// If it's a directory, walk through it
	var count int
	err = filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %q: %v\n", filePath, err)
			return nil
		}

		// Skip directories unless recursive mode is enabled
		if info.IsDir() {
			if filePath != path && !cfg.Recursive {
				return filepath.SkipDir
			}
			return nil
		}

		// Check if it's a video file
		if media.IsVideoFile(filePath, cfg.FileExtensions) {
			count++
			fmt.Printf("\nFound video file (%d): %s\n", count, filePath)
			if err := processFile(filePath, cfg, metadataProvider); err != nil {
				fmt.Printf("Error processing %s: %v\n", filePath, err)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking directory: %v", err)
	}

	if count == 0 {
		return fmt.Errorf("no video files found in: %s", path)
	}

	return nil
}

func generateFilename(originalPath string, info *media.VideoInfo, metadata *metadata.MovieMetadata, cfg *config.Config) string {
	dir := filepath.Dir(originalPath)
	ext := filepath.Ext(originalPath)

	// Get video stream info
	var res, codec string
	for _, stream := range info.Streams {
		if stream.CodecType == "video" {
			res = resolution.GetStandardResolution(stream.Width, stream.Height)
			codec = stream.CodecName
			break
		}
	}

	// Format the filename according to the pattern
	pattern := cfg.MovieFormat
	title := metadata.Title
	year := fmt.Sprintf("%d", metadata.Year)

	// Replace template variables
	newName := strings.NewReplacer(
		"{title}", title,
		"{year}", year,
		"{resolution}", res,
		"{codec}", codec,
	).Replace(pattern)

	// Apply scene style if enabled (use dots)
	if cfg.SceneStyle {
		cfg.Separator = "."
	}

	// Replace spaces with configured separator
	if cfg.Separator != " " {
		newName = strings.ReplaceAll(newName, " ", cfg.Separator)
	}

	// Apply lowercase if enabled
	if cfg.LowerCase {
		newName = strings.ToLower(newName)
	}

	// Clean up filename (remove invalid characters)
	newName = strings.Map(func(r rune) rune {
		switch r {
		case '/', '\\', ':', '*', '?', '"', '<', '>', '|':
			return '-'
		default:
			return r
		}
	}, newName)

	return filepath.Join(dir, newName+ext)
}

func confirmRename() bool {
	fmt.Print("\nDo you want to rename the file? (y/N): ")
	var response string
	fmt.Scanln(&response)
	return strings.ToLower(response) == "y"
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func main() {
	// Parse command line flags
	batchMode := flag.Bool("b", false, "process automatically without interactive prompts")
	recursive := flag.Bool("r", false, "search for files within nested directories")
	lowerCase := flag.Bool("l", false, "rename files using lowercase characters")
	sceneStyle := flag.Bool("s", false, "use dots in place of spaces")
	noOverwrite := flag.Bool("no-overwrite", false, "prevent relocation if it would overwrite a file")
	language := flag.String("lang", "en", "metadata language (ISO 639-1 code)")
	movieFormat := flag.String("movie-format", "{title} ({year}) [{resolution} {codec}]", "movie filename format")
	tvFormat := flag.String("tv-format", "{series} S{season:02d}E{episode:02d} [{resolution} {codec}]", "TV episode filename format")
	previewMode := flag.Bool("preview", false, "Preview mode: show what would be done without making changes")
	noMetadata := flag.Bool("no-metadata", false, "Skip online metadata lookup")
	separator := flag.String("separator", "", "Character to use as separator in filenames (default: space, use '.' for scene style)")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: vidkit [options] <file_or_directory>")
		fmt.Println("\nOptions:")
		fmt.Println("  -b, --batch      process automatically without interactive prompts")
		fmt.Println("  -r, --recursive  search for files within nested directories")
		fmt.Println("  -l, --lower      rename files using lowercase characters")
		fmt.Println("  -s, --scene      use dots in place of spaces (shortcut for --separator '.')")
		fmt.Println("  --separator      character to use as separator in filenames")
		fmt.Println("  --no-overwrite   prevent relocation if it would overwrite a file")
		fmt.Println("  --lang <code>    metadata language (ISO 639-1 code)")
		fmt.Println("  --movie-format <format> movie filename format")
		fmt.Println("  --tv-format <format>   TV episode filename format")
		fmt.Println("  --preview        show what would be done without making changes")
		fmt.Println("  --no-metadata    skip online metadata lookup")
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Override config with command line flags
	cfg.BatchMode = *batchMode
	cfg.Recursive = *recursive
	cfg.LowerCase = *lowerCase
	cfg.SceneStyle = *sceneStyle
	cfg.NoOverwrite = *noOverwrite
	cfg.Language = *language
	cfg.MovieFormat = *movieFormat
	cfg.TVFormat = *tvFormat
	cfg.PreviewMode = *previewMode
	cfg.NoMetadata = *noMetadata
	if *separator != "" {
		cfg.Separator = *separator
	}

	// Validate configuration
	if err := config.ValidateConfig(cfg); err != nil {
		fmt.Printf("Configuration error: %v\n", err)
		os.Exit(1)
	}

	// Create metadata provider if needed
	var metadataProvider *metadata.TMDbProvider
	if !cfg.NoMetadata {
		metadataProvider, err = metadata.NewTMDbProvider(cfg.TMDbAPIKey)
		if err != nil {
			fmt.Printf("Error initializing TMDb: %v\n", err)
			os.Exit(1)
		}
	}

	// Process each path argument
	for _, path := range args {
		if err := processPath(path, cfg, metadataProvider); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}
}
