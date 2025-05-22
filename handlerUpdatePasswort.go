package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Muto1907/Chirpy/internal/auth"
	"github.com/Muto1907/Chirpy/internal/database"
)

func (cfg *apiConfig) updatePasswordAndEmail(w http.ResponseWriter, r *http.Request) {
	type requestParameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
	}

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "access token is missing", err)
		return
	}
	id, err := auth.ValidateJWT(accessToken, cfg.secretKey)
	if err != nil {
		respondWithError(w, 401, "Couldn't validate access token", err)
		return
	}
	var params *requestParameters
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "couldn't decode Request", err)
		return
	}
	pw, err := auth.HashPassWord(params.Password)
	if err != nil {
		respondWithError(w, 500, "couldn't hash Password", err)
		return
	}
	setUserPWParams := database.SetUserPasswordParams{
		HashedPassword: pw,
		Email:          params.Email,
		ID:             id,
	}
	user, err := cfg.queries.SetUserPassword(context.Background(), setUserPWParams)
	if err != nil {
		respondWithError(w, 500, "coultn't set new password in DB", err)
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			Id:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}
