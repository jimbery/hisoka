package httpclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type AnimeDetails struct {
	Title    string
	MalID    int
	Image    string
	Synopsis string
}

type AnimeSearchResults struct {
	Pagination struct {
		Count int
		Total int
	}
	Data []AnimeDetails
}

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

// RealHTTPClient implements the HTTPClient interface using real HTTP requests
type RealHTTPClient struct{}

func (c *RealHTTPClient) Get(url string) (*http.Response, error) {
	return http.Get(url)
}

type JixenAnimeSearchBody struct {
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
	Data []struct {
		MalID  int    `json:"mal_id"`
		URL    string `json:"url"`
		Images struct {
			Jpg struct {
				ImageURL      string `json:"image_url"`
				SmallImageURL string `json:"small_image_url"`
				LargeImageURL string `json:"large_image_url"`
			} `json:"jpg"`
			Webp struct {
				ImageURL      string `json:"image_url"`
				SmallImageURL string `json:"small_image_url"`
				LargeImageURL string `json:"large_image_url"`
			} `json:"webp"`
		} `json:"images"`
		Trailer struct {
			YoutubeID string `json:"youtube_id"`
			URL       string `json:"url"`
			EmbedURL  string `json:"embed_url"`
			Images    struct {
				ImageURL        string `json:"image_url"`
				SmallImageURL   string `json:"small_image_url"`
				MediumImageURL  string `json:"medium_image_url"`
				LargeImageURL   string `json:"large_image_url"`
				MaximumImageURL string `json:"maximum_image_url"`
			} `json:"images"`
		} `json:"trailer"`
		Approved bool `json:"approved"`
		Titles   []struct {
			Type  string `json:"type"`
			Title string `json:"title"`
		} `json:"titles"`
		Title         string   `json:"title"`
		TitleEnglish  string   `json:"title_english"`
		TitleJapanese string   `json:"title_japanese"`
		TitleSynonyms []string `json:"title_synonyms"`
		Type          string   `json:"type"`
		Source        string   `json:"source"`
		Episodes      int      `json:"episodes"`
		Status        string   `json:"status"`
		Airing        bool     `json:"airing"`
		Aired         struct {
			From time.Time `json:"from"`
			To   time.Time `json:"to"`
			Prop struct {
				From struct {
					Day   int `json:"day"`
					Month int `json:"month"`
					Year  int `json:"year"`
				} `json:"from"`
				To struct {
					Day   int `json:"day"`
					Month int `json:"month"`
					Year  int `json:"year"`
				} `json:"to"`
			} `json:"prop"`
			String string `json:"string"`
		} `json:"aired"`
		Duration   string  `json:"duration"`
		Rating     string  `json:"rating"`
		Score      float64 `json:"score"`
		ScoredBy   int     `json:"scored_by"`
		Rank       int     `json:"rank"`
		Popularity int     `json:"popularity"`
		Members    int     `json:"members"`
		Favorites  int     `json:"favorites"`
		Synopsis   string  `json:"synopsis"`
		Background any     `json:"background"`
		Season     string  `json:"season"`
		Year       int     `json:"year"`
		Broadcast  struct {
			Day      string `json:"day"`
			Time     string `json:"time"`
			Timezone string `json:"timezone"`
			String   string `json:"string"`
		} `json:"broadcast"`
		Producers []struct {
			MalID int    `json:"mal_id"`
			Type  string `json:"type"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		} `json:"producers"`
		Licensors []struct {
			MalID int    `json:"mal_id"`
			Type  string `json:"type"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		} `json:"licensors"`
		Studios []struct {
			MalID int    `json:"mal_id"`
			Type  string `json:"type"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		} `json:"studios"`
		Genres []struct {
			MalID int    `json:"mal_id"`
			Type  string `json:"type"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		} `json:"genres"`
		ExplicitGenres []any `json:"explicit_genres"`
		Themes         []struct {
			MalID int    `json:"mal_id"`
			Type  string `json:"type"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		} `json:"themes"`
		Demographics []any `json:"demographics"`
	} `json:"data"`
}

// takes in a query string and searches for animes with the Jixen open API
func SearchAnime(q string) (AnimeSearchResults, error) {
	if len(q) < 3 {
		return AnimeSearchResults{}, fmt.Errorf("search should be more than 3 characters")
	}

	// put in util
	jixenURL := os.Getenv("JIKAN_BASE_URL")

	if jixenURL == "" {
		return AnimeSearchResults{}, fmt.Errorf("failed to load environment variables: %s", jixenURL)
	}

	resp, err := http.Get(jixenURL + "anime?q=" + q + "?limit=10")
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
