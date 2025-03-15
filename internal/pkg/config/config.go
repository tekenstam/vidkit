package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds application configuration
type Config struct {
	// Common options
	TMDbAPIKey     string `json:"tmdb_api_key"`
	BatchMode      bool   `json:"batch_mode"`
	Recursive      bool   `json:"recursive"`
	LowerCase      bool   `json:"lowercase"`
	SceneStyle     bool   `json:"scene_style"`
	Separator      string `json:"separator"`
	FileExtensions []string `json:"file_extensions"`
	Language       string `json:"language"`
	NoOverwrite    bool   `json:"no_overwrite"`
	NoMetadata     bool   `json:"no_metadata"`
	PreviewMode    bool   `json:"preview_mode"`
	EnableMetadata bool   `json:"enable_metadata"`

	// Formatting options
	MovieFormat string `json:"movie_format"`
	TVFormat    string `json:"tv_format"`
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
		BatchMode:      false,
		Recursive:      false,
		LowerCase:      false,
		SceneStyle:     false,
		Separator:      " ",
		FileExtensions: []string{".mp4", ".mkv", ".avi", ".mov", ".wmv", ".m4v", ".mpg", ".mpeg"},
		Language:       "en",
		NoOverwrite:    false,
		NoMetadata:     false,
		PreviewMode:    false,
		EnableMetadata: true,
		MovieFormat:    "{title} ({year}) [{resolution} {codec}]",
		TVFormat:       "{title} S{season:02d}E{episode:02d} {episode_title} [{resolution} {codec}]",
	}

	// Try to read config file
	configPath := ConfigFilePath()
	data, err := os.ReadFile(configPath)
	
	// If file doesn't exist, create default config
	if errors.Is(err, os.ErrNotExist) {
		// Ensure directory exists
		err = os.MkdirAll(filepath.Dir(configPath), 0755)
		if err != nil {
			return cfg, err
		}
		
		// Write default config
		data, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			return cfg, err
		}
		
		err = os.WriteFile(configPath, data, 0644)
		if err != nil {
			return cfg, err
		}
		
		return cfg, nil
	}
	
	// If other error occurred while reading
	if err != nil {
		return cfg, err
	}
	
	// Parse config file
	err = json.Unmarshal(data, cfg)
	if err != nil {
		return cfg, err
	}
	
	return cfg, nil
}

// ValidateConfig validates the configuration
func ValidateConfig(cfg *Config) error {
	if cfg.Separator == "" {
		cfg.Separator = " "
	}
	
	if cfg.SceneStyle {
		cfg.Separator = "."
	}
	
	if cfg.Language == "" {
		cfg.Language = "en"
	}
	
	if cfg.MovieFormat == "" {
		cfg.MovieFormat = "{title} ({year}) [{resolution} {codec}]"
	}
	
	if cfg.TVFormat == "" {
		cfg.TVFormat = "{title} S{season:02d}E{episode:02d} {episode_title} [{resolution} {codec}]"
	}
	
	return nil
}

// SaveConfig saves the configuration to the default location
func SaveConfig(config *Config) error {
	configPath := ConfigFilePath()
	
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}
	
	// Write config
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}
	
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %v", err)
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
		Language: "en",
		MovieFormat: "{title} ({year}) [{resolution} {codec}]",
		TVFormat: "{title} S{season:02d}E{episode:02d} {episode_title} [{resolution} {codec}]",
		FileExtensions: []string{
			".mp4", ".mkv", ".avi", ".mov", ".wmv",
			".flv", ".webm", ".m4v", ".mpg", ".mpeg", ".3gp",
		},
		Separator: " ", // Default to spaces
	}
}
