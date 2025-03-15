package metadata

import (
	"fmt"
	"testing"

	tmdb "github.com/cyruzin/golang-tmdb"
)

func TestExtractMovieInfo(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     MovieSearch
	}{
		{
			name:     "Standard format with year",
			filename: "Big Buck Bunny (2008).mp4",
			want:     MovieSearch{Title: "Big Buck Bunny", Year: 2008},
		},
		{
			name:     "With resolution and codec",
			filename: "Big Buck Bunny (2008) [1080p x264].mp4",
			want:     MovieSearch{Title: "Big Buck Bunny", Year: 2008},
		},
		{
			name:     "With scene style dots",
			filename: "Big.Buck.Bunny.2008.1080p.x264.mp4",
			want:     MovieSearch{Title: "Big Buck Bunny", Year: 2008},
		},
		{
			name:     "Without year",
			filename: "Big Buck Bunny.mp4",
			want:     MovieSearch{Title: "Big Buck Bunny", Year: 0},
		},
		{
			name:     "With multiple brackets",
			filename: "Big Buck Bunny [2008] [1080p].mp4",
			want:     MovieSearch{Title: "Big Buck Bunny", Year: 2008},
		},
		{
			name:     "With quality tags",
			filename: "Big Buck Bunny (2008) HDRip 1080p.mp4",
			want:     MovieSearch{Title: "Big Buck Bunny", Year: 2008},
		},
		{
			name:     "With release group",
			filename: "[GROUP] Big Buck Bunny (2008).mp4",
			want:     MovieSearch{Title: "Big Buck Bunny", Year: 2008},
		},
		{
			name:     "With multiple years (use first)",
			filename: "Big Buck Bunny (2008) (2009).mp4",
			want:     MovieSearch{Title: "Big Buck Bunny", Year: 2008},
		},
		{
			name:     "With invalid year format",
			filename: "Big Buck Bunny (208).mp4",
			want:     MovieSearch{Title: "Big Buck Bunny (208)", Year: 0},
		},
		{
			name:     "With multiple resolutions",
			filename: "Big Buck Bunny (2008) [1080p] [720p].mp4",
			want:     MovieSearch{Title: "Big Buck Bunny", Year: 2008},
		},
		{
			name:     "With common release tags",
			filename: "Big Buck Bunny (2008) EXTENDED 1080p BRRip x264.mp4",
			want:     MovieSearch{Title: "Big Buck Bunny", Year: 2008},
		},
		{
			name:     "With dots and underscores",
			filename: "Big_Buck_Bunny.2008.1080p_x264.mp4",
			want:     MovieSearch{Title: "Big Buck Bunny", Year: 2008},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractMovieInfo(tt.filename)
			if got.Title != tt.want.Title {
				t.Errorf("ExtractMovieInfo() title = %v, want %v", got.Title, tt.want.Title)
			}
			if got.Year != tt.want.Year {
				t.Errorf("ExtractMovieInfo() year = %v, want %v", got.Year, tt.want.Year)
			}
		})
	}
}

type mockTMDbClient struct {
	searchResults *tmdb.SearchMovies
	movieDetails  *tmdb.MovieDetails
	shouldError   bool
}

func (m *mockTMDbClient) GetSearchMovies(query string, urlOptions map[string]string) (*tmdb.SearchMovies, error) {
	if m.shouldError {
		return nil, fmt.Errorf("API error")
	}
	return m.searchResults, nil
}

func (m *mockTMDbClient) GetMovieDetails(id int, urlOptions map[string]string) (*tmdb.MovieDetails, error) {
	if m.shouldError {
		return nil, fmt.Errorf("API error")
	}
	return m.movieDetails, nil
}

func TestSearchMovie(t *testing.T) {
	movieResult := map[string]interface{}{
		"id":    int64(1),
		"title": "Big Buck Bunny",
	}

	tests := []struct {
		name            string
		search          MovieSearch
		language        string
		mockClient      mockTMDbClient
		want            *MovieMetadata
		wantErr         bool
		wantErrMessage  string
	}{
		{
			name: "Successful search with exact match",
			search: MovieSearch{
				Title: "Big Buck Bunny",
				Year:  2008,
			},
			language: "en",
			mockClient: mockTMDbClient{
				searchResults: &tmdb.SearchMovies{
					Page:         1,
					TotalResults: 1,
					TotalPages:   1,
					Results:      []interface{}{movieResult},
				},
				movieDetails: &tmdb.MovieDetails{
					Title:       "Big Buck Bunny",
					ReleaseDate: "2008-04-10",
					Overview:    "A test movie",
				},
			},
			want: &MovieMetadata{
				Title:    "Big Buck Bunny",
				Year:     2008,
				Overview: "A test movie",
			},
			wantErr: false,
		},
		{
			name: "No results found",
			search: MovieSearch{
				Title: "Nonexistent Movie",
				Year:  2000,
			},
			language: "en",
			mockClient: mockTMDbClient{
				searchResults: &tmdb.SearchMovies{
					Results: []interface{}{},
				},
			},
			want:           nil,
			wantErr:        true,
			wantErrMessage: "no movies found matching 'Nonexistent Movie'",
		},
		{
			name: "API error during search",
			search: MovieSearch{
				Title: "Test Movie",
				Year:  2000,
			},
			language: "en",
			mockClient: mockTMDbClient{
				shouldError: true,
			},
			want:           nil,
			wantErr:        true,
			wantErrMessage: "failed to search movie: API error",
		},
		{
			name: "Search without year",
			search: MovieSearch{
				Title: "Big Buck Bunny",
			},
			language: "en",
			mockClient: mockTMDbClient{
				searchResults: &tmdb.SearchMovies{
					Results: []interface{}{movieResult},
				},
				movieDetails: &tmdb.MovieDetails{
					Title:       "Big Buck Bunny",
					ReleaseDate: "2008-04-10",
					Overview:    "A test movie",
				},
			},
			want: &MovieMetadata{
				Title:    "Big Buck Bunny",
				Year:     2008,
				Overview: "A test movie",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &TMDbProvider{client: &tt.mockClient}
			got, err := provider.SearchMovie(tt.search, tt.language)

			if tt.wantErr {
				if err == nil {
					t.Error("SearchMovie() error = nil, wantErr true")
					return
				}
				if tt.wantErrMessage != "" && err.Error() != tt.wantErrMessage {
					t.Errorf("SearchMovie() error = %v, wantErrMessage %v", err, tt.wantErrMessage)
				}
				return
			}

			if err != nil {
				t.Errorf("SearchMovie() unexpected error = %v", err)
				return
			}

			if got.Title != tt.want.Title {
				t.Errorf("SearchMovie() title = %v, want %v", got.Title, tt.want.Title)
			}
			if got.Year != tt.want.Year {
				t.Errorf("SearchMovie() year = %v, want %v", got.Year, tt.want.Year)
			}
			if got.Overview != tt.want.Overview {
				t.Errorf("SearchMovie() overview = %v, want %v", got.Overview, tt.want.Overview)
			}
		})
	}
}

func TestNewTMDbProvider(t *testing.T) {
	tests := []struct {
		name    string
		apiKey  string
		wantErr bool
	}{
		{
			name:    "Valid API key",
			apiKey:  "valid_key",
			wantErr: false,
		},
		{
			name:    "Empty API key",
			apiKey:  "",
			wantErr: false, // TMDb client initialization doesn't validate the key
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := NewTMDbProvider(tt.apiKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTMDbProvider() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && provider == nil {
				t.Error("NewTMDbProvider() provider is nil, want non-nil")
			}
		})
	}
}
