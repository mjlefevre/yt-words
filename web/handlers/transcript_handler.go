package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mjlefevre/sanoja/pkg/transcript"
)

// TranscriptHandler handles transcript-related HTTP requests
type TranscriptHandler struct {
	client *transcript.Client
}

// NewTranscriptHandler creates a new TranscriptHandler
func NewTranscriptHandler() *TranscriptHandler {
	return &TranscriptHandler{
		client: transcript.NewClient(),
	}
}

// GetTranscript handles requests to fetch transcripts
func (h *TranscriptHandler) GetTranscript(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check for help query parameter
	if _, wantHelp := r.URL.Query()["help"]; wantHelp {
		w.Header().Set("Content-Type", "text/plain")
		helpText := `Sanoja Transcript API Usage:

Fetch a transcript:
  GET /transcript?v=VIDEO_ID
  GET /transcript?url=YOUTUBE_URL
  GET /transcript?videoId=VIDEO_ID

Parameters:
  v         - YouTube video ID
  url       - Full YouTube video URL
  videoId   - Alternative to 'v' parameter
  json      - Add this flag to get JSON response
  help      - Show this help message

Examples:
  /transcript?v=k82RwXqZHY8
  /transcript?url=https://www.youtube.com/watch?v=k82RwXqZHY8
  /transcript?v=k82RwXqZHY8&json

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
