package main

import (
	"encoding/json"
	"net/http"
)

func handlerPostChirps(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := Chirp{}
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
	respondWithJson(w, http.StatusCreated, params)
}

func isValidChirp(text string) (bool, string) {
	const maxLength = 140
	if len(text) > maxLength {
		return false, "Chirp is too long"
	}
	return true, ""
}
