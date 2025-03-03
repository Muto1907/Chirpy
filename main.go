package main

import (
	"log"
	"net/http"
)

func main() {
	const filePathRoot = "."
	const readinessPath = "/healthz"
	const port = "8080"
	serveMux := http.NewServeMux()
	serveMux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot))))
	serveMux.HandleFunc(readinessPath, handler)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}
	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	log.Fatal(server.ListenAndServe())

}

func handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte{'O', 'K'})
}
