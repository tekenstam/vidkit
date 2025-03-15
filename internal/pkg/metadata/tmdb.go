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
	Genres   []string
}

// TVShowSearch represents a TV show search request
type TVShowSearch struct {
	Title        string
	Year         int
	Season       int
	Episode      int
	EpisodeTitle string
}

// TVShowMetadata represents TV show metadata from TMDb
type TVShowMetadata struct {
	Title        string
	Year         int
	Overview     string
	Season       int
	Episode      int
	EpisodeTitle string
	SeasonCount  int
	Network      string
	AirDate      string
	Status       string
	Genres       []string
}

// MetadataProvider defines the interface for metadata providers
type MetadataProvider interface {
	SearchMovie(search MovieSearch, language string) (*MovieMetadata, error)
	SearchTVShow(search TVShowSearch, language string) (*TVShowMetadata, error)
}

// TMDbClient defines the interface for TMDb operations
// Only includes the methods we actually use from the TMDb client
type TMDbClient interface {
	GetSearchMovies(query string, urlOptions map[string]string) (*tmdb.SearchMovies, error)
	GetMovieDetails(id int, urlOptions map[string]string) (*tmdb.MovieDetails, error)
}

// TMDbProvider implements movie metadata lookup using TMDb
type TMDbProvider struct {
	client TMDbClient
}

// Ensure TMDbProvider implements MetadataProvider
var _ MetadataProvider = (*TMDbProvider)(nil)

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

	// Extract genre names
	genreNames := make([]string, 0, len(movie.Genres))
	for _, genre := range movie.Genres {
		genreNames = append(genreNames, genre.Name)
	}

	return &MovieMetadata{
		Title:    movie.Title,
		Year:     year,
		Overview: movie.Overview,
		Genres:   genreNames,
	}, nil
}

// SearchTVShow attempts to search for a TV show using TMDb (not fully implemented)
func (p *TMDbProvider) SearchTVShow(search TVShowSearch, language string) (*TVShowMetadata, error) {
	// This is a stub implementation since TMDb's API for TV shows works differently
	// In a real implementation, we would use GetSearchTV and GetTVDetails
	return nil, fmt.Errorf("TMDb TV show search not fully implemented. Use TvMaze provider instead")
}

// ExtractMovieInfo extracts movie information from a filename
func ExtractMovieInfo(filename string) MovieSearch {
	// Extract base name without extension
	basename := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))

	// Check if it's a TV show pattern first
	tvInfo := ExtractTVShowInfo(filename)
	if tvInfo.Season > 0 && tvInfo.Episode > 0 {
		// This is a TV show filename, not a movie
		return MovieSearch{Title: "", Year: 0}
	}

	// Look for year pattern (YYYY) in filename, but only in delimiters
	// Only accept years in parentheses or square brackets
	yearPattern := regexp.MustCompile(`\((\d{4})\)|\[(\d{4})\]`)
	year := 0
	title := basename

	// Try all possible year patterns
	for _, matches := range yearPattern.FindAllStringSubmatch(basename, -1) {
		for i := 1; i < len(matches); i++ {
			if matches[i] != "" {
				if y, err := strconv.Atoi(matches[i]); err == nil {
					year = y
					// Remove year from title
					title = yearPattern.ReplaceAllString(basename, " ")
					break
				}
			}
		}
		if year != 0 {
			break
		}
	}

	// Only clean up the title if we found a year in valid delimiters
	if year > 0 {
		// Clean up the title by removing common patterns
		title = cleanTitle(title)
	} else {
		// For files without valid delimited years, just return the original
		// This preserves formats like "The.Matrix.1999.mp4"
		return MovieSearch{
			Title: basename,
			Year:  0,
		}
	}

	return MovieSearch{
		Title: title,
		Year:  year,
	}
}

