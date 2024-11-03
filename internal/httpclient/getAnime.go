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

type AnimeDetailsFull struct {
	Title     string
	MalID     int
	Image     string
	Synopsis  string
	Year      int
	Episodes  int
	Rating    string
	Trailer   string
	Genres    []GenreInfo
	Streaming []Streaming
}

type GenreInfo struct {
	Type string
	Name string
	URL  string
}

type SteamingOutput struct {
	Name string
	URL  string
}

type GetAnimeResponse struct {
	Data Data `json:"data"`
}

type Data struct {
	MalID          int             `json:"mal_id"`
	URL            string          `json:"url"`
	Images         Images          `json:"images"`
	Trailer        Trailer         `json:"trailer"`
	Approved       bool            `json:"approved"`
	Titles         []Title         `json:"titles"`
	Title          string          `json:"title"`
	TitleEnglish   string          `json:"title_english"`
	TitleJapanese  string          `json:"title_japanese"`
	TitleSynonyms  []string        `json:"title_synonyms"`
	Type           string          `json:"type"`
	Source         string          `json:"source"`
	Episodes       int             `json:"episodes"`
	Status         string          `json:"status"`
	Airing         bool            `json:"airing"`
	Aired          Aired           `json:"aired"`
	Duration       string          `json:"duration"`
	Rating         string          `json:"rating"`
	Score          float64         `json:"score"`
	ScoredBy       int             `json:"scored_by"`
	Rank           int             `json:"rank"`
	Popularity     int             `json:"popularity"`
	Members        int             `json:"members"`
	Favorites      int             `json:"favorites"`
	Synopsis       string          `json:"synopsis"`
	Background     string          `json:"background"`
	Season         string          `json:"season"`
	Year           int             `json:"year"`
	Broadcast      Broadcast       `json:"broadcast"`
	Producers      []Producer      `json:"producers"`
	Licensors      []Licensor      `json:"licensors"`
	Studios        []Studio        `json:"studios"`
	Genres         []Genre         `json:"genres"`
	ExplicitGenres []ExplicitGenre `json:"explicit_genres"`
	Themes         []Theme         `json:"themes"`
	Demographics   []Demographic   `json:"demographics"`
	Relations      []Relation      `json:"relations"`
	Theme          ThemeDetails    `json:"theme"`
	External       []External      `json:"external"`
	Streaming      []Streaming     `json:"streaming"`
}

type Images struct {
	Jpg  ImageDetails `json:"jpg"`
	Webp ImageDetails `json:"webp"`
}

type ImageDetails struct {
	ImageURL      string `json:"image_url"`
	SmallImageURL string `json:"small_image_url"`
	LargeImageURL string `json:"large_image_url"`
}

type Trailer struct {
	YoutubeID string `json:"youtube_id"`
	URL       string `json:"url"`
	EmbedURL  string `json:"embed_url"`
}

type Title struct {
	Type  string `json:"type"`
	Title string `json:"title"`
}

type Aired struct {
	From string `json:"from"`
	To   string `json:"to"`
	Prop Prop   `json:"prop"`
}

type Prop struct {
	From   Date   `json:"from"`
	To     Date   `json:"to"`
	String string `json:"string"`
}

type Date struct {
	Day   int `json:"day"`
	Month int `json:"month"`
	Year  int `json:"year"`
}

type Broadcast struct {
	Day      string `json:"day"`
	Time     string `json:"time"`
	Timezone string `json:"timezone"`
	String   string `json:"string"`
}

type Producer struct {
	MalID int    `json:"mal_id"`
	Type  string `json:"type"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

type Licensor struct {
	MalID int    `json:"mal_id"`
	Type  string `json:"type"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

type Studio struct {
	MalID int    `json:"mal_id"`
	Type  string `json:"type"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

type Genre struct {
	MalID int    `json:"mal_id"`
	Type  string `json:"type"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

type ExplicitGenre struct {
	MalID int    `json:"mal_id"`
	Type  string `json:"type"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

type Theme struct {
	MalID int    `json:"mal_id"`
	Type  string `json:"type"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

type Demographic struct {
	MalID int    `json:"mal_id"`
	Type  string `json:"type"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

type Relation struct {
	Relation string  `json:"relation"`
	Entry    []Entry `json:"entry"`
}

type Entry struct {
	MalID int    `json:"mal_id"`
	Type  string `json:"type"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

type ThemeDetails struct {
	Openings []string `json:"openings"`
	Endings  []string `json:"endings"`
}

type External struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Streaming struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

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
	defer resp.Body.Close() // Ensure the response body is closed

	// Check for non-200 status codes
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
