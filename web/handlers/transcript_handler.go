package handlers

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/mjlefevre/sanoja/pkg/transcript"
	"github.com/mjlefevre/sanoja/web/templates"
)

// TranscriptHandler handles transcript-related HTTP requests
type TranscriptHandler struct {
	client *transcript.Client
	port   int
}

// NewTranscriptHandler creates a new TranscriptHandler
func NewTranscriptHandler(port int) *TranscriptHandler {
	return &TranscriptHandler{
		client: transcript.NewClient(),
		port:   port,
	}
}

// GetTranscript handles requests to fetch transcripts
func (h *TranscriptHandler) GetTranscript(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check for bookmarklet request
	if _, wantBookmarklet := r.URL.Query()["bookmarklet"]; wantBookmarklet {
		tmpl, err := template.ParseFS(templates.Files, "bookmarklet.html")
		if err != nil {
			http.Error(w, "Error loading bookmarklet template", http.StatusInternalServerError)
			return
		}

		// Replace PORT in template with actual port
		var buf bytes.Buffer
		err = tmpl.Execute(&buf, struct{ Port int }{Port: h.port})
		if err != nil {
			http.Error(w, "Error processing template", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.Write(buf.Bytes())
		return
	}

	// Check for help query parameter
	if _, wantHelp := r.URL.Query()["help"]; wantHelp {
		w.Header().Set("Content-Type", "text/plain")
		helpText := `Sanoja Transcript API Usage:

Fetch a transcript:
  GET /ytt?v=VIDEO_ID
  GET /ytt?url=YOUTUBE_URL
  GET /ytt?videoId=VIDEO_ID

Parameters:
  v           - YouTube video ID
  url         - Full YouTube video URL
  videoId     - Alternative to 'v' parameter
  json        - Add this flag to get JSON response
  bookmarklet - Get a bookmarklet for easy transcript fetching from browser
  help        - Show this help message
  
Examples:
  /ytt?v=k82RwXqZHY8
  /ytt?url=https://www.youtube.com/watch?v=k82RwXqZHY8
  /ytt?v=k82RwXqZHY8&json

Response Formats:
  - Default: Plain text
  - JSON: Add 'json' parameter or set Accept: application/json header
    Returns: {"text": "transcript content"}`
		w.Write([]byte(helpText))
		return
	}

	videoID := r.URL.Query().Get("v")
	if videoID == "" {
		videoID = r.URL.Query().Get("videoId")
	}

	// If no direct video ID, try to extract from URL
	if videoID == "" {
		url := r.URL.Query().Get("url")
		if url != "" {
			videoID = transcript.ExtractVideoID(url)
		}
	}

	if videoID == "" {
		http.Error(w, "Missing video ID or URL", http.StatusBadRequest)
		return
	}

	transcriptText, err := h.client.GetTranscriptString(videoID)
	if err != nil {
		http.Error(w, "Error fetching transcript: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if JSON format is requested
	_, wantJSON := r.URL.Query()["json"]
	acceptHeader := r.Header.Get("Accept")
	if wantJSON || acceptHeader == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"text": transcriptText,
		})
		return
	}

	// Default to plain text
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(transcriptText))
}
