package main

import (
	"log"
	"net/http"
)

func main() {
	const filePathRoot = "."
	const port = "8080"
	serveMux := http.NewServeMux()
	serveMux.Handle("/", http.FileServer(http.Dir(filePathRoot)))
	server := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}
	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	log.Fatal(server.ListenAndServe())
}
