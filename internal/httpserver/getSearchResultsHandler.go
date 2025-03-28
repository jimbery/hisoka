package httpserver

import (
	"encoding/json"
	"fmt"
	"hisoka/internal/httpclient"
	"net/http"
)

func getSearchResults(w http.ResponseWriter, r *http.Request) {
	enableCors(w, r)

	fmt.Println("Received /search request")

	searchTerm := r.FormValue("q")
	if searchTerm == "" {
		http.Error(w, "Missing search query", http.StatusBadRequest)
		return
	}

	results, err := httpclient.SearchAnime(searchTerm)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error searching for anime: %v", err), http.StatusInternalServerError)
		return
	}

	jsonResp, err := json.Marshal(results)
	if err != nil {
		http.Error(w, "Error marshalling JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResp)
	if err != nil {
		http.Error(w, "Error writing to output", http.StatusInternalServerError)
	}
}
