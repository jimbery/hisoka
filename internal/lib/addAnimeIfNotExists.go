package lib

import (
	"fmt"
	"hisoka/internal/httpclient"
	"hisoka/internal/storage"
	"log"
)

func AddAnimeIfNotExists(dbx *storage.Service, malID int) (id *int, err error) {
	animeVoteData, err := dbx.GetAnimeVoteDataByMalID(malID)
	if err != nil {
		log.Println("error getting GetAnimeVoteDataByMalID", err)
		return nil, fmt.Errorf("GetAnimeVoteDataByMalID %s", err)
	}

	if animeVoteData == nil {
		anime, err := httpclient.GetAnime(malID)
		if err != nil {
			log.Println("error getting anime from Jixen", err)
			return nil, fmt.Errorf("GetAnime %s", err)
		}

		id, err = dbx.InsertNewAnime(malID, anime.Title, 0, 0)
		if err != nil {
			log.Println("error InsertNewAnime", err)
			return nil, fmt.Errorf("InsertNewAnime %s", err)
		}

		return id, nil
	}

	return &animeVoteData.ID, nil
}
