package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"web-server-bootdotdev/internal/database"
)

func handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	db, err := database.NewDB("./database.json")
	if err != nil {
		log.Printf("Problem establishing connection to database. Error: %v", err)
		return
	}
	chirps, err := db.GetChirp()
	if err != nil {
		log.Printf("Problem getting Chirps from database. Error: %v\n", err)
		return
	}

	respondWithJson(w, http.StatusOK, chirps)
}

func handlerPostChirps(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := database.Chirp{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	isValid, msg := isValidChirp(params.Body)

	if !isValid {
		respondWithError(w, http.StatusBadRequest, msg)
	}

	params.Body = silenceProfanities(params.Body)
	log.Printf("Incoming chirp: %s\n", string(params.Body))

	// add to db
	db, err := database.NewDB("./database.json")
	if err != nil {
		log.Printf("error creating database: %v\n", err)
	}

	chirp, err := db.CreateChirp(params.Body)
	if err != nil {
		log.Printf("Problem creating chirp. Error: %v", err)
	}

	log.Printf("Outgoing chirp: %s\n", chirp.Body)

	respondWithJson(w, http.StatusCreated, chirp)

}

func isValidChirp(text string) (bool, string) {
	const maxLength = 140
	if len(text) > maxLength {
		return false, "Chirp is too long"
	}
	return true, ""
}

func handlerGetChirpById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		log.Printf("Problem casting query param to int. Error: %v", err)
		respondWithError(w, http.StatusBadRequest, "")
		return
	}

	db, err := database.NewDB(dbPath)
	if err != nil {
		log.Printf("Problem creating db connection. Error: %v", err)
		respondWithError(w, http.StatusInternalServerError, "")
		return
	}

	allChirps, err := db.GetChirp()
	if err != nil {
		log.Printf("Problem getting chirps. Error: %v", err)
		respondWithError(w, http.StatusInternalServerError, "")
		return
	}

	if id > len(allChirps) || id < 1 {
		log.Printf("Id out of range. Id: %v, length of db: %v", id, len(allChirps))
		respondWithError(w, http.StatusNotFound, "Chirp does not exist")
		return
	}

	respondWithJson(w, http.StatusOK, allChirps[id-1])
	log.Printf("ChirpByID request successfully sent")
}
