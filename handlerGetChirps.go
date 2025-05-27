package main

import (
	"context"
	"net/http"

	"github.com/Muto1907/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {

	result := []Chirp{}
	query_param := r.URL.Query().Get("author_id")
	var chirpsDb []database.Chirp
	var err error
	if query_param == "" {
		chirpsDb, err = cfg.queries.GetChirps(r.Context())
		if err != nil {
			respondWithError(w, 500, "Couldn't get Chirps", err)
		}
	} else {
		author_id, err := uuid.Parse(query_param)
		if err != nil {
			respondWithError(w, 500, "Couldn't parse author_id", err)
		}
		chirpsDb, err = cfg.queries.GetChirpsByAuthor(context.Background(), author_id)
		if err != nil {
			respondWithError(w, 500, "Couldn't get Chirps", err)
		}
	}

	for _, chirp := range chirpsDb {
		result = append(result, Chirp{
			Id:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID,
		})
	}
	respondWithJSON(w, 200, result)
}
