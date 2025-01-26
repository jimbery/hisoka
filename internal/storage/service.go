package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

type Service struct {
	DB *sql.DB
}

func NewDBStore() (*Service, error) {
	connStr := os.Getenv("db")
	fmt.Println(connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return &Service{DB: db}, nil
}

func (s *Service) Close() error {
	return s.DB.Close()
}
