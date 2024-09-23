package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"web-server-bootdotdev/internal/auth"
)

type Chirp struct {
	ID       int    `json:"id"`
	Body     string `json:"body"`
	AuthorID int    `json:"author_id"`
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.DB.GetChirp()
	if err != nil {
		log.Printf("Problem getting Chirps from database. Error: %v\n", err)
		return
	}

	respondWithJson(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerCreateChirps(w http.ResponseWriter, r *http.Request) {
	//Each post request also needs a token in the authorization header
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	//Get user id from token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No jwt token")
		return
	}
	userID, err := auth.ParseJWTToken(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't parse jwt token")
		return
	}

	isValid, msg := isValidChirp(params.Body)

	if !isValid {
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	params.Body = silenceProfanities(params.Body)
	log.Printf("Incoming chirp: %s\n", string(params.Body))

	chirp, err := cfg.DB.CreateChirp(params.Body, userID)
	if err != nil {
		log.Printf("Problem creating chirp. Error: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create Chirp")
		return
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

func (cfg *apiConfig) handlerGetChirpById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		log.Printf("Problem casting query param to int. Error: %v", err)
		respondWithError(w, http.StatusBadRequest, "")
		return
	}

	allChirps, err := cfg.DB.GetChirp()
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
