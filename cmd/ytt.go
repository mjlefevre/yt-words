package main

import (
	"fmt"

	"github.com/mjlefevre/sanoja/pkg/transcript"
	"github.com/spf13/cobra"
)

var yttCmd = &cobra.Command{
	Use:   "ytt [VIDEO]",
	Short: "Get a YouTube video transcript from a YouTube video ID or URL",
	Long: `Get a YouTube video transcript and generate output.

VIDEO can be either:
  - A YouTube video ID (e.g., k82RwXqZHY8)
  - A full YouTube URL (e.g., https://www.youtube.com/watch?v=k82RwXqZHY8)

Examples:
  sanoja ytt k82RwXqZHY8
  sanoja ytt https://www.youtube.com/watch?v=k82RwXqZHY8`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := args[0]
		videoID := transcript.ExtractVideoID(input)
		if videoID == "" {
			return fmt.Errorf("invalid YouTube URL or Video ID: %s", input)
		}

		client := transcript.NewClient()
		transcriptText, err := client.GetTranscriptString(videoID)
		if err != nil {
			return fmt.Errorf("error fetching transcript: %v", err)
		}

		fmt.Printf("Transcript for video %s:\n%s\n", videoID, transcriptText)
		return nil
	},
}
