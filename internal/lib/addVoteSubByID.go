package lib

import (
	"hisoka/internal/storage"
	"log"
)

func AddVoteSubByID(dbx *storage.Service, id int) (*storage.AnimeVoteData, error) {
	animeVoteData, err := dbx.GetAnimeVoteDataByID(id)
	if err != nil {
		log.Println("error getting GetAnimeVoteDataByMalId", err)
		return nil, err
	}

	if animeVoteData == nil {
		log.Println("anime does not exist", err)

		return nil, err
	}

	err = dbx.AddVoteSubByID(animeVoteData.ID)
	if err != nil {
		log.Println("error adding AddVoteSubById", err)
		return nil, err
	}

	animeVoteDataOutput, err := dbx.GetAnimeVoteDataByID(id)
	if err != nil {
		log.Println("error getting GetAnimeVoteDataByMalId", err)
		return nil, err
	}

	return animeVoteDataOutput, nil
}
