package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"web-server-bootdotdev/internal/database"

	"golang.org/x/crypto/bcrypt"
)

func handlerPostUsers(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := database.User{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	log.Printf("Incoming user: %s\n", string(params.Email))

	// add to db
	db, err := database.NewDB(dbPath)
	if err != nil {
		log.Printf("error creating database: %v\n", err)
		//respond with error
		return
	}

	fmt.Printf("Pre-hashed: %s\n", params.Password)
	hashedPassword, err := hashPassword(params.Password)
	if err != nil {
		log.Printf("Error when hashing the password: %v", err)
		//respondwitherror
		return
	}
	fmt.Printf("Post-hashed: %s\n", hashedPassword)

	usr, err := db.CreateUser(params.Email, hashedPassword)
	if err != nil {
		log.Printf("Problem creating user. Error: %v", err)
		//respond with error
		return
	}

	log.Printf("Outgoing user: %s\n", usr.Email)

	respondWithJson(w, http.StatusCreated, struct {
		Id    int
		Email string
	}{
		Id:    usr.Id,
		Email: usr.Email,
	})
}

func hashPassword(p string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
