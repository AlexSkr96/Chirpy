package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/AlexSkr96/Chirpy/internal/database"
)

func createChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parametrs struct {
		Body   string `json:"body"`
		UserID string `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parametrs{}
	err := decoder.Decode(&params)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if len(params.Body) > 140 {
		http.Error(w, "Chirp can't be longger then 140 symbols", 400)
		return
	}

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleanedOutput := getCleanedBody(params.Body, badWords)
	userID, err := uuid.Parse(params.UserID)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}
	chirp, err := cfg.DB.CreateChirp(context.Background(), database.CreateChirpParams{UserID: userID, Body: cleanedOutput})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(chirp)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}
