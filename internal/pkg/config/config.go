package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the application configuration
type Config struct {
	TMDbAPIKey      string   `json:"tmdb_api_key"`
	BatchMode       bool     `json:"batch_mode"`
	Recursive       bool     `json:"recursive"`
	LowerCase       bool     `json:"lower_case"`
	SceneStyle      bool     `json:"scene_style"`
	NoOverwrite     bool     `json:"no_overwrite"`
	Language        string   `json:"language"`
	MovieFormat     string   `json:"movie_format"`
	TVFormat        string   `json:"tv_format"`
	IgnorePatterns  []string `json:"ignore_patterns"`
	FileExtensions  []string `json:"file_extensions"`
	PreviewMode     bool     `json:"preview_mode"`
	NoMetadata      bool     `json:"no_metadata"`
	Separator       string   `json:"separator"`
}

// DefaultConfig returns a configuration with default values
func DefaultConfig() *Config {
	return &Config{
		Language: "en",
		MovieFormat: "{title} ({year}) [{resolution} {codec}]",
		TVFormat: "{series} S{season:02d}E{episode:02d} [{resolution} {codec}]",
		FileExtensions: []string{
			".mp4", ".mkv", ".avi", ".mov", ".wmv",
			".flv", ".webm", ".m4v", ".mpg", ".mpeg", ".3gp",
		},
		Separator: " ", // Default to spaces
	}
}

// LoadConfig loads the configuration from the default location
func LoadConfig() (*Config, error) {
	configPath := getConfigPath()
	
	// Try to load existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Create default config
			config := DefaultConfig()
			
			// Ensure directory exists
			if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
				return nil, fmt.Errorf("failed to create config directory: %v", err)
			}
			
			// Write default config
			if err := SaveConfig(config); err != nil {
				return nil, err
			}
			
			return config, nil
		}
		return nil, fmt.Errorf("failed to read config: %v", err)
	}
	
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %v", err)
	}

	// Set default separator if not specified
	if config.Separator == "" {
		config.Separator = " "
	}
	
	return &config, nil
}

// SaveConfig saves the configuration to the default location
func SaveConfig(config *Config) error {
	configPath := getConfigPath()
	
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

// ValidateConfig checks if the configuration is valid
func ValidateConfig(config *Config) error {
	if !config.NoMetadata && config.TMDbAPIKey == "" {
		return fmt.Errorf("TMDb API key is required for metadata lookup")
	}
	return nil
}

// getConfigPath returns the path to the configuration file
func getConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".vidkit.json"
	}
	return filepath.Join(homeDir, ".config", "vidkit", "config.json")
}
