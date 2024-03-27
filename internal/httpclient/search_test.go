package httpclient

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
)

type MockHTTPClient struct{}

var testAnimeSearchResults = AnimeSearchResults{
	Pagination: struct {
		Count int
		Total int
	}{
		Count: 2,
		Total: 10,
	},
	Data: []AnimeDetails{
		{
			Title:    "Naruto",
			MalID:    1,
			Image:    "naruto.webp",
			Synopsis: "The story of Naruto Uzumaki, a young ninja who seeks recognition from his peers and dreams of becoming the Hokage, the leader of his village.",
		},
		{
			Title:    "One Piece",
			MalID:    2,
			Image:    "one_piece.webp",
			Synopsis: "Follows the adventures of Monkey D. Luffy and his pirate crew in order to find the greatest treasure ever left by the legendary Pirate, Gold Roger.",
		},
	},
}

func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
	if strings.Contains(url, "error") {
		return nil, errors.New("error occurred")
	}

	jsonData := []byte(`{
		"pagination": {
			"last_visible_page": 10,
			"has_next_page": true,
			"current_page": 1,
			"items": {
				"count": 2,
				"total": 10,
				"per_page": 10
			}
		},
		"data": [
			{
				"mal_id": 1,
				"url": "https://example.com",
				"images": {
					"jpg": {
						"image_url": "naruto.jpg",
						"small_image_url": "https://example.com/small_image.jpg",
						"large_image_url": "https://example.com/large_image.jpg"
					},
					"webp": {
						"image_url": "naruto.webp",
						"small_image_url": "https://example.com/small_image.webp",
						"large_image_url": "https://example.com/large_image.webp"
					}
				},
				"approved": true,
				"titles": [
					{
						"type": "main",
						"title": "Example Anime"
					}
				],
				"title": "Naruto",
				"synopsis": "The story of Naruto Uzumaki, a young ninja who seeks recognition from his peers and dreams of becoming the Hokage, the leader of his village.",
				"season": "Winter 2000",
				"year": 2000
			},
			{
				"mal_id": 2,
				"url": "https://example.com",
				"images": {
					"jpg": {
						"image_url": "one_piece.jpg",
						"small_image_url": "https://example.com/small_image.jpg",
						"large_image_url": "https://example.com/large_image.jpg"
					},
					"webp": {
						"image_url": "one_piece.webp",
						"small_image_url": "https://example.com/small_image.webp",
						"large_image_url": "https://example.com/large_image.webp"
					}
				},
				"approved": true,
				"titles": [
					{
						"type": "main",
						"title": "Example Anime"
					}
				],
				"title": "One Piece",
				"synopsis": "Follows the adventures of Monkey D. Luffy and his pirate crew in order to find the greatest treasure ever left by the legendary Pirate, Gold Roger.",
				"season": "Winter 2000",
				"year": 2000
			}
		]
	}
	`)

	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(string(jsonData))),
	}, nil
}

func TestSearchAnime(t *testing.T) {
	// Set up test environment
	os.Setenv("JIKAN_BASE_URL", "http://example.com/")

	// Test cases
	tests := []struct {
		name           string
		query          string
		expectedOutput AnimeSearchResults
		expectedError  string
		httpClient     HTTPClient
	}{
		{
			name:           "Valid query",
			query:          "Example",
			expectedOutput: testAnimeSearchResults,
			expectedError:  "",
			httpClient:     &MockHTTPClient{},
		},
		{
			name:           "Query less than 3 characters",
			query:          "ab",
			expectedOutput: AnimeSearchResults{},
			expectedError:  "search should be more than 3 characters",
			httpClient:     &MockHTTPClient{},
		},
		{
			name:           "HTTP client returns error",
			query:          "error",
			expectedOutput: AnimeSearchResults{},
			expectedError:  "error connecting to http client error occurred",
			httpClient:     &MockHTTPClient{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output, err := SearchAnime(tc.query, tc.httpClient)

			if !reflect.DeepEqual(output, tc.expectedOutput) {
				t.Errorf("got %v; want %v", output, tc.expectedOutput)
			}
			if err != nil {
				if err.Error() != tc.expectedError {
					t.Errorf("le cry got %v, want %v", err, tc.expectedError)
				}
			}

		})
	}
}
