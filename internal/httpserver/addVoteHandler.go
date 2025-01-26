package httpserver

import (
	"encoding/json"
	"hisoka/internal/lib"
	"log"
	"net/http"
)

func addVote(w http.ResponseWriter, r *http.Request) {
	enableCors(w, r) // Enable CORS at the start
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var input voteInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		log.Println("Error decoding JSON:", err)
		return
	}

	// Ensure the anime exists
	id, err := lib.AddAnimeIfNotExists(input.MalID)
	if err != nil {
		http.Error(w, "Failed to add anime", http.StatusBadRequest)
		log.Println("Error adding anime:", err)
		return
	}

	log.Printf("Received vote for mal_id %d with vote type %s\n", input.MalID, input.VoteType)

	var animeVoteDataOutput interface{}
	switch input.VoteType {
	case DubVote:
		animeVoteDataOutput, err = lib.AddVoteDubByID(*id)
	case SubVote:
		animeVoteDataOutput, err = lib.AddVoteSubByID(*id)
	default:
		http.Error(w, "Invalid vote type", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Failed to add vote", http.StatusInternalServerError)
		log.Println("Error adding vote:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(animeVoteDataOutput)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Println("Error encoding response:", err)
	}
}
