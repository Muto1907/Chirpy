package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/Muto1907/Chirpy/internal/auth"
	"github.com/Muto1907/Chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) createChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Error decoding request: ", err)
		return
	}
	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get Authorization Header", err)
	}
	id, err := auth.ValidateJWT(bearer, cfg.secretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Authorization error, user not logged in", err)
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long", nil)
		return
	}
	cleanedBody := replaceBadWords(params.Body)
	createChirpParams := database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: id, //params.UserId,
	}
	chirpDb, err := cfg.queries.CreateChirp(r.Context(), createChirpParams)
	if err != nil {
		respondWithError(w, 500, "Error creating Chirp in DB", err)
		return
	}
	returnChirp := Chirp{
		Id:        chirpDb.ID,
		CreatedAt: chirpDb.CreatedAt,
		UpdatedAt: chirpDb.UpdatedAt,
		Body:      chirpDb.Body,
		UserId:    chirpDb.UserID,
	}

	respondWithJSON(w, 201, returnChirp)
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
