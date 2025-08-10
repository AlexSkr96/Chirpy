package main

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/google/uuid"

	"github.com/AlexSkr96/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerChirpsGetByID(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	dbChirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		UserID:    dbChirp.UserID,
		Body:      dbChirp.Body,
	})
}

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	var dbChirps []database.Chirp
	var err error

	sortOrder := r.URL.Query().Get("sort")
	if sortOrder == "" {
		sortOrder = "asc"
	} else if sortOrder != "desc" && sortOrder != "asc" {
		http.Error(w, fmt.Sprintf("Invalid sort order: %v", sortOrder), http.StatusBadRequest)
		return
	}

	AuthorID := r.URL.Query().Get("author_id")
	if AuthorID != "" {
		AuthorUUID, err := uuid.Parse(AuthorID)
		if err != nil {
			http.Error(w, "Couldn't parse author ID", http.StatusInternalServerError)
			return
		}
		dbChirps, err = cfg.db.GetChirpsByAuthorID(r.Context(), AuthorUUID)
	} else {
		dbChirps, err = cfg.db.GetAllChirps(r.Context())
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
	}

	slices.SortFunc(dbChirps, func(a, b database.Chirp) int {
		toReturn := 0
		if a.CreatedAt.Before(b.CreatedAt) {
			toReturn -= 1
		} else if a.CreatedAt.After(b.CreatedAt) {
			toReturn += 1
		}
		if sortOrder == "desc" {
			toReturn *= -1
		}
		return toReturn
	})

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			UserID:    dbChirp.UserID,
			Body:      dbChirp.Body,
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}
