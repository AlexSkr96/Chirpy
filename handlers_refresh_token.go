package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/AlexSkr96/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerVerifyRefreshToken(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, "No refresh token found", http.StatusUnauthorized)
		return
	}

	refreshToken, err := cfg.db.GetRefreshToken(context.Background(), token)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}
	if refreshToken.RevokedAt.Valid {
		http.Error(w, "Refresh token revoked", http.StatusUnauthorized)
		return
	}
	if time.Now().After(refreshToken.ExpiresAt) {
		http.Error(w, "Refresh token expired", http.StatusUnauthorized)
		return
	}

	jwt, err := auth.MakeJWT(refreshToken.UserID, cfg.secret)
	if err != nil {
		http.Error(w, "Failed to create JWT", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, 200, response{
		Token: jwt,
	})
}

func (cfg *apiConfig) handlerRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, fmt.Sprintf("No refresh token found: %v", err), http.StatusNoContent)
		return
	}

	err = cfg.db.RevokeRefreshToken(context.Background(), token)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to revoke refresh token: %v", err), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, 204, struct{}{})
}
