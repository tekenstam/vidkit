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

	// Check file extensions
	expectedExts := []string{
		".mp4", ".mkv", ".avi", ".mov", ".wmv",
		".flv", ".webm", ".m4v", ".mpg", ".mpeg", ".3gp",
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
			name: "Valid config with API key",
			config: &Config{
				TMDbAPIKey: "valid_key",
				NoMetadata: false,
			},
			wantError: false,
		},
		{
			name: "Missing API key but metadata disabled",
			config: &Config{
				TMDbAPIKey: "",
				NoMetadata: true,
			},
			wantError: false,
		},
		{
			name: "Missing API key with metadata enabled",
			config: &Config{
				TMDbAPIKey: "",
				NoMetadata: false,
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

	// Create a test config
	testConfig := &Config{
		TMDbAPIKey:  "test_key",
		Language:    "fr",
		MovieFormat: "test-{title}-{year}",
		Separator:   "-",
		BatchMode:   true,
	}

	// Set up a temporary config path
	configPath := filepath.Join(tmpDir, "config.json")
	originalPath := getConfigPath
	SetConfigPath(func() string {
		return configPath
	})
	defer SetConfigPath(originalPath)

	// Test saving config
	if err := SaveConfig(testConfig); err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	// Test loading config
	loadedConfig, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// Compare saved and loaded configs
	if !reflect.DeepEqual(testConfig, loadedConfig) {
		t.Errorf("LoadConfig() = %v, want %v", loadedConfig, testConfig)
	}

	// Test loading with missing file (should create default)
	os.Remove(configPath)
	defaultConfig, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() with missing file error = %v", err)
	}

	// Verify default values
	if defaultConfig.Language != "en" {
		t.Errorf("Default config language = %v, want en", defaultConfig.Language)
	}
}

func TestConfigPermissions(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "vidkit-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Set up a temporary config path
	configPath := filepath.Join(tmpDir, "config.json")
	originalPath := getConfigPath
	SetConfigPath(func() string {
		return configPath
	})
	defer SetConfigPath(originalPath)

	// Create a test config
	testConfig := DefaultConfig()
	testConfig.TMDbAPIKey = "test_key"

	// Save config
	if err := SaveConfig(testConfig); err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	// Check file permissions
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("Failed to stat config file: %v", err)
	}

	// Config file should be readable/writable only by owner
	expectedPerm := os.FileMode(0644)
	if info.Mode().Perm() != expectedPerm {
		t.Errorf("Config file permissions = %v, want %v", info.Mode().Perm(), expectedPerm)
	}
}
