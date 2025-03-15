package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	// Check default values
	if cfg.Language != "en" {
		t.Errorf("Default language = %v, want %v", cfg.Language, "en")
	}

	if cfg.MovieFormat != "{title} ({year}) [{resolution} {codec}]" {
		t.Errorf("Default movie format = %v, want %v", cfg.MovieFormat, "{title} ({year}) [{resolution} {codec}]")
	}

	if cfg.Separator != " " {
		t.Errorf("Default separator = %v, want space", cfg.Separator)
	}

	// Check default providers
	if cfg.MovieProvider != ProviderTMDb {
		t.Errorf("Default movie provider = %v, want %v", cfg.MovieProvider, ProviderTMDb)
	}

	if cfg.TVProvider != ProviderTVMaze {
		t.Errorf("Default TV provider = %v, want %v", cfg.TVProvider, ProviderTVMaze)
	}

	// Check file extensions - updated to match the current defaults
	expectedExts := []string{
		".mp4", ".mkv", ".avi", ".mov", ".wmv", ".m4v", ".mpg", ".mpeg", ".webm", ".flv", ".ts", ".m2ts", ".mts", ".mxf",
	}
	if !reflect.DeepEqual(cfg.FileExtensions, expectedExts) {
		t.Errorf("Default file extensions = %v, want %v", cfg.FileExtensions, expectedExts)
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name      string
		config    *Config
		wantError bool
	}{
		{
			name: "Valid config",
			config: &Config{
				Separator:     " ",
				Language:      "en",
				MovieFormat:   "{title} ({year})",
				TVFormat:      "{title} S{season:02d}E{episode:02d}",
				MovieProvider: ProviderTMDb,
				TVProvider:    ProviderTVMaze,
				TMDbAPIKey:    "test_key",
			},
			wantError: false,
		},
		{
			name:      "Empty config is valid",
			config:    &Config{},
			wantError: false,
		},
		{
			name: "Config with formats but without API key",
			config: &Config{
				MovieFormat:   "{title} ({year})",
				TVFormat:      "{title} S{season:02d}E{episode:02d}",
				MovieProvider: ProviderTMDb,
				NoMetadata:    true, // No metadata lookup, so no API key needed
			},
			wantError: false,
		},
		{
			name: "TMDb provider without API key",
			config: &Config{
				MovieFormat:    "{title} ({year})",
				TVFormat:       "{title} S{season:02d}E{episode:02d}",
				MovieProvider:  ProviderTMDb,
				EnableMetadata: true,
				NoMetadata:     false,
			},
			wantError: true,
		},
		{
			name: "OMDb provider without API key",
			config: &Config{
				MovieFormat:    "{title} ({year})",
				TVFormat:       "{title} S{season:02d}E{episode:02d}",
				MovieProvider:  ProviderOMDb,
				EnableMetadata: true,
				NoMetadata:     false,
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.config)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateConfig() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestLoadAndSaveConfig(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "vidkit-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Set the config path to our test directory
	originalPath := ConfigFilePath
	defer SetConfigPath(originalPath)

	configPath := filepath.Join(tmpDir, "config.json")
	SetConfigPath(func() string {
		return configPath
	})

	// Test saving a config
	testConfig := &Config{
		TMDbAPIKey:     "test_key",
		OMDbAPIKey:     "test_omdb_key",
		TVDbAPIKey:     "test_tvdb_key",
		Language:       "es",
		Separator:      ".",
		MovieFormat:    "Custom {title} ({year})",
		TVFormat:       "Custom {title} S{season:02d}E{episode:02d}",
		FileExtensions: []string{".mp4", ".mkv"},
		MovieProvider:  ProviderOMDb,
		TVProvider:     ProviderTVDb,
	}

	err = SaveConfig(testConfig)
	if err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("Config file was not created at %s", configPath)
	}

	// Test loading the config
	loadedConfig, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// Check loaded values match what we saved
	if loadedConfig.TMDbAPIKey != testConfig.TMDbAPIKey {
		t.Errorf("Loaded API key = %v, want %v", loadedConfig.TMDbAPIKey, testConfig.TMDbAPIKey)
	}

	if loadedConfig.OMDbAPIKey != testConfig.OMDbAPIKey {
		t.Errorf("Loaded OMDb API key = %v, want %v", loadedConfig.OMDbAPIKey, testConfig.OMDbAPIKey)
	}

	if loadedConfig.TVDbAPIKey != testConfig.TVDbAPIKey {
		t.Errorf("Loaded TVDb API key = %v, want %v", loadedConfig.TVDbAPIKey, testConfig.TVDbAPIKey)
	}

	if loadedConfig.Language != testConfig.Language {
		t.Errorf("Loaded language = %v, want %v", loadedConfig.Language, testConfig.Language)
	}

	if loadedConfig.MovieFormat != testConfig.MovieFormat {
		t.Errorf("Loaded movie format = %v, want %v", loadedConfig.MovieFormat, testConfig.MovieFormat)
	}

	if loadedConfig.MovieProvider != testConfig.MovieProvider {
		t.Errorf("Loaded movie provider = %v, want %v", loadedConfig.MovieProvider, testConfig.MovieProvider)
	}

	if loadedConfig.TVProvider != testConfig.TVProvider {
		t.Errorf("Loaded TV provider = %v, want %v", loadedConfig.TVProvider, testConfig.TVProvider)
	}
}

func TestConfigPermissions(t *testing.T) {
	// Skip this test on platforms where file permissions behave differently
	if os.Getenv("SKIP_PERMISSION_TEST") != "" {
		t.Skip("Skipping permission test")
	}

	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "vidkit-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Set the config path to our test directory
	originalPath := ConfigFilePath
	defer SetConfigPath(originalPath)

	configDir := filepath.Join(tmpDir, ".config", "vidkit")
	configPath := filepath.Join(configDir, "config.json")
	SetConfigPath(func() string {
		return configPath
	})

	// Ensure directory doesn't exist yet
	if _, err := os.Stat(configDir); !os.IsNotExist(err) {
		t.Fatalf("Test directory already exists: %v", configDir)
	}

	// Test creating config in new directory
	cfg := DefaultConfig()
	err = SaveConfig(cfg)
	if err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	// Check directory was created
	dirInfo, err := os.Stat(configDir)
	if err != nil {
		t.Fatalf("Failed to stat config dir: %v", err)
	}
	if !dirInfo.IsDir() {
		t.Errorf("%s is not a directory", configDir)
	}

	// Check file permissions
	fileInfo, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("Failed to stat config file: %v", err)
	}

	// File should be readable and writable by user
	expectedMode := os.FileMode(0644)
	if fileInfo.Mode().Perm() != expectedMode {
		t.Errorf("Config file permissions = %v, want %v", fileInfo.Mode().Perm(), expectedMode)
	}
}
