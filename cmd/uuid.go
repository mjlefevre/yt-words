package main

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var (
	count int
)

// uuidCmd represents the uuid command
var uuidCmd = &cobra.Command{
	Use:   "uuid",
	Short: "Generate one or more UUIDs",
	Long: `Generate one or more random UUIDs (v4).

Examples:
  sanoja uuid           # Generate one UUID
  sanoja uuid -n 5      # Generate 5 UUIDs`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if count < 1 {
			return fmt.Errorf("number of UUIDs must be positive, got %d", count)
		}

		var uuids []string
		for i := 0; i < count; i++ {
			id := uuid.New()
			uuids = append(uuids, id.String())
		}
		fmt.Println(strings.Join(uuids, "\n"))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(uuidCmd)
	uuidCmd.Flags().IntVarP(&count, "number", "n", 1, "Number of UUIDs to generate")
}
