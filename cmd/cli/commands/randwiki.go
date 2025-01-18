package commands

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mjlefevre/sanoja/internal/browser"
	"github.com/spf13/cobra"
	"golang.org/x/net/html"
)

var textOnly bool

var randwikiCmd = &cobra.Command{
	Use:   "randwiki",
	Short: "Get a random Wikipedia article",
	Long: `Get a random Wikipedia article.

Examples:
  sanoja randwiki        # Get random article with HTML
  sanoja randwiki -t     # Get random article text only`,
	RunE: func(cmd *cobra.Command, args []string) error {
		url := "https://en.wikipedia.org/wiki/Special:Random"

		// Create request with headers
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return fmt.Errorf("error creating request: %v", err)
		}

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

		// Parse HTML
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return fmt.Errorf("error parsing page: %v", err)
		}

		// Get article title
		title := doc.Find("#firstHeading").First().Text()
		if title == "" {
			return fmt.Errorf("could not find article title")
		}

		// Get article content
		bodyContent := doc.Find("#bodyContent").First()
		if bodyContent.Length() == 0 {
			return fmt.Errorf("could not find article body content")
		}

		content := bodyContent.Find("#mw-content-text .mw-parser-output").First()
		if content.Length() == 0 {
			return fmt.Errorf("could not find article content")
		}

		// Remove unwanted elements
		//content.Find("table").Remove()
		content.Find("style").Remove()  // Remove style tags
		content.Find("script").Remove() // Remove script tags
		// Remove HTML comments
		content.Contents().Each(func(i int, s *goquery.Selection) {
			if html.CommentNode == s.Get(0).Type {
				s.Remove()
			}
		})
		content.Find(".mw-editsection").Remove()                          // Remove edit links
		content.Find(".reference").Remove()                               // Remove reference numbers
		content.Find(".reflist").Remove()                                 // Remove references section
		content.Find(".navbox").Remove()                                  // Remove navigation boxes
		content.Find(".mw-parser-output").RemoveClass("mw-parser-output") // Remove parser output class

		var output string
		if textOnly {
			// Get only paragraph text
			var texts []string
			content.Find("p").Each(func(i int, s *goquery.Selection) {
				text := strings.TrimSpace(s.Text())
				if text != "" {
					texts = append(texts, text)
				}
			})
			output = strings.Join(texts, "\n\n")
		} else {
			output, err = content.Html()
			if err != nil {
				return fmt.Errorf("error getting HTML: %v", err)
			}
		}

		fmt.Printf("Title: %s\n\n%s\n", title, output)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(randwikiCmd)
	randwikiCmd.Flags().BoolVarP(&textOnly, "text", "t", false, "Output text only (no HTML)")
}
