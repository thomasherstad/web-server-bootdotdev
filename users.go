package main

import (
	"encoding/json"
	"log"
	"net/http"
	"web-server-bootdotdev/internal/database"
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
	}

	usr, err := db.CreateUser(params.Email)
	if err != nil {
		log.Printf("Problem creating user. Error: %v", err)
	}

	log.Printf("Outgoing user: %s\n", usr.Email)

	respondWithJson(w, http.StatusCreated, usr)
}
