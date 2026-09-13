package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/AFEScalante/Chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

type parameters struct {
	Body string `json:"body"`
	UserID string `json:"user_id"`
}

func validateChirps(w http.ResponseWriter, r *http.Request) parameters {
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return parameters{}
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return parameters{}
	}

	cleanedBody := cleanProfanity(params.Body)
	return parameters{
		Body: cleanedBody,
		UserID: params.UserID,
	}
}

func (cfg *apiConfig) handlerCreateChirps(w http.ResponseWriter, r *http.Request) {
	params := validateChirps(w, r)
	if params.Body == "" {
		return
	}

	newChirpParams := database.CreateChirpParams{
		ID:        uuid.New(), 
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Body:      params.Body,
		UserID:    uuid.MustParse(params.UserID),
	}
	createdChirp, err := cfg.db.CreateChirp(r.Context(), newChirpParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create chirp", err)
		return
	}
	newChirp := Chirp{
		ID:        createdChirp.ID,
		CreatedAt: createdChirp.CreatedAt,
		UpdatedAt: createdChirp.UpdatedAt,
		Body:      createdChirp.Body,
		UserID:    createdChirp.UserID,
	}
	respondWithJSON(w, http.StatusCreated, newChirp)
}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get chirps", err)
		return
	}
	var response []Chirp
	for _, c := range chirps {
		response = append(response, Chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	if chirpID == "" {
		respondWithError(w, http.StatusBadRequest, "Chirp ID is required", nil)
		return
	}
	id, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Chirp ID", err)
		return
	}
	chirp, err := cfg.db.GetChirpByID(r.Context(), id)
	if err == sql.ErrNoRows {
		respondWithError(w, http.StatusNotFound, "Chirp not found", nil)
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get chirp", err)
		return
	}
	response := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	respondWithJSON(w, http.StatusOK, response)
}