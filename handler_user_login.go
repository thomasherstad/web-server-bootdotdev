package main

import (
	"encoding/json"
	"net/http"
	"web-server-bootdotdev/internal/auth"
)

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds,omitempty"`
	}

	type response struct {
		User
		Token string `json:"token"`
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

	var expiresInSeconds int
	if params.ExpiresInSeconds != 0 {
		expiresInSeconds = params.ExpiresInSeconds
	} else {
		expiresInSeconds = 86400 //24 hours
	}

	token, err := auth.CreateJWTToken(expiresInSeconds, user.ID, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create token")
		return
	}

	respondWithJson(w, http.StatusOK, response{
		User: User{
			ID:    user.ID,
			Email: user.Email,
		},
		Token: token,
	})
}
