package httpserver

import (
	"encoding/json"
	"fmt"
	"hisoka/internal/httpclient"
	"net/http"
	"strconv"
	"strings"
)

func getAnime(w http.ResponseWriter, r *http.Request) {
	enableCors(w, r)
	fmt.Println("Received /anime request")

	path := r.URL.Path

	parts := strings.Split(path, "/")

	if len(parts) > 3 {
		MalIDString := parts[3]

		MalID, err := strconv.Atoi(MalIDString)
		if err != nil {
			http.Error(w, fmt.Sprintf("anime ID must be a number: %v", err), http.StatusInternalServerError)
			return
		}

		results, err := httpclient.GetAnime(MalID)
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
	} else {
		fmt.Fprintf(w, "missing id from url")
	}

}
