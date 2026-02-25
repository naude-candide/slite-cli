package cmd

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/naude-candide/slite-cli/internal/config"
	"github.com/naude-candide/slite-cli/internal/slite"
)

var (
	jsonOutput bool
	debug      bool
	baseURL    string
	timeout    time.Duration
	apiKey     string
)

var rootCmd = &cobra.Command{
	Use:   "slite",
	Short: "CLI wrapper for the Slite API",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output raw JSON")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug logging")
	rootCmd.PersistentFlags().StringVar(&baseURL, "base-url", slite.DefaultBaseURL, "Slite API base URL")
	rootCmd.PersistentFlags().DurationVar(&timeout, "timeout", 15*time.Second, "HTTP timeout (e.g. 10s)")

	cobra.OnInitialize(func() {
		apiKey = config.APIKey()
	})
}
