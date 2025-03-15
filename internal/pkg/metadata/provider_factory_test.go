package metadata

import (
	"fmt"
	"testing"

	"github.com/tekenstam/vidkit/internal/pkg/config"
)

func TestCreateMovieProvider(t *testing.T) {
	tests := []struct {
		name          string
		providerType  config.ProviderType
		apiKey        string
		expectError   bool
		expectedType  string
	}{
		{
			name:         "TMDb Provider",
			providerType: config.ProviderTMDb,
			apiKey:       "test_tmdb_key",
			expectError:  false,
			expectedType: "*metadata.TMDbProvider",
		},
		{
			name:         "OMDb Provider",
			providerType: config.ProviderOMDb,
			apiKey:       "test_omdb_key",
			expectError:  false,
			expectedType: "*metadata.OMDbProvider",
		},
		{
			name:         "TMDb Provider with empty key",
			providerType: config.ProviderTMDb,
			apiKey:       "",
			expectError:  true,
			expectedType: "",
		},
		{
			name:         "OMDb Provider with empty key",
			providerType: config.ProviderOMDb,
			apiKey:       "",
			expectError:  true,
			expectedType: "",
		},
		{
			name:         "Unknown Provider",
			providerType: "unknown",
			apiKey:       "test_key",
			expectError:  true,
			expectedType: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create config with appropriate provider and key
			cfg := &config.Config{
				MovieProvider: tt.providerType,
			}
			
			// Set the appropriate API key based on provider type
			switch tt.providerType {
			case config.ProviderTMDb:
				cfg.TMDbAPIKey = tt.apiKey
			case config.ProviderOMDb:
				cfg.OMDbAPIKey = tt.apiKey
			}
			
			// Call the factory function
			provider, err := CreateMovieProvider(cfg)
			
			// Check if we got expected error status
			if (err != nil) != tt.expectError {
				t.Errorf("CreateMovieProvider() error = %v, expectError %v", err, tt.expectError)
				return
			}
			
			// If we expect an error, no need to check the provider type
			if tt.expectError {
				return
			}
			
			// Check if we got the expected provider type
			providerType := fmt.Sprintf("%T", provider)
			if providerType != tt.expectedType {
				t.Errorf("CreateMovieProvider() provider type = %v, want %v", providerType, tt.expectedType)
			}
		})
	}
}

func TestCreateTVShowProvider(t *testing.T) {
	tests := []struct {
		name          string
		providerType  config.ProviderType
		apiKey        string
		expectError   bool
		expectedType  string
	}{
		{
			name:         "TvMaze Provider",
			providerType: config.ProviderTVMaze,
			apiKey:       "", // No API key needed for TvMaze
			expectError:  false,
			expectedType: "*metadata.TvMazeProvider",
		},
		{
			name:         "TVDb Provider",
			providerType: config.ProviderTVDb,
			apiKey:       "test_tvdb_key",
			expectError:  false,
			expectedType: "*metadata.TVDbProvider",
		},
		{
			name:         "TVDb Provider with empty key",
			providerType: config.ProviderTVDb,
			apiKey:       "",
			expectError:  true,
			expectedType: "",
		},
		{
			name:         "Unknown Provider",
			providerType: "unknown",
			apiKey:       "test_key",
			expectError:  true,
			expectedType: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create config with appropriate provider and key
			cfg := &config.Config{
				TVProvider: tt.providerType,
			}
			
			// Set the appropriate API key based on provider type
			if tt.providerType == config.ProviderTVDb {
				cfg.TVDbAPIKey = tt.apiKey
			}
			
			// Call the factory function
			provider, err := CreateTVShowProvider(cfg)
			
			// Check if we got expected error status
			if (err != nil) != tt.expectError {
				t.Errorf("CreateTVShowProvider() error = %v, expectError %v", err, tt.expectError)
				return
			}
			
			// If we expect an error, no need to check the provider type
			if tt.expectError {
				return
			}
			
			// Check if we got the expected provider type
			providerType := fmt.Sprintf("%T", provider)
			if providerType != tt.expectedType {
				t.Errorf("CreateTVShowProvider() provider type = %v, want %v", providerType, tt.expectedType)
			}
		})
	}
}

func TestGetProvider(t *testing.T) {
	tests := []struct {
		name         string
		isTV         bool
		movieType    config.ProviderType
		tvType       config.ProviderType
		expectError  bool
		expectedType string
	}{
		{
			name:         "Get Movie Provider",
			isTV:         false,
			movieType:    config.ProviderTMDb,
			tvType:       config.ProviderTVMaze,
			expectError:  false,
			expectedType: "*metadata.TMDbProvider",
		},
		{
			name:         "Get TV Provider",
			isTV:         true,
			movieType:    config.ProviderTMDb,
			tvType:       config.ProviderTVMaze,
			expectError:  false,
			expectedType: "*metadata.TvMazeProvider",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create config with test API keys
			cfg := &config.Config{
				MovieProvider: tt.movieType,
				TVProvider:    tt.tvType,
				TMDbAPIKey:    "test_tmdb_key",
				OMDbAPIKey:    "test_omdb_key",
				TVDbAPIKey:    "test_tvdb_key",
			}
			
			// Call the function
			provider, err := GetProvider(cfg, tt.isTV)
			
			// Check if we got expected error status
			if (err != nil) != tt.expectError {
				t.Errorf("GetProvider() error = %v, expectError %v", err, tt.expectError)
				return
			}
			
			// If we expect an error, no need to check the provider type
			if tt.expectError {
				return
			}
			
			// Check if we got the expected provider type
			providerType := fmt.Sprintf("%T", provider)
			if providerType != tt.expectedType {
				t.Errorf("GetProvider() provider type = %v, want %v", providerType, tt.expectedType)
			}
		})
	}
}
