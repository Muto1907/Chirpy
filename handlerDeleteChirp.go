package main

import (
	"context"
	"net/http"

	"github.com/Muto1907/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Couldn't get Access Token", err)
		return
	}
	id, err := auth.ValidateJWT(accessToken, cfg.secretKey)
	if err != nil {
		respondWithError(w, 401, "Couldn't validate acces token", err)
		return
	}
	chirpId, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, 500, "Couldn't parse chirpId", err)
		return
	}
	chirp, err := cfg.queries.GetChirp(context.Background(), chirpId)
	if err != nil {
		respondWithError(w, 404, "Chirp not found", err)
		return
	}
	if chirp.UserID != id {
		respondWithError(w, 403, "Unauthorized Delete Attempt", err)
		return
	}
	err = cfg.queries.DeleteChirp(context.Background(), chirpId)
	if err != nil {
		respondWithError(w, 500, "Couldn't delete Chirp", err)
		return
	}
	w.WriteHeader(204)
}
