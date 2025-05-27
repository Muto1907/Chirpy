package main

import (
	"context"
	"net/http"
	"sort"

	"github.com/Muto1907/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {

	result := []Chirp{}
	query_param_author := r.URL.Query().Get("author_id")
	query_param_sort_order := r.URL.Query().Get("sort")
	var chirpsDb []database.Chirp
	var err error
	if query_param_author == "" {
		chirpsDb, err = cfg.queries.GetChirps(r.Context())
		if err != nil {
			respondWithError(w, 500, "Couldn't get Chirps", err)
		}
	} else {
		author_id, err := uuid.Parse(query_param_author)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Couldn't parse author_id", err)
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
	if query_param_sort_order == "desc" {
		sort.Slice(result,
			func(i int, j int) bool {
				return result[i].CreatedAt.After(result[j].CreatedAt)
			})
	}
	respondWithJSON(w, 200, result)
}
