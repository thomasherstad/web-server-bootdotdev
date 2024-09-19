package database

import (
	"encoding/json"
	"os"
	"sort"
	"sync"
)

//TODO:
// - use mutex when accessing the file
// - need to send correct status codes when errors pop up

type DB struct {
	path string
	mu   *sync.RWMutex
}

type Chirp struct {
	Id   int    `json:"id"`
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
		mu:   &sync.RWMutex{},
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
	data.Body = body

	database, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	data.Id = len(database.Chirps) + 1
	database.Chirps[data.Id] = data

	//Save to the disk
	err = db.writeDB(database)
	if err != nil {
		return Chirp{}, err
	}

	return data, nil
}

// GetChirp returns all chirps in the database
func (db *DB) GetChirp() ([]Chirp, error) {
	database, err := db.loadDB()
	if err != nil {
		return []Chirp{}, err
	}

	var chirps []Chirp
	for _, chirp := range database.Chirps {
		chirps = append(chirps, chirp)
	}

	//TODO: Needs to be sorted

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].Id < chirps[j].Id
	})

	return chirps, nil
}

// ensureDB creates a new database file if it doesn't exist
// TODO: RETURN ERRORS PROPERLY
func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)

	database := DBStructure{
		Chirps: make(map[int]Chirp),
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
