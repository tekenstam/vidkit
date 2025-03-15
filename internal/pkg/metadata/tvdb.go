package metadata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// TVDbProvider implements TV show metadata lookup using TheTVDB API
type TVDbProvider struct {
	apiKey      string
	baseURL     string
	apiToken    string
	client      *http.Client
	tokenExpiry time.Time
}

// TVDbLoginRequest represents a login request to the TVDb API
type TVDbLoginRequest struct {
	ApiKey string `json:"apikey"`
}

// TVDbLoginResponse represents the login response from the TVDb API
type TVDbLoginResponse struct {
	Token string `json:"token"`
}

// TVDbSeriesResponse represents the series response from the TVDb API
type TVDbSeriesResponse struct {
	Data struct {
		ID            int    `json:"id"`
		SeriesName    string `json:"seriesName"`
		FirstAired    string `json:"firstAired"`
		Status        string `json:"status"`
		Network       string `json:"network"`
		Overview      string `json:"overview"`
		Genre         []string `json:"genre"`
		SeriesID      int    `json:"seriesId"`
		ImdbID        string `json:"imdbId"`
	} `json:"data"`
}

// TVDbSearchResponse represents the search response from the TVDb API
type TVDbSearchResponse struct {
	Data []struct {
		ID            int    `json:"id"`
		SeriesName    string `json:"seriesName"`
		FirstAired    string `json:"firstAired"`
		Status        string `json:"status"`
		Network       string `json:"network"`
		Overview      string `json:"overview"`
	} `json:"data"`
}

// TVDbEpisodeResponse represents the episode response from the TVDb API
type TVDbEpisodeResponse struct {
	Data struct {
		ID            int    `json:"id"`
		EpisodeName   string `json:"episodeName"`
		FirstAired    string `json:"firstAired"`
		Overview      string `json:"overview"`
		AiredSeason   int    `json:"airedSeason"`
		AiredEpisodeNumber int `json:"airedEpisodeNumber"`
	} `json:"data"`
}

// TVDbSeasonsResponse represents the seasons response from the TVDb API
type TVDbSeasonsResponse struct {
	Data []struct {
		ID            int    `json:"id"`
		AiredSeasons  []string `json:"airedSeasons"`
	} `json:"data"`
}

// NewTVDbProvider creates a new TVDb metadata provider
func NewTVDbProvider(apiKey string) (*TVDbProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("TVDb API key is required")
	}
	
	return &TVDbProvider{
		apiKey:      apiKey,
		baseURL:     "https://api.thetvdb.com",
		client:      &http.Client{Timeout: 10 * time.Second},
		tokenExpiry: time.Time{}, // Zero time, will be updated on first login
	}, nil
}

// Ensure TVDbProvider implements MetadataProvider
var _ MetadataProvider = (*TVDbProvider)(nil)

// getToken gets or refreshes the API token
func (p *TVDbProvider) getToken() error {
	// Check if token is still valid
	if p.apiToken != "" && time.Now().Before(p.tokenExpiry) {
		return nil
	}
	
	// Prepare login request
	loginURL := fmt.Sprintf("%s/login", p.baseURL)
	loginData := TVDbLoginRequest{
		ApiKey: p.apiKey,
	}
	
	jsonData, err := json.Marshal(loginData)
	if err != nil {
		return fmt.Errorf("failed to marshal login request: %v", err)
	}
	
	// Execute login request
	req, err := http.NewRequest("POST", loginURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create login request: %v", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute login request: %v", err)
	}
	defer resp.Body.Close()
	
	// Parse login response
	var loginResp TVDbLoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return fmt.Errorf("failed to decode login response: %v", err)
	}
	
	// Update token and expiry (tokens are valid for 24 hours)
	p.apiToken = loginResp.Token
	p.tokenExpiry = time.Now().Add(23 * time.Hour) // Slightly shorter to be safe
	
	return nil
}

// sendRequest sends an authenticated request to the TVDb API
func (p *TVDbProvider) sendRequest(method, path string, body []byte) (*http.Response, error) {
	// Ensure we have a valid token
	if err := p.getToken(); err != nil {
		return nil, err
	}
	
	// Prepare request
	reqURL := fmt.Sprintf("%s%s", p.baseURL, path)
	req, err := http.NewRequest(method, reqURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	
	// Add authentication header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.apiToken))
	
	if len(body) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}
	
	// Execute request
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	
	// Check for authentication errors
	if resp.StatusCode == http.StatusUnauthorized {
		resp.Body.Close()
		
		// Token might have expired, try to refresh
		p.apiToken = ""
		if err := p.getToken(); err != nil {
			return nil, err
		}
		
		// Retry request with new token
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.apiToken))
		return p.client.Do(req)
	}
	
	return resp, nil
}

