package storage

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func (s *Service) InsertNewAnime(malID int, name string, dubVote int, subVote int) (id *int, err error) {
	err = s.DB.QueryRow(`
		INSERT INTO anime (name, dub_vote, sub_vote, created_at, updated_at, mal_id)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, $4)
		RETURNING id
	`, name, dubVote, subVote, malID).Scan(&id)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	return id, nil
}

func (s *Service) GetAnimeVoteDataByMalID(malID int) (anime *AnimeVoteData, err error) {
	anime = &AnimeVoteData{}

	err = s.DB.QueryRow(`
		SELECT id, name, dub_vote, sub_vote, mal_id, created_at, updated_at
		FROM anime
		WHERE mal_id = $1
	`, malID).Scan(
		&anime.ID,
		&anime.Name,
		&anime.DubVote,
		&anime.SubVote,
		&anime.MalID,
		&anime.CreatedAt,
		&anime.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching anime vote data: %v", err)
	}

	return anime, nil
}

func (s *Service) GetAnimeVoteDataByID(id int) (anime *AnimeVoteData, err error) {
	anime = &AnimeVoteData{}

	err = s.DB.QueryRow(`
		SELECT id, name, dub_vote, sub_vote, mal_id, created_at, updated_at
		FROM anime
		WHERE id = $1
	`, id).Scan(
		&anime.ID,
		&anime.Name,
		&anime.DubVote,
		&anime.SubVote,
		&anime.MalID,
		&anime.CreatedAt,
		&anime.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching anime vote data: %v", err)
	}

	return anime, nil
}

func (s *Service) AddVoteSubByID(id int) error {
	result, err := s.DB.Exec(`
		UPDATE anime
		SET sub_vote = sub_vote + 1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("error updating anime dub vote data: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no anime found with id %d", id)
	}

	return nil
}

func (s *Service) AddVoteDubByID(id int) error {
	result, err := s.DB.Exec(`
		UPDATE anime
		SET dub_vote = dub_vote + 1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("error updating anime dub vote data: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no anime found with id %d", id)
	}

	return nil
}

// integration tests

func (s *Service) DeleteAnimeByMalID(malID int) error {
	_, err := s.DB.Exec("DELETE from anime WHERE mal_id = $1", malID)
	if err != nil {
		return fmt.Errorf("error updating deleting anime: %v", err)
	}

	return nil
}
