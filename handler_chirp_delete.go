package main

import (
	"log"
	"net/http"
	"strconv"
	"web-server-bootdotdev/internal/auth"
)

func (cfg *apiConfig) handlerDeleteChirpById(w http.ResponseWriter, r *http.Request) {
	chirpID, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		log.Printf("Problem casting query param to int. Error: %v", err)
		respondWithError(w, http.StatusBadRequest, "")
		return
	}

	//Get userID from token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No token in header")
		return
	}

	userID, err := auth.ParseJWTToken(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	chirp, err := cfg.DB.GetChirpByID(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Chirp does not exist")
		return
	}

	//Check if userID is the same as the chirp's authorID
	if userID != chirp.AuthorID {
		respondWithError(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
		return
	}

	//deleteChirp()
	cfg.DB.DeleteChirpByID(chirpID)

	w.WriteHeader(http.StatusNoContent)
}