// SearchMovie attempts to search for a movie using TVDb
func (p *TVDbProvider) SearchMovie(search MovieSearch, language string) (*MovieMetadata, error) {
	// TVDb is primarily for TV shows, not ideal for movies
	return nil, fmt.Errorf("movie search is not supported by TVDb. Use TMDb or OMDb provider instead")
}

// SearchTVShow searches for a TV show using TVDb
func (p *TVDbProvider) SearchTVShow(search TVShowSearch, language string) (*TVShowMetadata, error) {
	// Search for the TV show
	searchPath := fmt.Sprintf("/search/series?name=%s", url.QueryEscape(search.Title))
	
	// Add year if available
	if search.Year > 0 {
		searchPath = fmt.Sprintf("%s&year=%d", searchPath, search.Year)
	}
	
	resp, err := p.sendRequest("GET", searchPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to search TV show: %v", err)
	}
	defer resp.Body.Close()
	
	// Parse search response
	var searchResp TVDbSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %v", err)
	}
	
	// Check if we found anything
	if len(searchResp.Data) == 0 {
		return nil, fmt.Errorf("no TV shows found matching '%s'", search.Title)
	}
	
	// Get the first result
	seriesID := searchResp.Data[0].ID
	
	// Get detailed series information
	seriesPath := fmt.Sprintf("/series/%d", seriesID)
	seriesResp, err := p.sendRequest("GET", seriesPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get series details: %v", err)
	}
	defer seriesResp.Body.Close()
	
	// Parse series response
	var series TVDbSeriesResponse
	if err := json.NewDecoder(seriesResp.Body).Decode(&series); err != nil {
		return nil, fmt.Errorf("failed to decode series response: %v", err)
	}
	
	// Get information about seasons
	seasonsPath := fmt.Sprintf("/series/%d/episodes/summary", seriesID)
	seasonsResp, err := p.sendRequest("GET", seasonsPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get seasons information: %v", err)
	}
	defer seasonsResp.Body.Close()
	
	// Parse seasons response
	var seasons TVDbSeasonsResponse
	if err := json.NewDecoder(seasonsResp.Body).Decode(&seasons); err != nil {
		return nil, fmt.Errorf("failed to decode seasons response: %v", err)
	}
	
	// Extract season count
	seasonCount := 0
	if len(seasons.Data) > 0 {
		seasonCount = len(seasons.Data[0].AiredSeasons)
	}
	
	// Extract year from first aired date
	year := 0
	if series.Data.FirstAired != "" {
		if t, err := time.Parse("2006-01-02", series.Data.FirstAired); err == nil {
			year = t.Year()
		}
	}
	
	// Create the result metadata
	metadata := &TVShowMetadata{
		Title:       series.Data.SeriesName,
		Year:        year,
		Overview:    series.Data.Overview,
		Network:     series.Data.Network,
		Status:      series.Data.Status,
		SeasonCount: seasonCount,
		Genres:      series.Data.Genre,
		Season:      search.Season,
		Episode:     search.Episode,
	}
	
	// If we have season and episode information, get episode details
	if search.Season > 0 && search.Episode > 0 {
		episodePath := fmt.Sprintf("/series/%d/episodes/query?airedSeason=%d&airedEpisode=%d", 
			seriesID, search.Season, search.Episode)
		
		episodeResp, err := p.sendRequest("GET", episodePath, nil)
		if err != nil {
			return metadata, nil // Return what we have so far, episode info is optional
		}
		defer episodeResp.Body.Close()
		
		// Parse episode response
		var episode TVDbEpisodeResponse
		if err := json.NewDecoder(episodeResp.Body).Decode(&episode); err != nil {
			return metadata, nil // Return what we have so far
		}
		
		// Update metadata with episode information
		metadata.EpisodeTitle = episode.Data.EpisodeName
		if episode.Data.FirstAired != "" {
			metadata.AirDate = episode.Data.FirstAired
		}
	}
	
	return metadata, nil
}
