package commands

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sanoja",
	Short: "Process YouTube video transcripts",
	Long: `Sanoja is a tool for processing YouTube video transcripts.
    
Example usage:
  sanoja yt https://www.youtube.com/watch?v=abc123xyz
  sanoja yt abc123xyz`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(ytCmd)
}
