package cmd

import (
	"fmt"

	"github.com/naude/slite-cli/internal/slite"
)

func newClient() (*slite.Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("missing API key: set SLITE_API_KEY")
	}

	return slite.NewClient(slite.Config{
		APIKey:  apiKey,
		BaseURL: baseURL,
		Timeout: timeout,
		Debug:   debug,
	})
}
