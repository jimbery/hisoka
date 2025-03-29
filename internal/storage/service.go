package storage

import (
	"database/sql"
	"fmt"
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
		log.Fatal("failed to open database connection:", err)
		return nil, fmt.Errorf("failed to open db connection %s", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(2 * time.Minute)

	err = db.Ping()
	if err != nil {
		log.Fatal("failed to ping the database:", err)
		return nil, fmt.Errorf("fail to ping db %s", err)
	}

	return &Service{DB: db}, nil
}

func (s *Service) Close() error {
	return s.DB.Close()
}
