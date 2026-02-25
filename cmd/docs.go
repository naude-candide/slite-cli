package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/naude-candide/slite-cli/internal/output"
)

var (
	docsOwner  string
	docsLimit  int
	docsOffset int
)

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Manage Slite docs",
}

var docsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List docs",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		notes, err := client.ListNotes(context.Background(), docsOwner, docsLimit, docsOffset)
		if err != nil {
			return err
		}

		return output.RenderNotes(notes, jsonOutput)
	},
}

var docsGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a single doc by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		note, err := client.GetNote(context.Background(), args[0])
		if err != nil {
			return err
		}

		return output.RenderNote(note, jsonOutput)
	},
}

func init() {
	docsListCmd.Flags().StringVar(&docsOwner, "owner", "", "Owner user ID")
	docsListCmd.Flags().IntVar(&docsLimit, "limit", 20, "Page size")
	docsListCmd.Flags().IntVar(&docsOffset, "offset", 0, "Offset")

	docsCmd.AddCommand(docsListCmd)
	docsCmd.AddCommand(docsGetCmd)
	rootCmd.AddCommand(docsCmd)
}
