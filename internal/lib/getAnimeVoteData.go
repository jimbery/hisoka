package lib

import (
	"fmt"
	"hisoka/internal/storage"
	"log"
)

func GetAnimnVoteDataByMalID(dbx *storage.Service, malID int) (*storage.AnimeVoteData, error) {
	animeVoteData, err := dbx.GetAnimeVoteDataByMalID(malID)
	if err != nil {
		log.Println("error getting GetAnimeVoteDataByMalId", err)
		return nil, fmt.Errorf("GetAnimeVoteDataByMalId %s", err)
	}

	if animeVoteData == nil {
		log.Println("anime does not exist")
		return nil, fmt.Errorf("anime does not exist")
	}

	return animeVoteData, nil
}
