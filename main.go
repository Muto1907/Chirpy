package main

import (
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
	serveMux.HandleFunc("GET /admin/metrics", cfg.metricsHandler)
	serveMux.HandleFunc("GET /api/healthz", handlerReady)
	serveMux.HandleFunc("POST /admin/reset", cfg.resetHandler)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}
	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	log.Fatal(server.ListenAndServe())

}

type apiConfig struct {
	fileServerHits atomic.Int32
}
