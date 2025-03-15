package metadata

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// OMDbProvider implements movie metadata lookup using the Open Movie Database API
type OMDbProvider struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// OMDbResponse represents the response from the OMDb API
type OMDbResponse struct {
	Title    string `json:"Title"`
	Year     string `json:"Year"`
	Rated    string `json:"Rated"`
	Released string `json:"Released"`
	Runtime  string `json:"Runtime"`
	Genre    string `json:"Genre"`
	Director string `json:"Director"`
	Writer   string `json:"Writer"`
	Actors   string `json:"Actors"`
	Plot     string `json:"Plot"`
	Language string `json:"Language"`
	Country  string `json:"Country"`
	Awards   string `json:"Awards"`
	Poster   string `json:"Poster"`
	Ratings  []struct {
		Source string `json:"Source"`
		Value  string `json:"Value"`
	} `json:"Ratings"`
	Metascore  string `json:"Metascore"`
	ImdbRating string `json:"imdbRating"`
	ImdbVotes  string `json:"imdbVotes"`
	ImdbID     string `json:"imdbID"`
	Type       string `json:"Type"`
	DVD        string `json:"DVD"`
	BoxOffice  string `json:"BoxOffice"`
	Production string `json:"Production"`
	Website    string `json:"Website"`
	Response   string `json:"Response"`
	Error      string `json:"Error,omitempty"`
}

// OMDbSearchResponse represents the search response from the OMDb API
type OMDbSearchResponse struct {
	Search       []OMDbSearchResult `json:"Search"`
	TotalResults string             `json:"totalResults"`
	Response     string             `json:"Response"`
	Error        string             `json:"Error,omitempty"`
}

// OMDbSearchResult represents a search result from the OMDb API
type OMDbSearchResult struct {
	Title  string `json:"Title"`
	Year   string `json:"Year"`
	ImdbID string `json:"imdbID"`
	Type   string `json:"Type"`
	Poster string `json:"Poster"`
}

// NewOMDbProvider creates a new OMDb metadata provider
func NewOMDbProvider(apiKey string) (*OMDbProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("OMDb API key is required")
	}

	return &OMDbProvider{
		apiKey:  apiKey,
		baseURL: "http://www.omdbapi.com/",
		client:  &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Ensure OMDbProvider implements MetadataProvider
var _ MetadataProvider = (*OMDbProvider)(nil)

// SearchMovie searches for a movie using OMDb
func (p *OMDbProvider) SearchMovie(search MovieSearch, language string) (*MovieMetadata, error) {
	// First, search by title to get the IMDb ID
	searchURL, err := url.Parse(p.baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OMDb URL: %v", err)
	}

	query := searchURL.Query()
	query.Set("apikey", p.apiKey)
	query.Set("s", search.Title)
	query.Set("type", "movie")

	// If we have a year, add it to improve search accuracy
	if search.Year > 0 {
		query.Set("y", strconv.Itoa(search.Year))
	}

	searchURL.RawQuery = query.Encode()

	// Execute the search request
	resp, err := p.client.Get(searchURL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to search movie: %v", err)
	}
	defer resp.Body.Close()

	// Parse the response
	var searchResp OMDbSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %v", err)
	}

	// Check if the search was successful
	if searchResp.Response != "True" {
		// Try a more forgiving search without year if we had one
		if search.Year > 0 {
			// Remove year and try again
			query.Del("y")
			searchURL.RawQuery = query.Encode()

			resp, err := p.client.Get(searchURL.String())
			if err != nil {
				return nil, fmt.Errorf("failed to search movie without year: %v", err)
			}
			defer resp.Body.Close()

			if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
				return nil, fmt.Errorf("failed to decode search response: %v", err)
			}

			if searchResp.Response != "True" {
				return nil, fmt.Errorf("no movies found matching '%s'", search.Title)
			}
		} else {
			return nil, fmt.Errorf("no movies found matching '%s'", search.Title)
		}
	}

	if len(searchResp.Search) == 0 {
		return nil, fmt.Errorf("no movies found matching '%s'", search.Title)
	}

	// Get the first result's IMDb ID
	imdbID := searchResp.Search[0].ImdbID

	// Now get the detailed information using the IMDb ID
	detailURL, _ := url.Parse(p.baseURL)
	detailQuery := detailURL.Query()
	detailQuery.Set("apikey", p.apiKey)
	detailQuery.Set("i", imdbID)
	detailQuery.Set("plot", "full")
	detailURL.RawQuery = detailQuery.Encode()

	// Execute the detail request
	detailResp, err := p.client.Get(detailURL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get movie details: %v", err)
	}
	defer detailResp.Body.Close()

	// Parse the detail response
	var movie OMDbResponse
	if err := json.NewDecoder(detailResp.Body).Decode(&movie); err != nil {
		return nil, fmt.Errorf("failed to decode movie details: %v", err)
	}

	// Check if the response was successful
	if movie.Response != "True" {
		return nil, fmt.Errorf("failed to get movie details: %s", movie.Error)
	}

	// Extract the year from the year field
	year := 0
	if movie.Year != "" {
		// OMDb sometimes returns year ranges like "2008-2013", we only want the first year
		yearPart := strings.Split(movie.Year, "–")[0] // Note: this is an en dash, not a hyphen
		yearPart = strings.Split(yearPart, "-")[0]    // Also handle regular hyphens
		if y, err := strconv.Atoi(yearPart); err == nil {
			year = y
		}
	}

	return &MovieMetadata{
		Title:    movie.Title,
		Year:     year,
		Overview: movie.Plot,
	}, nil
}

// SearchTVShow attempts to search for a TV show using OMDb
func (p *OMDbProvider) SearchTVShow(search TVShowSearch, language string) (*TVShowMetadata, error) {
	// OMDb is not ideal for TV show episode lookup
	// For a more complete implementation, it would be better to use a dedicated TV API
	return nil, fmt.Errorf("TV show search is not fully supported in OMDb. Use TVMaze or TVDb provider instead")
}
