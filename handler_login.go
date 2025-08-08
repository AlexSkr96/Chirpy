package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AlexSkr96/Chirpy/internal/auth"
	"github.com/AlexSkr96/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(context.Background(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("No user with email %v", params.Email), err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret)

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
		return
	}
	cfg.db.SaveRefreshToken(context.Background(), database.SaveRefreshTokenParams{
		UserID: user.ID,
		Token:  refreshToken,
	})

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:           user.ID,
			Email:        user.Email,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
			Token:        token,
			RefreshToken: refreshToken,
		},
	})
}
