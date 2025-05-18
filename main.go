package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/Muto1907/Chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error Loading env variables: %s", err)
	}
	dbURL := os.Getenv("DB_URL")
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatalf("No Secret Key found.")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error opening database: %s", err)
	}
	dbQueries := database.New(db)
	platform := os.Getenv("PLATFORM")
	const filePathRoot = "."
	const port = "8080"
	cfg := &apiConfig{
		atomic.Int32{},
		dbQueries,
		platform,
		secretKey,
	}
	serveMux := http.NewServeMux()
	serveMux.Handle("/app/", http.StripPrefix("/app", cfg.middlewareMetricsInc(http.FileServer(http.Dir(filePathRoot)))))
	serveMux.HandleFunc("GET /admin/metrics", cfg.metricsHandler)
	serveMux.HandleFunc("GET /api/healthz", handlerReady)
	serveMux.HandleFunc("POST /admin/reset", cfg.resetHandler)
	serveMux.HandleFunc("POST /api/users", cfg.createUser)
	serveMux.HandleFunc("POST /api/chirps", cfg.createChirp)
	serveMux.HandleFunc("POST /api/login", cfg.loginHandler)
	serveMux.HandleFunc("GET /api/chirps", cfg.getChirps)
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", cfg.getChirp)
	serveMux.HandleFunc("POST /api/refresh", cfg.refreshToken)
	serveMux.HandleFunc("POST /api/revoke", cfg.revokeToken)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}
	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	log.Fatal(server.ListenAndServe())

}

type apiConfig struct {
	fileServerHits atomic.Int32
	queries        *database.Queries
	platform       string
	secretKey      string
}

type User struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}
