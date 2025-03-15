package metadata

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOMDbProvider_SearchMovie(t *testing.T) {
	// Create a mock HTTP server for OMDb
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters
		q := r.URL.Query()
		apiKey := q.Get("apikey")
		searchTitle := q.Get("s")
		// searchYear removed as it was unused
		imdbID := q.Get("i")

		// Check if API key is present
		if apiKey == "" {
			http.Error(w, "No API key provided", http.StatusUnauthorized)
			return
		}

		// Handle search request
		if searchTitle != "" {
			// Return search results
			w.Header().Set("Content-Type", "application/json")
			if searchTitle == "The Matrix" {
				// Return search results
				response := `{
					"Search": [
						{
							"Title": "The Matrix",
							"Year": "1999",
							"imdbID": "tt0133093",
							"Type": "movie",
							"Poster": "https://m.media-amazon.com/images/M/MV5BNzQzOTk3OTAtNDQ0Zi00ZTVkLWI0MTEtMDllZjNkYzNjNTc4L2ltYWdlXkEyXkFqcGdeQXVyNjU0OTQ0OTY@._V1_SX300.jpg"
						},
						{
							"Title": "The Matrix Reloaded",
							"Year": "2003",
							"imdbID": "tt0234215",
							"Type": "movie",
							"Poster": "https://m.media-amazon.com/images/M/MV5BODE0MzZhZTgtYzkwYi00YmI5LThlZWYtOWRmNWE5ODk0NzMxXkEyXkFqcGdeQXVyNjU0OTQ0OTY@._V1_SX300.jpg"
						}
					],
					"totalResults": "2",
					"Response": "True"
				}`
				w.Write([]byte(response))
				return
			} else if searchTitle == "NonExistentMovie" {
				// Return empty search results
				response := `{
					"Response": "False",
					"Error": "Movie not found!"
				}`
				w.Write([]byte(response))
				return
			}
		}

		// Handle movie details request by IMDb ID
		if imdbID != "" {
			w.Header().Set("Content-Type", "application/json")
			if imdbID == "tt0133093" {
				// Return The Matrix details
				response := `{
					"Title": "The Matrix",
					"Year": "1999",
					"Rated": "R",
					"Released": "31 Mar 1999",
					"Runtime": "136 min",
					"Genre": "Action, Sci-Fi",
					"Director": "Lana Wachowski, Lilly Wachowski",
					"Writer": "Lilly Wachowski, Lana Wachowski",
					"Actors": "Keanu Reeves, Laurence Fishburne, Carrie-Anne Moss",
					"Plot": "When a beautiful stranger leads computer hacker Neo to a forbidding underworld, he discovers the shocking truth--the life he knows is the elaborate deception of an evil cyber-intelligence.",
					"Language": "English",
					"Country": "United States, Australia",
					"Awards": "Won 4 Oscars. 42 wins & 51 nominations total",
					"Poster": "https://m.media-amazon.com/images/M/MV5BNzQzOTk3OTAtNDQ0Zi00ZTVkLWI0MTEtMDllZjNkYzNjNTc4L2ltYWdlXkEyXkFqcGdeQXVyNjU0OTQ0OTY@._V1_SX300.jpg",
					"Ratings": [
						{
							"Source": "Internet Movie Database",
							"Value": "8.7/10"
						},
						{
							"Source": "Rotten Tomatoes",
							"Value": "88%"
						},
						{
							"Source": "Metacritic",
							"Value": "73/100"
						}
					],
					"Metascore": "73",
					"imdbRating": "8.7",
					"imdbVotes": "1,796,248",
					"imdbID": "tt0133093",
					"Type": "movie",
					"DVD": "21 Sep 1999",
					"BoxOffice": "$171,479,930",
					"Production": "N/A",
					"Website": "N/A",
					"Response": "True"
				}`
				w.Write([]byte(response))
				return
			} else {
				// Return error for unknown ID
				response := `{
					"Response": "False",
					"Error": "Error getting data."
				}`
				w.Write([]byte(response))
				return
			}
		}

		// Default error response
		http.Error(w, "Invalid request", http.StatusBadRequest)
	}))
	defer server.Close()

	tests := []struct {
		name      string
		search    MovieSearch
		language  string
		wantTitle string
		wantYear  int
		wantErr   bool
	}{
		{
			name: "Successful search",
			search: MovieSearch{
				Title: "The Matrix",
				Year:  1999,
			},
			language:  "en",
			wantTitle: "The Matrix",
			wantYear:  1999,
			wantErr:   false,
		},
		{
			name: "Movie not found",
			search: MovieSearch{
				Title: "NonExistentMovie",
				Year:  2025,
			},
			language:  "en",
			wantTitle: "",
			wantYear:  0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a provider that uses our mock server
			provider := &OMDbProvider{
				apiKey:  "test_api_key",
				baseURL: server.URL,
				client:  server.Client(),
			}

			// Call the method under test
			got, err := provider.SearchMovie(tt.search, tt.language)

			// Check if we got the expected error state
			if (err != nil) != tt.wantErr {
				t.Errorf("OMDbProvider.SearchMovie() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If we expected an error, no need to check the result
			if tt.wantErr {
				return
			}

			// Check if we got the expected title
			if got.Title != tt.wantTitle {
				t.Errorf("OMDbProvider.SearchMovie() title = %v, want %v", got.Title, tt.wantTitle)
			}

			// Check if we got the expected year
			if got.Year != tt.wantYear {
				t.Errorf("OMDbProvider.SearchMovie() year = %v, want %v", got.Year, tt.wantYear)
			}

			// Check if we got a non-empty overview
			if got.Overview == "" {
				t.Errorf("OMDbProvider.SearchMovie() overview is empty")
			}
		})
	}
}

func TestNewOMDbProvider(t *testing.T) {
	tests := []struct {
		name    string
		apiKey  string
		wantErr bool
	}{
		{
			name:    "Valid API key",
			apiKey:  "test_api_key",
			wantErr: false,
		},
		{
			name:    "Empty API key",
			apiKey:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewOMDbProvider(tt.apiKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewOMDbProvider() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
