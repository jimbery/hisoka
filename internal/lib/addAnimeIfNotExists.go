package lib

import (
	"hisoka/internal/httpclient"
	"hisoka/internal/storage"
	"log"
)

func AddAnimeIfNotExists(malID int) error {
	dbx, _ := storage.NewDBStore()

	defer func() {
		if err := dbx.Close(); err != nil {
			log.Fatal("Failed to close the database:", err)
		}
	}()

	animeVoteData, err := dbx.GetAnimeVoteDataByMalID(malID)
	if err != nil {
		log.Println("error getting GetAnimeVoteDataByMalId", err)
		return err
	}

	if animeVoteData == nil {
		anime, err := httpclient.GetAnime(malID)
		if err != nil {
			log.Println("error getting anime from Jixen", err)
			return err
		}

		err = dbx.InsertNewAnime(malID, anime.Title, 0, 0)
		if err != nil {
			log.Println("error InsertNewAnime", err)
			return err
		}
	}

	return nil
}
