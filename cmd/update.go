package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	updateVersion string
	updateRepo    string
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update slite to the latest release",
	RunE: func(cmd *cobra.Command, args []string) error {
		exePath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("resolve executable path: %w", err)
		}
		installDir := filepath.Dir(exePath)

		installerURL := "https://raw.githubusercontent.com/naude-candide/slite-cli/main/scripts/install.sh"
		resp, err := http.Get(installerURL)
		if err != nil {
			return fmt.Errorf("download installer: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return fmt.Errorf("download installer: status=%d", resp.StatusCode)
		}

		script, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read installer: %w", err)
		}

		bash := exec.Command("bash")
		bash.Stdin = bytes.NewReader(script)
		bash.Stdout = os.Stdout
		bash.Stderr = os.Stderr
		bash.Env = append(os.Environ(),
			"INSTALL_DIR="+installDir,
			"VERSION="+updateVersion,
			"REPO="+updateRepo,
			"SKIP_API_KEY_PROMPT=1",
		)

		if err := bash.Run(); err != nil {
			return fmt.Errorf("run installer: %w", err)
		}

		fmt.Printf("Updated slite in %s\n", installDir)
		return nil
	},
}

func init() {
	updateCmd.Flags().StringVar(&updateVersion, "version", "latest", "Release tag to install (e.g. v0.1.4)")
	updateCmd.Flags().StringVar(&updateRepo, "repo", "naude-candide/slite-cli", "GitHub repository to install from (owner/repo)")
	rootCmd.AddCommand(updateCmd)
}
