package main

import (
	"encoding/json"
	"log"
	"net/http"
	"web-server-bootdotdev/internal/auth"
	"web-server-bootdotdev/internal/database"
)

// TODO: Don't allow multiple users with the same email
func handlerPostUsers(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := database.User{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	log.Printf("Incoming user: %s\n", string(params.Email))

	// connect to db
	db, err := database.NewDB(dbPath)
	if err != nil {
		log.Printf("error creating database: %v\n", err)
		//respond with error
		return
	}

	// Check if there is already a user with that email
	_, err = db.GetUserByEmail(params.Email)
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

	usr, err := db.CreateUser(params.Email, hashedPassword)
	if err != nil {
		log.Printf("Problem creating user. Error: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	log.Printf("Outgoing user: %s\n", usr.Email)

	respondWithJson(w, http.StatusCreated, struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	}{
		Id:    usr.Id,
		Email: usr.Email,
	})
}

func handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := database.User{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	db, err := database.NewDB(dbPath)
	if err != nil {
		log.Printf("error creating database: %v\n", err)
		//respond with error
		return
	}

	user, err := db.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	correctPassword := auth.ComparePasswords(user.Password, params.Password)
	if !correctPassword {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	respondWithJson(w, http.StatusOK, struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	}{
		Id:    user.Id,
		Email: user.Email,
	})
}
