package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Muto1907/Chirpy/internal/auth"
	"github.com/Muto1907/Chirpy/internal/database"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
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
	dur := time.Hour

	token, err := auth.MakeJWT(user.ID, cfg.secretKey, dur)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt generate JWT", err)
		return
	}
	rf_token, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate Refresh Token string", err)
	}
	rft_params := database.CreateRefreshTokenParams{
		Token:  rf_token,
		UserID: user.ID,
	}
	rft, err := cfg.queries.CreateRefreshToken(context.Background(), rft_params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't insert Refresh Token into database", err)
		return
	}
	responseUser := response{User{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	},
		token,
		rft.Token,
	}
	respondWithJSON(w, http.StatusOK, responseUser)
}
