package httpserver

import (
	"encoding/json"
	"fmt"
	"hisoka/internal/httpclient"
	"net/http"
)

func Listen() {
	http.HandleFunc("/search", getSearchResults)

	err := http.ListenAndServe(":3333", nil)
	if err != nil {
		fmt.Println("error starting server")
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func getSearchResults(w http.ResponseWriter, r *http.Request) {
	taskCompleted := make(chan bool)
	enableCors(&w)
	fmt.Printf("got /q request\n")
	var test httpclient.AnimeSearchResults

	searchTerm := r.FormValue("q")

	go func() {
		var err error
		test, err = httpclient.SearchAnime(searchTerm, &httpclient.RealHTTPClient{})
		if err != nil {

			http.Error(w, `Error searching for anime + $err`, http.StatusInternalServerError)
		}

		// Signal that the task has completed by sending true to the channel
		close(taskCompleted)
	}()

	// Wait for the goroutine to complete before proceeding
	<-taskCompleted

	jsonResp, _ := json.Marshal(test)
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(jsonResp)
	if err != nil {
		http.Error(w, "Error writng to output ", http.StatusInternalServerError)
	}
}
