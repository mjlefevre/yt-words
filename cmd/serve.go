package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mjlefevre/sanoja/internal/handlers"
	"github.com/spf13/cobra"
)

var (
	port int
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the web server",
	Long: `Start the web server that provides transcript functionality via HTTP.
By default, the server runs on port 3000. Use --port flag to specify a different port.

The server exposes the following endpoints:
  GET /ytt - Get YouTube video transcripts
  GET /ytt?help - View API documentation

Example:
  sanoja serve
  sanoja serve --port 8080`,
	RunE: func(cmd *cobra.Command, args []string) error {
		addr := fmt.Sprintf(":%d", port)

		// Create handlers
		transcriptHandler := handlers.NewTranscriptHandler(port)

		// Setup routes
		http.HandleFunc("/ytt", transcriptHandler.GetTranscript)

		log.Printf("Starting server on http://localhost%s", addr)
		return http.ListenAndServe(addr, nil)
	},
}

func init() {
	serveCmd.Flags().IntVarP(&port, "port", "p", 3000, "Port to run the server on")
	rootCmd.AddCommand(serveCmd)
}
