package httpclient

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetAnime(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/anime/1/full" {
				t.Errorf("Expected to request '/anime/1/full', got: %s", r.URL.Path)
			}
			response := GetAnimeResponse{
				Data: Data{
					MalID: 1,
					Images: Images{
						Webp: ImageDetails{
							ImageURL: "http://example.com/image.webp",
						},
					},
					Title:    "Test Anime",
					Synopsis: "This is a test synopsis.",
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

		result, err := GetAnime(1)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		expected := AnimeDetailsFull{
			MalID:    1,
			Image:    "http://example.com/image.webp",
			Title:    "Test Anime",
			Synopsis: "This is a test synopsis.",
		}

		if result.Episodes != expected.Episodes {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})

	t.Run("missing JIKAN_BASE_URL", func(t *testing.T) {
		os.Unsetenv("JIKAN_BASE_URL")

		_, err := GetAnime(1)
		if err == nil {
			t.Fatal("Expected an error, got none")
		}
		expectedError := "JIKAN_BASE_URL environment variable is not set"
		if err.Error() != expectedError {
			t.Errorf("Expected error message to be '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("http request error", func(t *testing.T) {
		os.Setenv("JIKAN_BASE_URL", "http://invalid-url/")

		_, err := GetAnime(1)
		if err == nil {
			t.Fatal("Expected an error, got none")
		}
	})

	t.Run("non-200 status code", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound) // Simulate a 404 error
		}))
		defer server.Close()

		os.Setenv("JIKAN_BASE_URL", server.URL+"/")

		_, err := GetAnime(1)
		if err == nil {
			t.Fatal("Expected an error, got none")
		}
		expectedError := "unexpected status: 404 Not Found"
		if err.Error() != expectedError {
			t.Errorf("Expected error message to be '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("json unmarshal error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write([]byte(`{"data": { "mal_id": 1, "images": { "webp": { "image_url": "http://example.com/image.webp" }, "wrong_field": {}}`))
			if err != nil {
				log.Fatal(err)
			}
		}))
		defer server.Close()

		os.Setenv("JIKAN_BASE_URL", server.URL+"/")

		_, err := GetAnime(1)
		if err == nil {
			t.Fatal("Expected an error, got none")
		}
	})
}
