package database

import (
	"encoding/json"
	"os"
	"sync"
)

//TODO:
// - need to send correct status codes when errors pop up

type DB struct {
	path string
	mu   *sync.RWMutex
}

type DBStructure struct {
	Chirps        map[int]Chirp           `json:"chirps"`
	Users         map[int]User            `json:"users"`
	RefreshTokens map[string]RefreshToken `json:"refresh_tokens"` //refresh tokens -> user id
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	db := DB{
		path: path,
		mu:   &sync.RWMutex{},
	}

	err := db.ensureDB()
	if err != nil {
		return &DB{}, err
	}

	return &db, nil
}

// ensureDB creates a new database file if it doesn't exist
// TODO: RETURN ERRORS PROPERLY
func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)

	database := DBStructure{
		Chirps:        make(map[int]Chirp),
		Users:         make(map[int]User),
		RefreshTokens: make(map[string]RefreshToken),
	}

	if os.IsNotExist(err) {
		os.Create(db.path)
		db.writeDB(database)
		return nil
	}
	return nil
}

func (db *DB) loadDB() (DBStructure, error) {
	// db.mu.Lock()
	// defer db.mu.Unlock()

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

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	jsonBytes, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, jsonBytes, 0777)
	if err != nil {
		return err
	}
	return nil
}
