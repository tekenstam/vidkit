package metadata

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// TvMazeProvider implements TV show metadata lookup using TvMaze API
type TvMazeProvider struct {
	baseURL string
	client  *http.Client
}

// Ensure TvMazeProvider implements MetadataProvider
var _ MetadataProvider = (*TvMazeProvider)(nil)

// TvMazeShow represents a TV show from TvMaze API
type TvMazeShow struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	Language     string   `json:"language"`
	Genres       []string `json:"genres"`
	Status       string   `json:"status"`
	Runtime      int      `json:"runtime"`
	Premiered    string   `json:"premiered"`
	OfficialSite string   `json:"officialSite"`
	Network      struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Country struct {
			Name     string `json:"name"`
			Code     string `json:"code"`
			Timezone string `json:"timezone"`
		} `json:"country"`
	} `json:"network"`
	Summary string `json:"summary"`
	Updated int64  `json:"updated"`
	Rating  struct {
		Average float64 `json:"average"`
	} `json:"rating"`
	Embedded struct {
		Seasons []struct {
			ID           int    `json:"id"`
			Number       int    `json:"number"`
			Name         string `json:"name"`
			EpisodeOrder int    `json:"episodeOrder"`
		} `json:"seasons"`
	} `json:"_embedded"`
}

// TvMazeSearchResult represents a search result from TvMaze API
type TvMazeSearchResult struct {
	Score float64    `json:"score"`
	Show  TvMazeShow `json:"show"`
}

// TvMazeEpisode represents a TV episode from TvMaze API
type TvMazeEpisode struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Season  int    `json:"season"`
	Number  int    `json:"number"`
	Airdate string `json:"airdate"`
	Runtime int    `json:"runtime"`
	Summary string `json:"summary"`
	Type    string `json:"type"`
}

// NewTvMazeProvider creates a new TvMaze metadata provider
func NewTvMazeProvider() *TvMazeProvider {
	return &TvMazeProvider{
		baseURL: "https://api.tvmaze.com",
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

// SearchMovie returns an error because TvMaze is for TV shows, not movies
func (p *TvMazeProvider) SearchMovie(search MovieSearch, language string) (*MovieMetadata, error) {
	return nil, fmt.Errorf("TvMaze does not support movie lookups")
}

// SearchTVShow searches for a TV show using TvMaze API
func (p *TvMazeProvider) SearchTVShow(search TVShowSearch, language string) (*TVShowMetadata, error) {
	// Construct search URL
	query := url.QueryEscape(search.Title)
	searchURL := fmt.Sprintf("%s/search/shows?q=%s", p.baseURL, query)

	// Make HTTP request
	resp, err := p.client.Get(searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to search TV show: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to search TV show: %s", resp.Status)
	}

	// Parse response
	var searchResults []TvMazeSearchResult
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	err = json.Unmarshal(body, &searchResults)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	if len(searchResults) == 0 {
		return nil, fmt.Errorf("no TV shows found matching '%s'", search.Title)
	}

	// Get the first result
	showID := searchResults[0].Show.ID

	// Get show details with seasons information
	showURL := fmt.Sprintf("%s/shows/%d?embed=seasons", p.baseURL, showID)
	showResp, err := p.client.Get(showURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get show details: %v", err)
	}
	defer showResp.Body.Close()

	if showResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get show details: %s", showResp.Status)
	}

	// Parse show details
	var show TvMazeShow
	showBody, err := io.ReadAll(showResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	err = json.Unmarshal(showBody, &show)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	// Extract year from premiere date
	year := 0
	if show.Premiered != "" {
		if t, err := time.Parse("2006-01-02", show.Premiered); err == nil {
			year = t.Year()
		}
	}

	// Clean HTML tags from summary
	summary := cleanHtmlTags(show.Summary)

	// Basic metadata without episode info
	metadata := &TVShowMetadata{
		Title:       show.Name,
		Year:        year,
		Overview:    summary,
		SeasonCount: len(show.Embedded.Seasons),
		Network:     show.Network.Name,
		Status:      show.Status,
		Genres:      show.Genres,
	}

	// If season and episode are provided, get episode details
	if search.Season > 0 && search.Episode > 0 {
		episodeURL := fmt.Sprintf("%s/shows/%d/episodebynumber?season=%d&number=%d", p.baseURL, showID, search.Season, search.Episode)
		episodeResp, err := p.client.Get(episodeURL)
		if err != nil {
			// Return show info without episode details
			return metadata, nil
		}
		defer episodeResp.Body.Close()

		if episodeResp.StatusCode != http.StatusOK {
			// Return show info without episode details
			return metadata, nil
		}

		// Parse episode details
		var episode TvMazeEpisode
		episodeBody, err := io.ReadAll(episodeResp.Body)
		if err != nil {
			// Return show info without episode details
			return metadata, nil
		}

		err = json.Unmarshal(episodeBody, &episode)
		if err != nil {
			// Return show info without episode details
			return metadata, nil
		}

		// Add episode details to metadata
		metadata.Season = episode.Season
		metadata.Episode = episode.Number
		metadata.EpisodeTitle = episode.Name
		metadata.AirDate = episode.Airdate
	}

	return metadata, nil
}

// Helper function to clean HTML tags from text
func cleanHtmlTags(html string) string {
	// Remove the <p> and </p> tags that are common in TvMaze responses
	cleaned := strings.ReplaceAll(html, "<p>", "")
	cleaned = strings.ReplaceAll(cleaned, "</p>", "")
	cleaned = strings.ReplaceAll(cleaned, "<b>", "")
	cleaned = strings.ReplaceAll(cleaned, "</b>", "")
	cleaned = strings.ReplaceAll(cleaned, "<i>", "")
	cleaned = strings.ReplaceAll(cleaned, "</i>", "")

	// Replace HTML entities
	cleaned = strings.ReplaceAll(cleaned, "&nbsp;", " ")
	cleaned = strings.ReplaceAll(cleaned, "&amp;", "&")
	cleaned = strings.ReplaceAll(cleaned, "&lt;", "<")
	cleaned = strings.ReplaceAll(cleaned, "&gt;", ">")
	cleaned = strings.ReplaceAll(cleaned, "&quot;", "\"")

	return strings.TrimSpace(cleaned)
}
