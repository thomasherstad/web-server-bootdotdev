package main

import (
	"net/http"
	"time"
	"web-server-bootdotdev/internal/auth"
)

func (cfg *apiConfig) HandlerUserRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No token in header")
		return
	}

	//Check if refresh token exists
	refToken, err := cfg.DB.GetRefreshToken(refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Refresh token doesn't exist")
		return
	}

	//Check if refresh token is valid
	expiration := refToken.Expiry
	now := time.Now()
	if now.After(expiration) {
		respondWithError(w, http.StatusUnauthorized, "Refresh token has expired")
		return
	}

	user, err := cfg.DB.GetUserById(refToken.UserId)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't fetch user from refresh token")
		return
	}

	//Create new access token
	newToken, err := auth.CreateJWTToken(user.ID, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create jwt token")
	}

	respondWithJson(w, http.StatusOK, response{
		Token: newToken,
	})

}
