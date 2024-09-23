package main

import (
	"encoding/json"
	"log"
	"net/http"
	"web-server-bootdotdev/internal/auth"
)

type User struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"-"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerPostUsers(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	// Check if there is already a user with that email
	_, err = cfg.DB.GetUserByEmail(params.Email)
	if err == nil {
		respondWithError(w, http.StatusConflict, "Email already exists.")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error when hashing the password: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	usr, err := cfg.DB.CreateUser(params.Email, hashedPassword)
	if err != nil {
		log.Printf("Problem creating user. Error: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	log.Printf("Outgoing user: %s\n", usr.Email)

	respondWithJson(w, http.StatusCreated, response{
		User: User{
			ID:          usr.ID,
			Email:       usr.Email,
			IsChirpyRed: usr.IsChirpyRed,
		},
	})
}
