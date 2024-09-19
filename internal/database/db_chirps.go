package database

import "sort"

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
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
