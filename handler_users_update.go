package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AlexSkr96/Chirpy/internal/auth"
	"github.com/AlexSkr96/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
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
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't decode parameters: %v"), err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, fmt.Sprintf("Couldn't get an access token: %v", err), 401)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		http.Error(w, fmt.Sprintf("Couldn't validate JWT: %v", err), 401)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("Couldn't hash password: %v", err), 500)
		return
	}

	user, err := cfg.db.UpdateUserPasswordAndEmail(r.Context(), database.UpdateUserPasswordAndEmailParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
		ID:             userID,
	})

	respondWithJSON(w, 200, response{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
	})
}
