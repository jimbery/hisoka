package storage

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq" // Import for side effects to register the driver
)

type Service struct {
	DB *sql.DB
}

func NewDBStore() (*Service, error) {
	// Set up the database connection
	connStr := os.Getenv("db")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	return &Service{DB: db}, nil
}

func (s *Service) Close() error {
	return s.DB.Close()
}

func (s *Service) InsertNewAnime(malID int, name string, dubVote int, subVote int) error {
	_, err := s.DB.Exec(`
		INSERT INTO anime (name, dub_vote, sub_vote, created_at, updated_at, mal_id)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, $4)
	`, name, dubVote, subVote, malID)

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *Service) GetAnimeVoteDataByMalID(malID int) (*AnimeVoteData, error) {
	var anime AnimeVoteData

	err := s.DB.QueryRow(`
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
			// Return nil if no rows are found
			return nil, nil
		}
		// Log other errors
		log.Printf("Error fetching anime vote data: %v\n", err)
		return nil, err
	}

	return &anime, nil
}

func (s *Service) GetAnimeVoteDataByID(id int) (*AnimeVoteData, error) {
	var anime AnimeVoteData

	err := s.DB.QueryRow(`
		SELECT id, name, dub_vote, sub_vote, mal_id, created_at, updated_at
		FROM anime
		WHERE mal_id = $1
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
			// Return nil if no rows are found
			return nil, nil
		}
		// Log other errors
		log.Printf("Error fetching anime vote data: %v\n", err)
		return nil, err
	}

	return &anime, nil
}

func (s *Service) AddVoteSubByID(id int) error {
	_, err := s.DB.Exec(`
		UPDATE anime
		SET sub_vote = sub_vote + 1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`, id)
	if err != nil {
		// Log other errors
		log.Printf("Error updating anime sub vote data: %v\n", err)
		return err
	}

	return nil
}

func (s *Service) AddVoteDubByID(id int) error {
	_, err := s.DB.Exec(`
		UPDATE anime
		SET dub_vote = dub_vote + 1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`, id)
	if err != nil {
		// Log other errors
		log.Printf("Error updating anime sub vote data: %v\n", err)
		return err
	}

	return nil
}
