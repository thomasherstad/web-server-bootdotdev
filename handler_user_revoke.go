package main

import (
	"net/http"
	"web-server-bootdotdev/internal/auth"
)

func (cfg *apiConfig) HandlerUserRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No token in header")
		return
	}

	err = cfg.DB.DeleteRefreshToken(refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token does not exist")
	}

	w.WriteHeader(http.StatusNoContent)
}
