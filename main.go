package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	const filePathRoot = "."
	const port = "8080"
	cfg := &apiConfig{
		atomic.Int32{},
	}
	serveMux := http.NewServeMux()
	serveMux.Handle("/app/", http.StripPrefix("/app", cfg.middlewareMetricsInc(http.FileServer(http.Dir(filePathRoot)))))
	serveMux.HandleFunc("GET /metrics", cfg.metricsHandler)
	serveMux.HandleFunc("GET /healthz", handlerReady)
	serveMux.HandleFunc("POST /reset", cfg.resetHandler)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}
	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	log.Fatal(server.ListenAndServe())

}

func handlerReady(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte{'O', 'K'})
}

type apiConfig struct {
	fileServerHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileServerHits.Load())))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, req *http.Request) {
	cfg.fileServerHits.Store(0)
	w.WriteHeader(http.StatusOK)
}
