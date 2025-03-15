package metadata

import (
	"fmt"

	"github.com/tekenstam/vidkit/internal/pkg/config"
)

// CreateMovieProvider creates the appropriate movie metadata provider based on configuration
func CreateMovieProvider(cfg *config.Config) (MetadataProvider, error) {
	switch cfg.MovieProvider {
	case config.ProviderTMDb:
		return NewTMDbProvider(cfg.TMDbAPIKey)
	case config.ProviderOMDb:
		return NewOMDbProvider(cfg.OMDbAPIKey)
	default:
		return nil, fmt.Errorf("unsupported movie provider: %s", cfg.MovieProvider)
	}
}

// CreateTVShowProvider creates the appropriate TV show metadata provider based on configuration
func CreateTVShowProvider(cfg *config.Config) (MetadataProvider, error) {
	switch cfg.TVProvider {
	case config.ProviderTVMaze:
		return NewTvMazeProvider(), nil
	case config.ProviderTVDb:
		return NewTVDbProvider(cfg.TVDbAPIKey)
	default:
		return nil, fmt.Errorf("unsupported TV show provider: %s", cfg.TVProvider)
	}
}

// GetProvider returns the appropriate provider for the type of content
func GetProvider(cfg *config.Config, isTV bool) (MetadataProvider, error) {
	if isTV {
		return CreateTVShowProvider(cfg)
	}
	return CreateMovieProvider(cfg)
}
