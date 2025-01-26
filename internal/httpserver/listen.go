package httpserver

import (
	"fmt"
	"hisoka/internal/storage"
	"log"
	"net/http"

	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(0.5, 5) // Limit to 2 requests per second with a burst of 5
func Listen() {
	service, err := storage.NewDBStore()
	if err != nil {
		log.Fatalf("Could not create service: %v", err)
	}
	defer service.DB.Close()

	http.Handle("/api/search", rateLimit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		getSearchResults(w, r)
	})))
	http.Handle("/api/anime/", rateLimit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		getAnime(w, r)
	})))
	http.Handle("/api/add-vote", rateLimit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		addVote(w, r, service)
	})))
	http.Handle("/api/anime-vote-data/", rateLimit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		getAnimeVoteData(w, r, service)
	})))

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Println("received health check")
		_, err := w.Write([]byte("OK"))
		if err != nil {
			http.Error(w, "Error writing to output", http.StatusInternalServerError)
		}
	})

	fmt.Println("Server is listening on port 3333...")
	err = http.ListenAndServe(":3333", nil)
	if err != nil {
		fmt.Println("error starting server:", err)
	}
}

func rateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func enableCors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
}
