package lib

import (
	"hisoka/internal/httpclient"
	"hisoka/internal/storage"
	"log"
)

func AddAnimeIfNotExists(malID int) (id *int, err error) {
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
		anime, err := httpclient.GetAnime(malID)
		if err != nil {
			log.Println("error getting anime from Jixen", err)
			return nil, err
		}

		id, err = dbx.InsertNewAnime(malID, anime.Title, 0, 0)
		if err != nil {
			log.Println("error InsertNewAnime", err)
			return nil, err
		}
	}

	return id, nil
}
