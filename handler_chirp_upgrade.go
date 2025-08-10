package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AlexSkr96/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUserChirpyRedUpgrade(w http.ResponseWriter, r *http.Request) {
	type data struct {
		UserID uuid.UUID `json:"user_id"`
	}
	type parameters struct {
		Event string `json:"event"`
		Data  data   `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get API key: %v", err), 401)
		return
	}
	if apiKey != cfg.polkaKey {
		http.Error(w, fmt.Sprintf("Invalid API key: %v", apiKey), 401)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	switch params.Event {
	case "user.upgraded":
		userID := params.Data.UserID
		err = cfg.db.UserSetIsChirpyRedTrueByID(r.Context(), userID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to upgrade user: %v", err), 404)
		}
	default:
	}
	respondWithJSON(w, 204, nil)
}
