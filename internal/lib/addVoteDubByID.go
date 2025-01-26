package lib

import (
	"hisoka/internal/storage"
	"log"
)

func AddVoteDubByID(dbx *storage.Service, id int) (*storage.AnimeVoteData, error) {
	animeVoteData, err := dbx.GetAnimeVoteDataByMalID(id)
	if err != nil {
		log.Println("error getting GetAnimeVoteDataByMalId", err)
		return nil, err
	}

	if animeVoteData == nil {
		log.Println("anime does not exist", err)

		return nil, err
	}

	err = dbx.AddVoteDubByID(animeVoteData.ID)
	if err != nil {
		log.Println("error adding AddVoteSubById", err)
		return nil, err
	}

	animeVoteDataOutput, err := dbx.GetAnimeVoteDataByMalID(id)
	if err != nil {
		log.Println("error getting GetAnimeVoteDataByMalId", err)
		return nil, err
	}

	return animeVoteDataOutput, nil
}
