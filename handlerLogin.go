package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Muto1907/Chirpy/internal/auth"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int64  `json:"expires_in_seconds"`
	}
	type response struct {
		User
		Token string `json:"token"`
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
		return
	}
	var dur time.Duration
	if params.ExpiresInSeconds == 0 || params.ExpiresInSeconds > 3600 {
		dur = time.Hour
	} else {
		dur = time.Second * time.Duration(params.ExpiresInSeconds)
	}
	token, err := auth.MakeJWT(user.ID, cfg.secretKey, dur)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt generate JWT", err)
		return
	}
	responseUser := response{User{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}, token,
	}
	respondWithJSON(w, http.StatusOK, responseUser)
}
