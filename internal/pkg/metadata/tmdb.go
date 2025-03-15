package metadata

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	tmdb "github.com/cyruzin/golang-tmdb"
)

// MovieSearch represents a movie search request
type MovieSearch struct {
	Title string
	Year  int
}

// MovieMetadata represents movie metadata from TMDb
type MovieMetadata struct {
	Title    string
	Year     int
	Overview string
}

// TMDbProvider implements movie metadata lookup using TMDb
type TMDbProvider struct {
	client *tmdb.Client
}

// NewTMDbProvider creates a new TMDb metadata provider
func NewTMDbProvider(apiKey string) (*TMDbProvider, error) {
	client, err := tmdb.Init(apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize TMDb client: %v", err)
	}
	return &TMDbProvider{
		client: client,
	}, nil
}

// SearchMovie searches for a movie using TMDb
func (p *TMDbProvider) SearchMovie(search MovieSearch, language string) (*MovieMetadata, error) {
	options := map[string]string{
		"language": language,
	}
	
	// If we have a year, add it to improve search accuracy
	if search.Year > 0 {
		options["year"] = strconv.Itoa(search.Year)
	}

	searchResults, err := p.client.GetSearchMovies(search.Title, options)
	if err != nil {
		return nil, fmt.Errorf("failed to search movie: %v", err)
	}

	if len(searchResults.Results) == 0 {
		// If no results with year, try without year
		if search.Year > 0 {
			delete(options, "year")
			searchResults, err = p.client.GetSearchMovies(search.Title, options)
			if err != nil {
				return nil, fmt.Errorf("failed to search movie: %v", err)
			}
		}
		if len(searchResults.Results) == 0 {
			return nil, fmt.Errorf("no movies found matching '%s'", search.Title)
		}
	}

	// Get the first result's details
	movieID := int(searchResults.Results[0].ID)
	movie, err := p.client.GetMovieDetails(movieID, options)
	if err != nil {
		return nil, fmt.Errorf("failed to get movie details: %v", err)
	}

	year := 0
	if movie.ReleaseDate != "" {
		if t, err := time.Parse("2006-01-02", movie.ReleaseDate); err == nil {
			year = t.Year()
		}
	}

	return &MovieMetadata{
		Title:    movie.Title,
		Year:     year,
		Overview: movie.Overview,
	}, nil
}

// ExtractMovieInfo extracts movie information from a filename
func ExtractMovieInfo(filename string) MovieSearch {
	// Extract base name without extension
	basename := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	
	// Look for year pattern (YYYY) in filename
	yearPattern := regexp.MustCompile(`\((\d{4})\)`)
	year := 0
	title := basename

	if matches := yearPattern.FindStringSubmatch(basename); len(matches) > 1 {
		if y, err := strconv.Atoi(matches[1]); err == nil {
			year = y
			// Remove the year from the title
			title = strings.TrimSpace(yearPattern.ReplaceAllString(basename, ""))
		}
	}

	// Clean up the title by removing common patterns
	title = strings.NewReplacer(
		"1080p", "",
		"720p", "",
		"480p", "",
		"360p", "",
		"h264", "",
		"x264", "",
		"[", "",
		"]", "",
		"(", "",
		")", "",
	).Replace(title)
	
	return MovieSearch{
		Title: strings.TrimSpace(title),
		Year:  year,
	}
}
