package main

import (
	"encoding/json"
	"log"
	"net/http"
	"web-server-bootdotdev/internal/auth"
)

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	correctPassword := auth.ComparePasswords(user.Password, params.Password)
	if !correctPassword {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	// The jwt should expire in 1 hour
	token, err := auth.CreateJWTToken(user.ID, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create access token")
		return
	}

	//Generate refresh token
	refreshToken, expiration, err := auth.GenerateRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create refresh token")
	}
	log.Printf("Generated refresh token: %s", refreshToken)

	//add expiration token to database
	userTokenized, err := cfg.DB.AddUserRefreshToken(user.ID, refreshToken, expiration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not save token")
	}

	respondWithJson(w, http.StatusOK, response{
		User: User{
			ID:    userTokenized.ID,
			Email: userTokenized.Email,
		},
		Token:        token,
		RefreshToken: refreshToken,
	})
}
