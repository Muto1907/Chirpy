package main

import (
	"context"
	"encoding/json"
	"net/http"

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
	decoder := json.NewDecoder(r.Body)
	params := &requestParameters{}
	err := decoder.Decode(params)
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
