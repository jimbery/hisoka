package httpclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func GetAnime(MalID int) (AnimeDetailsFull, error) {
	jikanURL := os.Getenv("JIKAN_BASE_URL")
	if jikanURL == "" {
		return AnimeDetailsFull{}, errors.New("JIKAN_BASE_URL environment variable is not set")
	}

	MalIDString := strconv.Itoa(MalID)
	resp, err := http.Get(jikanURL + "anime/" + MalIDString + "/full")
	if err != nil {
		return AnimeDetailsFull{}, fmt.Errorf("error connecting to HTTP client: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AnimeDetailsFull{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return AnimeDetailsFull{}, fmt.Errorf("error reading response body: %w", err)
	}

	var result GetAnimeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return AnimeDetailsFull{}, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	genres := []GenreInfo{}

	for _, g := range result.Data.Genres {
		gen := GenreInfo{
			Name: g.Name,
		}
		genres = append(genres, gen)
	}

	streaming := []Streaming{}
	streaming = append(streaming, result.Data.Streaming...)

	output := AnimeDetailsFull{
		MalID:     result.Data.MalID,
		Image:     result.Data.Images.Jpg.LargeImageURL,
		Title:     result.Data.Title,
		Synopsis:  result.Data.Synopsis,
		Year:      result.Data.Year,
		Episodes:  result.Data.Episodes,
		Rating:    result.Data.Rating,
		Genres:    genres,
		Streaming: streaming,
		Trailer:   result.Data.Trailer.EmbedURL,
	}

	return output, nil
}
