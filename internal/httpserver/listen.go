package httpserver

import (
	"encoding/json"
	"fmt"
	"hisoka/internal/httpclient"
	"net/http"

	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(0.5, 5) // Limit to 2 requests per second with a burst of 5

func Listen() {
	http.Handle("/search", rateLimit(http.HandlerFunc(getSearchResults)))

	fmt.Println("Server is listening on port 3333...")
	err := http.ListenAndServe(":3333", nil)
	if err != nil {
		fmt.Println("error starting server:", err)
	}
}

// Middleware to handle rate limiting
func rateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r) // Call the next handler
	})
}

func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func getSearchResults(w http.ResponseWriter, r *http.Request) {
	enableCors(w) // Enable CORS at the start

	fmt.Printf("Received /search request\n")

	searchTerm := r.FormValue("q")
	if searchTerm == "" {
		http.Error(w, "Missing search query", http.StatusBadRequest)
		return
	}

	// Call the search function directly, without a goroutine
	results, err := httpclient.SearchAnime(searchTerm, &httpclient.RealHTTPClient{})
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
