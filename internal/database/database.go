package database

import (
	"encoding/json"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type Chirp struct {
	Id   int    `json:id`
	Body string `json:"body"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	db := DB{
		path: path,
	}

	err := db.ensureDB()
	if err != nil {
		return &DB{}, err
	}

	return &db, nil
}

// CreateChirp creates a new chirp and saves it to the disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	var data Chirp
	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		return Chirp{}, nil
	}
	database, err := db.readDatabase()
	if err != nil {
		return Chirp{}, err
	}

	data.Id = len(database.Chirps) + 1
	database.Chirps[data.Id] = data

	//Save to the disk
	jsonBytes, err := json.Marshal(database)
	if err != nil {
		return Chirp{}, err
	}

	os.WriteFile(db.path, jsonBytes, fs.Perm())

	return data, nil
}

// GetChirp returns all chirps in the database
func (db *DB) GetChirp() ([]Chirp, error) {
	database, err := db.readDatabase()
	if err != nil {
		return []Chirp{}, err
	}

	var chirps []Chirp
	for _, chirp := range database.Chirps {
		chirps = append(chirps, chirp)
	}
	return chirps, nil
}

// ensureDB creates a new database file if it doesn't exist
// TODO: RETURN ERRORS PROPERLY
func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if os.IsNotExist(err) {
		os.Create(db.path)
		return nil
	}
	return nil
}

func (db *DB) readDatabase() (DBStructure, error) {
	fileData, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}

	var database DBStructure
	err = json.Unmarshal(fileData, &database)
	if err != nil {
		return DBStructure{}, err
	}

	return database, nil
}