// ExtractTVShowInfo extracts TV show information from a filename
func ExtractTVShowInfo(filename string) TVShowSearch {
	// Extract base name without extension
	basename := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))

	// Look for common TV show patterns
	// Pattern 1: ShowName.S01E02
	seasonEpisodePattern1 := regexp.MustCompile(`(?i)(.*?)[\s\._-]*s(\d{1,2})[\s\._-]*e(\d{1,2})(?:[\s\._-]*(.*))?`)

	// Pattern 2: ShowName.1x02
	seasonEpisodePattern2 := regexp.MustCompile(`(?i)(.*?)[\s\._-]*(\d{1,2})x(\d{1,2})(?:[\s\._-]*(.*))?`)

	// Pattern 3: ShowName.Season.1.Episode.2
	seasonEpisodePattern3 := regexp.MustCompile(`(?i)(.*?)[\s\._-]*(?:season|s)[\s\._-]*(\d{1,2})[\s\._-]*(?:episode|ep|e)[\s\._-]*(\d{1,2})(?:[\s\._-]*(.*))?`)

	// Try all season/episode patterns
	var title string
	var season, episode int
	var episodeTitle string

	if m := seasonEpisodePattern1.FindStringSubmatch(basename); len(m) >= 4 {
		title = m[1]
		season, _ = strconv.Atoi(m[2])
		episode, _ = strconv.Atoi(m[3])
		if len(m) > 4 {
			episodeTitle = m[4]
		}
	} else if m := seasonEpisodePattern2.FindStringSubmatch(basename); len(m) >= 4 {
		title = m[1]
		season, _ = strconv.Atoi(m[2])
		episode, _ = strconv.Atoi(m[3])
		if len(m) > 4 {
			episodeTitle = m[4]
		}
	} else if m := seasonEpisodePattern3.FindStringSubmatch(basename); len(m) >= 4 {
		title = m[1]
		season, _ = strconv.Atoi(m[2])
		episode, _ = strconv.Atoi(m[3])
		if len(m) > 4 {
			episodeTitle = m[4]
		}
	} else {
		// If no TV pattern is found, just return the title
		return TVShowSearch{
			Title: cleanTitle(basename),
		}
	}

	// Look for year pattern (YYYY) in title, but only in delimiters
	// Only accept years in parentheses or square brackets
	yearPattern := regexp.MustCompile(`\((\d{4})\)|\[(\d{4})\]`)
	year := 0

	// Try all possible year patterns
	for _, matches := range yearPattern.FindAllStringSubmatch(title, -1) {
		for i := 1; i < len(matches); i++ {
			if matches[i] != "" {
				if y, err := strconv.Atoi(matches[i]); err == nil {
					year = y
					// Remove year from title
					title = yearPattern.ReplaceAllString(title, " ")
					break
				}
			}
		}
		if year != 0 {
			break
		}
	}

	// Clean up the titles
	title = cleanTitle(title)

	// Clean up episode title - filter out quality and technical info
	if episodeTitle != "" {
		// List of common quality/codec terms to filter out
		qualityTerms := []string{"1080p", "720p", "480p", "HEVC", "h264", "x264", "HDRip", "BRRip", "BluRay", "WEB-DL", "HDTV"}

		cleanEpisodeTitle := episodeTitle

		for _, term := range qualityTerms {
			if strings.Contains(strings.ToLower(episodeTitle), strings.ToLower(term)) {
				// If episode title is just a quality indicator, set it to empty
				cleanEpisodeTitle = regexp.MustCompile(`(?i)`+term).ReplaceAllString(cleanEpisodeTitle, "")
			}
		}

		// If after removing all quality terms, only spaces remain, consider it empty
		if strings.TrimSpace(cleanEpisodeTitle) == "" {
			episodeTitle = ""
		} else {
			episodeTitle = cleanTitle(episodeTitle)
		}
	}

	return TVShowSearch{
		Title:        title,
		Year:         year,
		Season:       season,
		Episode:      episode,
		EpisodeTitle: episodeTitle,
	}
}

// Helper to clean up titles
func cleanTitle(title string) string {
	// Clean up the title by removing common patterns
	title = strings.NewReplacer(
		"1080p", "",
		"720p", "",
		"480p", "",
		"360p", "",
		"h264", "",
		"x264", "",
		"HDRip", "",
		"BRRip", "",
		"BluRay", "",
		"WEB-DL", "",
		"HDTV", "",
		"[", " ",
		"]", " ",
		"(", " ",
		")", " ",
		".", " ",
		"_", " ",
	).Replace(title)

	// Clean up extra spaces
	title = regexp.MustCompile(`\s+`).ReplaceAllString(strings.TrimSpace(title), " ")

	return title
}
