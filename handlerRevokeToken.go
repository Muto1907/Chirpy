package main

import (
	"context"
	"net/http"

	"github.com/Muto1907/Chirpy/internal/auth"
)

func (cfg *apiConfig) revokeToken(w http.ResponseWriter, r *http.Request) {
	HeaderToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Refresh Token is missing in Header", err)
		return
	}
	rft, err := cfg.queries.GetRefreshToken(context.Background(), HeaderToken)
	if err != nil {
		respondWithError(w, 500, "Coudln't retrieve Refreh Token from DB", err)
		return
	}
	err = cfg.queries.RevokeRefreshToken(context.Background(), rft.Token)
	if err != nil {
		respondWithError(w, 500, "Couldn't update Refresh Token to be revoked in DB", err)
		return
	}
	w.WriteHeader(204)
}
