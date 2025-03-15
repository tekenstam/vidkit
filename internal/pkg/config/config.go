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

// Config holds application configuration
type Config struct {
	// Common options
	TMDbAPIKey     string   `json:"tmdb_api_key"`
	OMDbAPIKey     string   `json:"omdb_api_key"`
	TVDbAPIKey     string   `json:"tvdb_api_key"`
	BatchMode      bool     `json:"batch_mode"`
	Recursive      bool     `json:"recursive"`
	LowerCase      bool     `json:"lowercase"`
	SceneStyle     bool     `json:"scene_style"`
	Separator      string   `json:"separator"`
	FileExtensions []string `json:"file_extensions"`
	Language       string   `json:"language"`
	NoOverwrite    bool     `json:"no_overwrite"`
	NoMetadata     bool     `json:"no_metadata"`
	PreviewMode    bool     `json:"preview_mode"`
	EnableMetadata bool     `json:"enable_metadata"`

	// Provider selection
	MovieProvider ProviderType `json:"movie_provider"`
	TVProvider    ProviderType `json:"tv_provider"`

	// Formatting options
	MovieFormat string `json:"movie_format"`
	TVFormat    string `json:"tv_format"`

	// Directory organization options
	MovieDirectory string `json:"movie_directory"` // Template for movie directory structure
	TVDirectory    string `json:"tv_directory"`    // Template for TV show directory structure
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
