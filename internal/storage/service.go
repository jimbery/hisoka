package storage

import (
	"database/sql"
	"log"
	"os"
	"time"
)

type Service struct {
	DB *sql.DB
}

func NewDBStore() (*Service, error) {
	connStr := os.Getenv("db")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open database connection:", err)
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(2 * time.Minute)

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping the database:", err)
		return nil, err
	}

	return &Service{DB: db}, nil
}

func (s *Service) Close() error {
	return s.DB.Close()
}
