package main

import (
	"encoding/json"
	"net/http"
)

type parameters struct {
	Body string `json:"body"`
}

type responseVal struct {
	CleanedBody string `json:"cleaned_body"`
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	cleanedBody := cleanProfanity(params.Body)
	respondWithJSON(w, http.StatusOK, responseVal{CleanedBody: cleanedBody})
}