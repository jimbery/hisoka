package lib

import (
	"hisoka/internal/storage"
	"log"
)

func AddVoteSubByMalID(malID int) (*storage.AnimeVoteData, error) {
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

		return nil, err
	}

	err = dbx.AddVoteSubByID(animeVoteData.ID)
	if err != nil {
		log.Println("error adding AddVoteSubById", err)
		return nil, err
	}

	animeVoteDataOutput, err := dbx.GetAnimeVoteDataByMalID(malID)
	if err != nil {
		log.Println("error getting GetAnimeVoteDataByMalId", err)
		return nil, err
	}

	return animeVoteDataOutput, nil
}
