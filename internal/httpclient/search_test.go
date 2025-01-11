package httpclient

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type AnimeData struct {
	MalID  int    `json:"mal_id"`
	URL    string `json:"url"`
	Images struct {
		Webp struct {
			ImageURL string `json:"image_url"`
		} `json:"webp"`
	} `json:"images"`
	Title    string `json:"title"`
	Synopsis string `json:"synopsis"`
}

type JixenAnimeSearchBodyTest struct {
	Pagination struct {
		LastVisiblePage int  `json:"last_visible_page"`
		HasNextPage     bool `json:"has_next_page"`
		CurrentPage     int  `json:"current_page"`
		Items           struct {
			Count   int `json:"count"`
			Total   int `json:"total"`
			PerPage int `json:"per_page"`
		} `json:"items"`
	} `json:"pagination"`
	Data []AnimeData `json:"data"`
}

func TestSearchAnime(t *testing.T) {
	t.Run("successful search", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := JixenAnimeSearchBodyTest{
				Pagination: struct {
					LastVisiblePage int  `json:"last_visible_page"`
					HasNextPage     bool `json:"has_next_page"`
					CurrentPage     int  `json:"current_page"`
					Items           struct {
						Count   int `json:"count"`
						Total   int `json:"total"`
						PerPage int `json:"per_page"`
					} `json:"items"`
				}{
					LastVisiblePage: 1,
					HasNextPage:     false,
					CurrentPage:     1,
					Items: struct {
						Count   int `json:"count"`
						Total   int `json:"total"`
						PerPage int `json:"per_page"`
					}{
						Count:   1,
						Total:   1,
						PerPage: 10,
					},
				},
				Data: []AnimeData{
					{
						MalID: 1,
						URL:   "http://example.com/anime/1",
						Images: struct {
							Webp struct {
								ImageURL string `json:"image_url"`
							} `json:"webp"`
						}{
							Webp: struct {
								ImageURL string `json:"image_url"`
							}{
								ImageURL: "http://example.com/image.webp",
							},
						},
						Title:    "Naruto",
						Synopsis: "A ninja story.",
					},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				log.Fatal(err)
			}
		}))
		defer server.Close()

		os.Setenv("JIKAN_BASE_URL", server.URL+"/")

		result, err := SearchAnime("naruto")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		expected := AnimeSearchResults{
			Pagination: struct {
				Count int
				Total int
			}{Count: 1, Total: 1},
			Data: []AnimeDetails{
				{
					MalID:    1,
					Image:    "http://example.com/image.webp",
					Title:    "Naruto",
					Synopsis: "A ninja story.",
				},
			},
		}

		if result.Data[0].Title != expected.Data[0].Title {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})

	t.Run("short query error", func(t *testing.T) {
		_, err := SearchAnime("na")
		if err == nil || err.Error() != "search should be more than 3 characters" {
			t.Fatalf("Expected error for short query, got: %v", err)
		}
	})

	t.Run("missing JIKAN_BASE_URL", func(t *testing.T) {
		os.Unsetenv("JIKAN_BASE_URL")
		_, err := SearchAnime("naruto")
		if err == nil || err.Error() != "failed to load environment variables: " {
			t.Fatalf("Expected error for missing JIKAN_BASE_URL, got: %v", err)
		}
	})

	t.Run("HTTP error", func(t *testing.T) {
		server := httptest.NewServer(http.NotFoundHandler())
		defer server.Close()

		os.Setenv("JIKAN_BASE_URL", server.URL+"/")

		_, err := SearchAnime("naruto")
		if err == nil {
			t.Fatalf("Expected an HTTP error, got: %v", err)
		}
	})
}
