package lib

import (
	"fmt"
	"hisoka/internal/storage"
	"log"
)

func AddVoteSubByID(dbx *storage.Service, id int) (*storage.AnimeVoteData, error) {
	animeVoteData, err := dbx.GetAnimeVoteDataByID(id)
	if err != nil {
		log.Println("error getting GetAnimeVoteDataByID", err)
		return nil, fmt.Errorf("GetAnimeVoteDataByID %s", err)
	}

	if animeVoteData == nil {
		log.Println("anime does not exist")
		return nil, fmt.Errorf("anime does not exist")
	}

	err = dbx.AddVoteSubByID(animeVoteData.ID)
	if err != nil {
		log.Println("error adding AddVoteSubById", err)
		return nil, fmt.Errorf("AddVoteSubById %s", err)
	}

	animeVoteDataOutput, err := dbx.GetAnimeVoteDataByID(id)
	if err != nil {
		log.Println("error getting GetAnimeVoteDataByID", err)
		return nil, fmt.Errorf("GetAnimeVoteDataByID %s", err)
	}

	return animeVoteDataOutput, nil
}
