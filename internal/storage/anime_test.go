package storage_test

import (
	"database/sql"
	"hisoka/internal/storage"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testDB *sql.DB
	dbx    *storage.Service
)

func TestMain(m *testing.M) {
	var cleanup func()
	var err error
	os.Setenv("db", "postgres://postgres:password@localhost:5432/postgres?sslmode=disable")

	testDB, cleanup = setupTestDB()

	dbx, err = storage.NewDBStore()
	if err != nil {
		log.Fatalf("Failed to open storage connection: %v", err)
	}

	code := m.Run()

	cleanup()

	os.Exit(code)
}

func setupTestDB() (*sql.DB, func()) {
	connStr := "host=localhost port=5432 user=postgres password=password dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	sqlFile, err := os.Open("migrations/000001_create_anime_table.up.sql")
	if err != nil {
		log.Fatalf("Failed to open SQL file: %v", err)
	}
	defer sqlFile.Close()

	migration, err := io.ReadAll(sqlFile)
	if err != nil {
		log.Fatalf("Failed to read SQL file: %v", err)
	}

	sql := strings.ReplaceAll(string(migration), "\r\n", "\n")

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}

	_, err = db.Exec(sql)
	if err != nil {
		log.Fatalf("Failed to execute migration: %v", err)
	}

	cleanup := func() {
		db.Close()
	}

	return db, cleanup
}

func TestInsertNewAnime(t *testing.T) {
	malID := rand.Intn(100000)
	name := "test"
	dubVote := 3
	subVote := 2

	defer func() {
		err := dbx.DeleteAnimeByMalID(malID)
		require.NoError(t, err)
	}()

	_, err := dbx.InsertNewAnime(malID, name, dubVote, subVote)
	require.NoError(t, err)

	voteData, err := dbx.GetAnimeVoteDataByMalID(malID)
	require.NoError(t, err)

	assert.Equal(t, malID, voteData.MalID, "mal_id is incorrect in db")
	assert.Equal(t, name, voteData.Name, "name is inconnect in db")
	assert.Equal(t, dubVote, voteData.DubVote, "dub_vote is incorrect in db")
	assert.Equal(t, subVote, voteData.SubVote, "sub_vote is incorrect in db")

}

func TestGetAnimeVoteDataByMalID(t *testing.T) {
	malID := rand.Intn(100000)
	name := "test"
	dubVote := 3
	subVote := 2

	defer func() {
		err := dbx.DeleteAnimeByMalID(malID)
		require.NoError(t, err)
	}()

	_, err := dbx.InsertNewAnime(malID, name, dubVote, subVote)
	require.NoError(t, err)

	voteData, err := dbx.GetAnimeVoteDataByMalID(malID)
	require.NoError(t, err)

	assert.Equal(t, voteData.DubVote, dubVote, "dub vote not expected value")
	assert.Equal(t, voteData.SubVote, subVote, "sub vote not expected value")
	assert.Equal(t, voteData.MalID, malID, "malID not expected value")
	assert.NotEmpty(t, voteData.CreatedAt)
	assert.NotEmpty(t, voteData.UpdatedAt)
	assert.NotEmpty(t, voteData.ID)
}

func TestGetAnimeVoteDataByMalIDNoData(t *testing.T) {
	malID := rand.Intn(100000)

	voteData, err := dbx.GetAnimeVoteDataByMalID(malID)
	require.NoError(t, err)

	assert.Nil(t, voteData, "expected anime vote data to be nil")
}

func TestAddVoteDubByID(t *testing.T) {
	malID := rand.Intn(100000)
	name := "test"
	dubVote := 0
	subVote := 0

	defer func() {
		err := dbx.DeleteAnimeByMalID(malID)
		require.NoError(t, err)
	}()

	id, err := dbx.InsertNewAnime(malID, name, dubVote, subVote)
	require.NoError(t, err)

	err = dbx.AddVoteDubByID(*id)
	require.NoError(t, err)

	voteData, err := dbx.GetAnimeVoteDataByMalID(malID)
	require.NoError(t, err)

	assert.Equal(t, dubVote+1, voteData.DubVote, "dubvote incorrect value")
	assert.Equal(t, subVote, voteData.SubVote, "subvote incorrect value")
}

func TestAddVoteDubByIDError(t *testing.T) {
	id := rand.Intn(100000)

	err := dbx.AddVoteDubByID(id)

	assert.ErrorContains(t, err, "no anime found with id")
}

func TestAddVoteSubByID(t *testing.T) {
	malID := rand.Intn(100000)
	name := "test"
	dubVote := 0
	subVote := 0

	defer func() {
		err := dbx.DeleteAnimeByMalID(malID)
		require.NoError(t, err)
	}()

	id, err := dbx.InsertNewAnime(malID, name, dubVote, subVote)
	require.NoError(t, err)

	err = dbx.AddVoteSubByID(*id)
	require.NoError(t, err)

	voteData, err := dbx.GetAnimeVoteDataByMalID(malID)
	require.NoError(t, err)

	assert.Equal(t, dubVote, voteData.DubVote, "dubvote incorrect value")
	assert.Equal(t, subVote+1, voteData.SubVote, "subvote incorrect value")
}

func TestAddVoteSubByIDError(t *testing.T) {
	id := rand.Intn(100000)

	err := dbx.AddVoteSubByID(id)

	assert.ErrorContains(t, err, "no anime found with id")
}
