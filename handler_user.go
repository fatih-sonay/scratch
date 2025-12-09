package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fatih-sonay/rssagg/internal/auth"
	"github.com/fatih-sonay/rssagg/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	type paramteres struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)

	params := paramteres{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing request body %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      params.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error creating user: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handlerGetUserByAPIKey(w http.ResponseWriter, r *http.Request) {

	// use GetAPIKey function from auth package to extract API key from request
	apiKey, err := auth.GetAPIKey(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Auth error: %v", err))
		return
	}

	user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Couldn't get user: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}
