package main

import (
	"context"
	"net/http"
	"time"

	"github.com/Muto1907/Chirpy/internal/auth"
)

func (cfg *apiConfig) refreshToken(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	BearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Refresh Token is missing in Header", err)
		return
	}
	token, err := cfg.queries.GetRefreshToken(context.Background(), BearerToken)
	if err != nil {
		respondWithError(w, 401, "Refresh Token does not exist", err)
		return
	}
	if token.RevokedAt.Valid {
		respondWithError(w, 401, "Refresh Token has been revoked", err)
	}
	if token.ExpiresAt.Before(time.Now().UTC()) {
		respondWithError(w, 401, "Refresh Token is expired", err)
		return
	}
	user, err := cfg.queries.GetUserByID(context.Background(), token.UserID)
	if err != nil {
		respondWithError(w, 500, "Unable to retrieve user from Refresh Token", err)
		return
	}
	jot, err := auth.MakeJWT(user.ID, cfg.secretKey, time.Hour)
	if err != nil {
		respondWithError(w, 500, "Unable to generate acces Token", err)
		return
	}
	respondWithJSON(w, http.StatusOK, response{jot})
}
