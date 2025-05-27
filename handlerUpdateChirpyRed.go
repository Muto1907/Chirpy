package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Muto1907/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) updateChirpyRed(w http.ResponseWriter, r *http.Request) {
	type requestParameters struct {
		Event string `json:"event"`
		Data  struct {
			UserId uuid.UUID `json:"user_id"`
		} `json:"data"`
	}
	type response struct{}
	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		respondWithError(w, 401, "API Key is missing in Request Header", err)
		return
	}
	if apiKey != cfg.polkakey {
		respondWithError(w, 401, "Wrong API Key", err)
		return
	}
	decoder := json.NewDecoder(r.Body)
	params := &requestParameters{}
	err = decoder.Decode(params)
	if err != nil {
		respondWithError(w, 500, "Couldn't decode webhook", err)
		return
	}
	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}
	err = cfg.queries.SetChirpyRed(context.Background(), params.Data.UserId)
	if err != nil {
		respondWithError(w, 404, "Couldn't set is_chirpy_red", err)
		return
	}
	respondWithJSON(w, 204, response{})
}
