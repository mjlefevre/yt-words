package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/mjlefevre/sanoja/internal/browser"
	"github.com/spf13/cobra"
)

var ytdbCmd = &cobra.Command{
	Use:   "ytdb [VIDEO]",
	Short: "Get a YouTube video description from a YouTube video ID or URL",
	Long: `Get a YouTube video description box content.

VIDEO can be either:
  - A YouTube video ID (e.g., k82RwXqZHY8)
  - A full YouTube URL (e.g., https://www.youtube.com/watch?v=k82RwXqZHY8)

Examples:
  sanoja ytdb k82RwXqZHY8
  sanoja ytdb https://www.youtube.com/watch?v=k82RwXqZHY8`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := args[0]
		videoID := extractVideoID(input)
		if videoID == "" {
			return fmt.Errorf("invalid YouTube URL or Video ID: %s", input)
		}

		url := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)
		fmt.Println("Fetching URL:", url)

		// Create request with headers
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return fmt.Errorf("error creating request: %v", err)
		}

		// Set default headers
		browser.SetDefaultHeaders(req)

		// Create client with timeout
		client := &http.Client{
			Timeout: 10 * time.Second,
		}

		// Get the page
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("error fetching page: %v", err)
		}
		defer resp.Body.Close()

		// Read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading response: %v", err)
		}

		// Convert body to string
		html := string(body)

		// Try different patterns to find description
		var description string

		// Pattern 1: Look for simpleText description
		re := regexp.MustCompile(`"description":\{"simpleText":"(.*?)"\}`)
		if matches := re.FindStringSubmatch(html); len(matches) > 1 {
			description = matches[1]
		}

		// Pattern 2: Look for description in runs format
		if description == "" {
			re = regexp.MustCompile(`"description":\{"runs":\[(.*?)\]\}`)
			if matches := re.FindStringSubmatch(html); len(matches) > 1 {
				// Extract text from runs
				runsData := matches[1]
				textRe := regexp.MustCompile(`"text":"(.*?)"`)
				texts := textRe.FindAllStringSubmatch(runsData, -1)

				var parts []string
				for _, match := range texts {
					if len(match) > 1 {
						parts = append(parts, match[1])
					}
				}
				description = strings.Join(parts, "")
			}
		}

		if description == "" {
			return fmt.Errorf("could not find description for video")
		}

		// Unescape unicode characters
		description = strings.ReplaceAll(description, "\\n", "\n")
		description = strings.ReplaceAll(description, "\\\"", "\"")
		description = strings.ReplaceAll(description, "\\r", "\r")
		description = strings.ReplaceAll(description, "\\t", "\t")

		// Clean up the text
		description = strings.TrimSpace(description)
		description = strings.ReplaceAll(description, "\n\n\n", "\n")

		fmt.Printf("Description for video %s:\n%s\n", videoID, description)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(ytdbCmd)
}

func extractVideoID(input string) string {
	// Handle full URL
	if strings.Contains(input, "youtube.com/watch?v=") {
		parts := strings.Split(input, "v=")
		if len(parts) != 2 {
			return ""
		}
		// Handle additional URL parameters
		id := strings.Split(parts[1], "&")[0]
		return id
	}

	// Handle short URL
	if strings.Contains(input, "youtu.be/") {
		parts := strings.Split(input, "youtu.be/")
		if len(parts) != 2 {
			return ""
		}
		return parts[1]
	}

	// Assume it's a direct video ID
	return input
}
