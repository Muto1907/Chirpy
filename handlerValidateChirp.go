package main

import (
	"encoding/json"
	"net/http"
)

func validateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnValues struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	param := parameters{}
	err := decoder.Decode(&param)
	if err != nil {
		respondWithError(w, 500, "Couln't decode parameters", err)
	}
	if len(param.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long", nil)
		return
	}
	respBody := returnValues{
		Valid: true,
	}
	respondWithJSON(w, 200, respBody)
}
