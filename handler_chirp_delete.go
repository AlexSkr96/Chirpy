package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/AlexSkr96/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid chirp ID: %v", err.Error()), http.StatusBadRequest)
		return
	}

	chirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting chirp: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error in handlerChirpDelete while getting bearer token: %v", err.Error()), http.StatusUnauthorized)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error in handlerChirpDelete while validating JWT: %v", err.Error()), http.StatusUnauthorized)
		return
	}

	if chirp.UserID == userID {
		err = cfg.db.DeleteChirpByID(r.Context(), chirpID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error deleting chirp: %v", err.Error()), 404)
			return
		}
		respondWithJSON(w, 204, nil)
	} else {
		http.Error(w, "Unauthorized", 403)
		return
	}
}
