package main

import (
	"encoding/json"
	"net/http"

	"github.com/Muto1907/Chirpy/internal/auth"
	"github.com/Muto1907/Chirpy/internal/database"
)

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode Request", err)
		return
	}
	password, err := auth.HashPassWord(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
	}
	userDb, err := cfg.queries.CreateUser(r.Context(), database.CreateUserParams{Email: params.Email, HashedPassword: password})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating User in databse", err)
		return
	}
	user := User{
		Id:        userDb.ID,
		CreatedAt: userDb.CreatedAt,
		UpdatedAt: userDb.UpdatedAt,
		Email:     userDb.Email,
	}
	respondWithJSON(w, 201, response{user})
}
