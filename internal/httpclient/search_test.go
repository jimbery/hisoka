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

func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
	if strings.Contains(url, "error") {
		return nil, errors.New("error occurred")
	}
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(`{"Data":[{"title_english":"Example Title"}]}`)),
	}, nil
}

func TestSearchAnime(t *testing.T) {
	// Set up test environment
	os.Setenv("JIKAN_BASE_URL", "http://example.com/")

	// Test cases
	tests := []struct {
		name           string
		query          string
		expectedOutput AnimeDetails
		expectedError  string
		httpClient     HTTPClient
	}{
		{
			name:           "Valid query",
			query:          "Example",
			expectedOutput: AnimeDetails{title: "Example Title"},
			expectedError:  "",
			httpClient:     &MockHTTPClient{},
		},
		{
			name:           "Query less than 3 characters",
			query:          "ab",
			expectedOutput: AnimeDetails{},
			expectedError:  "search should be more than 3 characters",
			httpClient:     &MockHTTPClient{},
		},
		{
			name:           "HTTP client returns error",
			query:          "error",
			expectedOutput: AnimeDetails{},
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
