package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mjlefevre/sanoja/web/handlers"
)

func main() {
	transcriptHandler := handlers.NewTranscriptHandler()

	// Register routes
	http.HandleFunc("/api/transcript", transcriptHandler.GetTranscript)

	// Add a simple health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	port := ":8080"
	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
