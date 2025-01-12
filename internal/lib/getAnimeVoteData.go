package lib

import (
	"fmt"
	"hisoka/internal/storage"
	"log"
)

func GetAnimnVoteDataByMalID(malID int) (*storage.AnimeVoteData, error) {
	dbx, _ := storage.NewDBStore()

	defer func() {
		if err := dbx.Close(); err != nil {
			log.Fatal("Failed to close the database:", err)
		}
	}()

	animeVoteData, err := dbx.GetAnimeVoteDataByMalID(malID)
	if err != nil {
		log.Println("error getting GetAnimeVoteDataByMalId", err)
		return nil, err
	}

	if animeVoteData == nil {
		log.Println("anime does not exist", err)

		return nil, fmt.Errorf("anime does not exist")
	}

	return animeVoteData, nil
}
