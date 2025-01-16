package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mjlefevre/sanoja/web/handlers"
)

func main() {
	port := 8080
	transcriptHandler := handlers.NewTranscriptHandler(port)

	// Register routes
	http.HandleFunc("/ytt", transcriptHandler.GetTranscript)

	// Add a simple health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting server on port %d", port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
