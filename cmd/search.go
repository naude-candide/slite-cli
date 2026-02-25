package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/naude-candide/slite-cli/internal/output"
)

var (
	searchLimit  int
	searchOffset int
	searchCursor string
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search docs",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		result, err := client.SearchNotes(context.Background(), args[0], searchLimit, searchOffset, searchCursor)
		if err != nil {
			return err
		}

		return output.RenderSearch(result, jsonOutput)
	},
}

func init() {
	searchCmd.Flags().IntVar(&searchLimit, "limit", 20, "Page size")
	searchCmd.Flags().IntVar(&searchOffset, "offset", 0, "Offset")
	searchCmd.Flags().StringVar(&searchCursor, "cursor", "", "Pagination cursor")
	rootCmd.AddCommand(searchCmd)
}
