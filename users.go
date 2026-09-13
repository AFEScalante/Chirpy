package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AFEScalante/Chirpy/internal/auth"
	"github.com/AFEScalante/Chirpy/internal/database"
	"github.com/google/uuid"
)

type userParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userLoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	ExpiresInSeconds *int `json:"expires_in_seconds"`
}

type userLoginResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

const (
	DEFAULT_EXPIRES_IN_SECONDS int = 3600
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := userParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	newUserParams := database.CreateUserParams{
		ID:           uuid.New(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Email:        params.Email,
		PasswordHash: hashedPassword,
	}

	createdUser, err := cfg.db.CreateUser(r.Context(), newUserParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}
	returnUser := User{
		ID: createdUser.ID,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
		Email: createdUser.Email,
	}
	respondWithJSON(w, http.StatusCreated, returnUser)
}

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := userLoginParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	if params.ExpiresInSeconds == nil {
		params.ExpiresInSeconds = new(int)
		*params.ExpiresInSeconds = DEFAULT_EXPIRES_IN_SECONDS
	}

	if *params.ExpiresInSeconds > DEFAULT_EXPIRES_IN_SECONDS {
		*params.ExpiresInSeconds = DEFAULT_EXPIRES_IN_SECONDS
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password", err)
		return
	}

	checkedPassword, err := auth.CheckPasswordHash(params.Password, user.PasswordHash)
	if err != nil || !checkedPassword {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Second * time.Duration(*params.ExpiresInSeconds))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	returnUser := userLoginResponse{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: token,
	}

	respondWithJSON(w, http.StatusOK, returnUser)
}
