package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/naude-candide/slite-cli/internal/output"
)

var meCmd = &cobra.Command{
	Use:   "me",
	Short: "Show current authenticated user",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		me, err := client.Me(context.Background())
		if err != nil {
			return err
		}

		return output.RenderMe(me, jsonOutput)
	},
}

func init() {
	rootCmd.AddCommand(meCmd)
}
