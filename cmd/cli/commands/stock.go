package commands

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mjlefevre/sanoja/internal/browser"
	"github.com/spf13/cobra"
)

var stockCmd = &cobra.Command{
	Use:   "stock [SYMBOL]",
	Short: "Get stock information from Yahoo Finance",
	Long: `Get stock information from Yahoo Finance.

Example:
  sanoja stock AAPL    # Get Apple stock information`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		symbol := strings.ToUpper(args[0])
		url := fmt.Sprintf("https://finance.yahoo.com/quote/%s", symbol)

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

		// Get title from main h1
		title := doc.Find("main h1").First().Text()
		if title == "" {
			return fmt.Errorf("could not find stock title")
		}

		// Get price section
		priceSection := doc.Find("section[data-testid='quote-price']").First()
		if priceSection.Length() == 0 {
			return fmt.Errorf("could not find price section")
		}

		// Get price text and clean it up
		priceText := priceSection.Text()
		priceText = strings.TrimSpace(priceText)
		priceText = strings.ReplaceAll(priceText, "  ", "\n")

		// Output results
		fmt.Printf("%s\n%s\n", title, priceText)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(stockCmd)
}
