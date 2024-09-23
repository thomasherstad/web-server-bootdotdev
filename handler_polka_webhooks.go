package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"web-server-bootdotdev/internal/auth"
	"web-server-bootdotdev/internal/database"
)

func (cfg *apiConfig) HandlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if apiKey != cfg.polkaApiKey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode json")
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	//Upgrade user
	err = cfg.DB.UpgradeUser(params.Data.UserID)
	if err != nil {
		if errors.Is(err, database.ErrNotExists) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user")
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
