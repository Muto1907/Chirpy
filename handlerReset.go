package main

import (
	"net/http"
)

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset iis only allowed for devs"))
		return
	}
	cfg.fileServerHits.Store(0)
	err := cfg.queries.DeleteUsers(req.Context())
	if err != nil {
		respondWithError(w, 500, "error deleting users", err)
		return
	}
	w.WriteHeader(200)
	w.Write([]byte("Deleted All users and reset hit count"))
}
