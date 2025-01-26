package httpserver

import (
	"encoding/json"
	"fmt"
	"hisoka/internal/lib"
	"hisoka/internal/storage"
	"net/http"
	"strconv"
	"strings"
)

func getAnimeVoteData(w http.ResponseWriter, r *http.Request, dbx *storage.Service) {
	enableCors(w, r) // Enable CORS at the start
	fmt.Println("Received /anime-vote-data request")

	path := r.URL.Path

	parts := strings.Split(path, "/")

	if len(parts) > 3 {
		MalIDString := parts[3]

		MalID, err := strconv.Atoi(MalIDString)
		if err != nil {
			http.Error(w, fmt.Sprintf("anime ID must be a number: %v", err), http.StatusInternalServerError)
			return
		}

		results, err := lib.GetAnimnVoteDataByMalID(dbx, MalID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error searching for anime: %v", err), http.StatusNotFound)
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
