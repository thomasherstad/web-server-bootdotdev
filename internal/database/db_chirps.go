package database

import (
	"errors"
	"sort"
)

type Chirp struct {
	Id       int    `json:"id"`
	Body     string `json:"body"`
	AuthorID int    `json:"author_id"`
}

// CreateChirp creates a new chirp and saves it to the disk
func (db *DB) CreateChirp(body string, userID int) (Chirp, error) {
	var data Chirp
	data.Body = body
	data.AuthorID = userID

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

func (db *DB) GetChirpByID(id int) (Chirp, error) {
	database, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := database.Chirps[id]
	if !ok {
		return Chirp{}, errors.New("id not found")
	}

	return chirp, nil
}

func (db *DB) DeleteChirpByID(id int) error {
	database, err := db.loadDB()
	if err != nil {
		return err
	}

	//check if it exists
	_, ok := database.Chirps[id]
	if !ok {
		return errors.New("id not found")
	}

	delete(database.Chirps, id)

	db.writeDB(database)

	return nil

}
