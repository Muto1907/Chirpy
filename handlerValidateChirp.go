package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

func validateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnValues struct {
		CleanedBody string `json:"cleaned_body"`
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
		CleanedBody: replaceBadWords(param.Body),
	}
	respondWithJSON(w, 200, respBody)
}

func replaceBadWords(in string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(in, " ")
	for i, word := range words {
		if slices.Contains(badWords, strings.ToLower(word)) {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
