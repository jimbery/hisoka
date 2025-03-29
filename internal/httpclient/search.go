package httpclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

type RealHTTPClient struct{}

func (c *RealHTTPClient) Get(url string) (*http.Response, error) {
	return http.Get(url)
}

func SearchAnime(q string) (AnimeSearchResults, error) {
	if len(q) < 3 {
		return AnimeSearchResults{}, fmt.Errorf("search should be more than 3 characters")
	}

	jixenURL := os.Getenv("JIKAN_BASE_URL")

	if jixenURL == "" {
		return AnimeSearchResults{}, fmt.Errorf("failed to load environment variables: %s", jixenURL)
	}

	resp, err := http.Get(jixenURL + "anime?q=" + q + "&limit=20&order_by=popularity")
	if err != nil {
		return AnimeSearchResults{}, fmt.Errorf("error connecting to http client %s", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return AnimeSearchResults{}, fmt.Errorf("error reading request body %s", err)
	}

	var result JixenAnimeSearchBody
	if err := json.Unmarshal(body, &result); err != nil {
		return AnimeSearchResults{}, fmt.Errorf("error unmarshaling json %s", err)
	}

	var output AnimeSearchResults

	output.Pagination.Count = result.Pagination.Items.Count
	output.Pagination.Total = result.Pagination.Items.Total

	for _, element := range result.Data {
		result := AnimeDetails{
			MalID:    element.MalID,
			Image:    element.Images.Webp.ImageURL,
			Title:    element.Title,
			Synopsis: element.Synopsis,
		}

		output.Data = append(output.Data, result)
	}

	return output, nil
}
