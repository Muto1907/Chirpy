package main

import "net/http"

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {

	result := []Chirp{}
	chirpsDb, err := cfg.queries.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, 500, "Couldn't get Chirps", err)
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
