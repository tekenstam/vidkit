// Package config provides configuration management for VidKit's application settings.
// It handles loading and saving user configurations, API keys for metadata providers,
// and operational settings that control VidKit's behavior.
//
// The configuration system supports multiple metadata providers for both movies and TV shows,
// allowing users to select their preferred data source. The package manages API keys,
// application preferences, and operational modes in a single configuration structure.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// ProviderType defines supported metadata provider types.
// These represent the available online services that can provide
// metadata for movies and TV shows for VidKit's analysis.
type ProviderType string

const (
	// Movie provider types
	ProviderTMDb ProviderType = "tmdb" // The Movie Database (primary movie provider)
	ProviderOMDb ProviderType = "omdb" // Open Movie Database (alternative movie provider)

	// TV show provider types
	ProviderTVMaze ProviderType = "tvmaze" // TVMaze (primary TV show provider)
	ProviderTVDb   ProviderType = "tvdb"   // The TV Database (alternative TV show provider)
)

// Config holds application configuration settings for VidKit.
// This structure is serialized to/from JSON when saving/loading configurations.
// It contains all user preferences and API keys needed for metadata lookups.
type Config struct {
	// API keys for metadata providers
	TMDbAPIKey     string   `json:"tmdb_api_key"` // API key for The Movie Database
	OMDbAPIKey     string   `json:"omdb_api_key"` // API key for Open Movie Database
	TVDbAPIKey     string   `json:"tvdb_api_key"` // API key for The TV Database
	BatchMode      bool     `json:"batch_mode"` // Run without interactive prompts
	Recursive      bool     `json:"recursive"` // Process directories recursively
	LowerCase      bool     `json:"lowercase"` // Convert filenames to lowercase
	SceneStyle     bool     `json:"scene_style"` // Use dots instead of spaces in filenames
	Separator      string   `json:"separator"`
	FileExtensions []string `json:"file_extensions"`
	Language       string   `json:"language"` // Preferred language for metadata (ISO 639-1 code)
	NoOverwrite    bool     `json:"no_overwrite"`
	NoMetadata     bool     `json:"no_metadata"`
	PreviewMode    bool     `json:"preview_mode"`
	EnableMetadata bool     `json:"enable_metadata"`

	// Provider preferences
	MovieProvider ProviderType `json:"movie_provider"` // Preferred movie metadata provider
	TVProvider    ProviderType `json:"tv_provider"`    // Preferred TV show metadata provider

	// Format templates
	MovieFormat string `json:"movie_format"` // Template for movie filename format
	TVFormat    string `json:"tv_format"`    // Template for TV show filename format
	MovieDirectory string `json:"movie_directory"` // Template for movie directory organization
	TVDirectory    string `json:"tv_directory"`    // Template for TV show directory organization
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
		LowerCase:      false,
		SceneStyle:     false,
		Separator:      " ",
		FileExtensions: []string{".mp4", ".mkv", ".avi", ".mov", ".wmv", ".m4v", ".mpg", ".mpeg"},
		Language:       "en",
		NoOverwrite:    false,
		NoMetadata:     false,
		MovieProvider:  ProviderTMDb,
		TVProvider:     ProviderTVMaze,
		MovieFormat:    "{title} ({year}) [{resolution} {codec}]",
		TVFormat:       "{title} S{season:02d}E{episode:02d} {episode_title} [{resolution} {codec}]",
		EnableMetadata: true,
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
	if cfg.MovieFormat == "" {
		cfg.MovieFormat = "{title} ({year}) [{resolution} {codec}]"
	}

	// Apply default TV format if needed
	if cfg.TVFormat == "" {
		cfg.TVFormat = "{title} S{season:02d}E{episode:02d} {episode_title} [{resolution} {codec}]"
	}

	// Check if metadata is enabled but no API key is provided
	if !cfg.NoMetadata && cfg.EnableMetadata {
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

// DefaultConfig returns a configuration with default values
func DefaultConfig() *Config {
	return &Config{
		TMDbAPIKey:     "",
		OMDbAPIKey:     "",
		TVDbAPIKey:     "",
		BatchMode:      false,
		Recursive:      false,
		LowerCase:      false,
		SceneStyle:     false,
		Separator:      " ",
		FileExtensions: []string{".mp4", ".mkv", ".avi", ".mov", ".wmv", ".m4v", ".mpg", ".mpeg", ".webm", ".flv", ".ts", ".m2ts", ".mts", ".mxf"},
		Language:       "en",
		NoOverwrite:    true,
		NoMetadata:     false,
		MovieProvider:  ProviderTMDb,
		TVProvider:     ProviderTVMaze,
		MovieFormat:    "{title} ({year}) [{resolution} {codec}]",
		TVFormat:       "{title} S{season:02d}E{episode:02d} {episode_title} [{resolution} {codec}]",
		EnableMetadata: true,
		MovieDirectory: "{genre}/{title} ({year})",
		TVDirectory:    "{genre}/{title}/Season {season}",
		OrganizeFiles:  true,
	}
}
