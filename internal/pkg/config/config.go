package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// ProviderType defines supported metadata provider types
type ProviderType string

const (
	// Movie provider types
	ProviderTMDb ProviderType = "tmdb"
	ProviderOMDb ProviderType = "omdb"

	// TV show provider types
	ProviderTVMaze ProviderType = "tvmaze"
	ProviderTVDb   ProviderType = "tvdb"
)

// Config holds application configuration settings for VidKit.
// This structure is serialized to/from JSON when saving/loading configurations.
// It contains all user preferences and API keys needed for metadata lookups.
type Config struct {
	// API keys for metadata providers
	TMDbAPIKey string `json:"tmdb_api_key"` // API key for The Movie Database
	OMDbAPIKey string `json:"omdb_api_key"` // API key for Open Movie Database
	TVDbAPIKey string `json:"tvdb_api_key"` // API key for The TV Database

	// Operational modes
	BatchMode bool `json:"batch_mode"` // Run without interactive prompts
	Recursive bool `json:"recursive"` // Process directories recursively

	// File handling preferences
	Lowercase  bool `json:"lowercase"`  // Convert filenames to lowercase
	SceneStyle bool `json:"scene_style"` // Use dots instead of spaces in filenames
	Separator      string   `json:"separator"`
	FileExtensions []string `json:"file_extensions"`
	Language       string   `json:"language"` // Preferred language for metadata (ISO 639-1 code)
	NoOverwrite    bool     `json:"no_overwrite"`
	NoMetadata     bool     `json:"no_metadata"`
	PreviewMode    bool     `json:"preview_mode"`
	OnlyVideo      bool     `json:"only_video"`

	// Provider preferences
	MovieProvider ProviderType `json:"movie_provider"` // Preferred movie metadata provider
	TVProvider    ProviderType `json:"tv_provider"`    // Preferred TV show metadata provider

	// Filename and directory templates
	MovieFilenameTemplate string `json:"movie_filename_template"` // Template for movie filename 
	TVFilenameTemplate    string `json:"tv_filename_template"`    // Template for TV show filename
	MovieDirectoryTemplate string `json:"movie_directory_template"` // Template for movie directory organization
	TVDirectoryTemplate    string `json:"tv_directory_template"`    // Template for TV show directory organization
	OrganizeFiles  bool   `json:"organize_files"` // Whether to move files to organized directories
}

// ConfigFilePath returns the path to the config file
var ConfigFilePath = func() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".config", "vidkit", "config.json")
}

// LoadConfig loads configuration from file
func LoadConfig() (*Config, error) {
	// Create default config
	cfg := &Config{
		TMDbAPIKey:     "",
		OMDbAPIKey:     "",
		TVDbAPIKey:     "",
		BatchMode:      false,
		Recursive:      false,
		Lowercase:      false,
		SceneStyle:     false,
		Separator:      " ",
		FileExtensions: []string{".mp4", ".mkv", ".avi", ".mov", ".wmv", ".m4v", ".mpg", ".mpeg", ".webm", ".flv", ".ts", ".m2ts", ".mts", ".mxf"},
		Language:       "en",
		NoOverwrite:    true,
		NoMetadata:     false,
		MovieProvider:  ProviderTMDb,
		TVProvider:     ProviderTVMaze,
		MovieFilenameTemplate: "{title} ({year}) [{resolution} {codec}]",
		TVFilenameTemplate:    "{title} S{season:02d}E{episode:02d} {episode_title} [{resolution} {codec}]",
		MovieDirectoryTemplate: "{genre}/{title} ({year})",
		TVDirectoryTemplate:    "{genre}/{title}/Season {season}",
		OrganizeFiles: false,
	}

	// Check if config file exists
	configPath := ConfigFilePath()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create directory structure
		configDir := filepath.Dir(configPath)
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return nil, err
		}

		// Save default config
		defaultConfig := DefaultConfig()
		if err := SaveConfig(defaultConfig); err != nil {
			return nil, err
		}
		return defaultConfig, nil
	}

	// Read config file
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// Parse config JSON
	if err := json.Unmarshal(configData, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// ValidateConfig validates the configuration
func ValidateConfig(cfg *Config) error {
	// Apply default values if needed
	if cfg.MovieFilenameTemplate == "" {
		cfg.MovieFilenameTemplate = "{title} ({year}) [{resolution} {codec}]"
	}

	// Apply default TV format if needed
	if cfg.TVFilenameTemplate == "" {
		cfg.TVFilenameTemplate = "{title} S{season:02d}E{episode:02d} {episode_title} [{resolution} {codec}]"
	}

	// Check if metadata is enabled but no API key is provided
	if !cfg.NoMetadata {
		// Check based on selected providers
		if cfg.MovieProvider == ProviderTMDb && cfg.TMDbAPIKey == "" {
			return errors.New("TMDb API key is required for metadata lookup (set tmdb_api_key in config.json)")
		}
		if cfg.MovieProvider == ProviderOMDb && cfg.OMDbAPIKey == "" {
			return errors.New("OMDb API key is required for metadata lookup (set omdb_api_key in config.json)")
		}
		if cfg.TVProvider == ProviderTVDb && cfg.TVDbAPIKey == "" {
			return errors.New("TVDb API key is required for metadata lookup (set tvdb_api_key in config.json)")
		}
	}

	// Apply scene style settings
	if cfg.SceneStyle && cfg.Separator == " " {
		cfg.Separator = "."
	}

	return nil
}

// SaveConfig saves the configuration to the default location
func SaveConfig(config *Config) error {
	configPath := ConfigFilePath()

	// Create directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("error creating config directory: %v", err)
	}

	// Marshal config to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling config: %v", err)
	}

	// Write config to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("error writing config file: %v", err)
	}

	return nil
}

// SetConfigPath allows setting a custom config path function (for testing)
func SetConfigPath(fn func() string) {
	ConfigFilePath = fn
}

// DefaultConfig returns a new Config instance with sensible defaults.
// This is used when no user configuration exists or when creating a new
// configuration file for the first time.
func DefaultConfig() *Config {
	return &Config{
		// Default API keys are empty and must be provided by the user
		TMDbAPIKey: "",
		OMDbAPIKey: "",
		TVDbAPIKey: "",

		// Default operational modes
		BatchMode: false,
		Recursive: false,

		// Default file handling preferences
		Lowercase:  false,
		SceneStyle: false,
		Separator:  " ",
		Language:   "en",
		
		NoOverwrite: false,
		NoMetadata:  false,
		PreviewMode: false,
		OnlyVideo:   true,

		// Default providers
		MovieProvider: ProviderTMDb,
		TVProvider:    ProviderTVMaze,

		// Default filename templates
		MovieFilenameTemplate: "{title} ({year}) [{resolution} {codec}]",
		TVFilenameTemplate:    "{title} - S{season:02d}E{episode:02d} - {episode_title}",
		
		// Default directory templates
		MovieDirectoryTemplate: "{genre}/{title} ({year})",
		TVDirectoryTemplate:    "{genre}/{title}/Season {season}",
		OrganizeFiles: false,

		// Default file extensions to process
		FileExtensions: []string{".mp4", ".mkv", ".avi", ".mov", ".m4v"},
	}
}
