package main

import (
	"encoding/json"
	"log"
	"net/http"
	"web-server-bootdotdev/internal/auth"
)

func (cfg *apiConfig) HandlerUserUpdate(w http.ResponseWriter, r *http.Request) {
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
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters")
	}

	//Get bearer jwt
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	//Parse token to get the user id
	userID, err := auth.ParseJWTToken(tokenString, cfg.jwtSecret)
	if err != nil {
		log.Printf("Problem parsing JWT token. Error: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Token is invalid")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash new password")
		return
	}

	user, err := cfg.DB.UpdateUser(userID, params.Email, hashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user info")
		return
	}

	respondWithJson(w, http.StatusOK, response{
		User: User{
			ID:          user.ID,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
	})
}
