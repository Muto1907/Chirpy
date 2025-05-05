package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Muto1907/Chirpy/internal/auth"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := &parameters{}
	err := decoder.Decode(params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode Request Body", err)
		return
	}
	user, err := cfg.queries.GetUserByEmail(context.Background(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't fetch user from db", err)
		return
	}
	err = auth.CheckPassWordHash(user.HashedPassword, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
	}

	responseUser := User{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	respondWithJSON(w, http.StatusOK, response{responseUser})
}
