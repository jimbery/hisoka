package httpserver

import (
	"fmt"
	"net/http"

	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(0.5, 5) // Limit to 2 requests per second with a burst of 5

func Listen() {
	// Handle specific routes
	http.Handle("/api/search", rateLimit(http.HandlerFunc(getSearchResults)))
	http.Handle("/api/anime/", rateLimit(http.HandlerFunc(getAnime)))
	http.Handle("/api/add_vote", rateLimit(http.HandlerFunc(addVote)))

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// Write an HTTP 200 OK status
		w.WriteHeader(http.StatusOK)

		fmt.Println("received health check")

		// Send a response body
		_, err := w.Write([]byte("OK"))
		if err != nil {
			http.Error(w, "Error writing to output", http.StatusInternalServerError)
		}
	})

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

func enableCors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Handle OPTIONS pre-flight request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
}
