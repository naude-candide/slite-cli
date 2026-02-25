package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/naude-candide/slite-cli/internal/slite"
)

const apiKeyExportPrefix = "export SLITE_API_KEY="

var (
	authNoPersist bool
	authFromStdin bool
	authCheck     bool
	authShellFile string
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage CLI authentication",
}

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Set and optionally persist your Slite API key",
	RunE: func(cmd *cobra.Command, args []string) error {
		key, err := readAPIKey(authFromStdin)
		if err != nil {
			return err
		}
		if key == "" {
			return fmt.Errorf("empty API key")
		}

		if err := verifyAPIKey(key); err != nil {
			return err
		}

		os.Setenv("SLITE_API_KEY", key)
		apiKey = key

		if authNoPersist {
			fmt.Println("API key is valid and set for current process (not persisted).")
			return nil
		}

		if err := upsertAPIKeyInProfile(authShellFile, key); err != nil {
			return err
		}
		fmt.Printf("API key is valid and saved to %s\n", authShellFile)
		fmt.Printf("Reload shell: source %s\n", authShellFile)
		return nil
	},
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show API key status",
	RunE: func(cmd *cobra.Command, args []string) error {
		envKey := strings.TrimSpace(os.Getenv("SLITE_API_KEY"))
		profileKey, err := readAPIKeyFromProfile(authShellFile)
		if err != nil {
			return err
		}

		if envKey != "" {
			fmt.Printf("env: set (%s)\n", maskKey(envKey))
		} else {
			fmt.Println("env: not set")
		}

		if profileKey != "" {
			fmt.Printf("profile (%s): set (%s)\n", authShellFile, maskKey(profileKey))
		} else {
			fmt.Printf("profile (%s): not set\n", authShellFile)
		}

		if !authCheck {
			return nil
		}

		key := envKey
		if key == "" {
			key = profileKey
		}
		if key == "" {
			return fmt.Errorf("no API key available to verify")
		}

		if err := verifyAPIKey(key); err != nil {
			return err
		}
		fmt.Println("api check: ok")
		return nil
	},
}

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove persisted API key and unset current process value",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := removeAPIKeyFromProfile(authShellFile); err != nil {
			return err
		}
		os.Unsetenv("SLITE_API_KEY")
		apiKey = ""
		fmt.Printf("Removed SLITE_API_KEY from %s\n", authShellFile)
		fmt.Println("Unset SLITE_API_KEY for current process")
		return nil
	},
}

func init() {
	authShellFile = defaultShellProfilePath()

	authCmd.PersistentFlags().StringVar(&authShellFile, "shell-file", authShellFile, "Shell profile file to read/write")

	authLoginCmd.Flags().BoolVar(&authNoPersist, "no-persist", false, "Do not save key to shell profile")
	authLoginCmd.Flags().BoolVar(&authFromStdin, "from-stdin", false, "Read API key from stdin")
	authStatusCmd.Flags().BoolVar(&authCheck, "check", false, "Validate key against Slite API")

	authCmd.AddCommand(authLoginCmd)
	authCmd.AddCommand(authStatusCmd)
	authCmd.AddCommand(authLogoutCmd)
	rootCmd.AddCommand(authCmd)
}

func defaultShellProfilePath() string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return ".zshrc"
	}
	return filepath.Join(home, ".zshrc")
}

func readAPIKey(fromStdin bool) (string, error) {
	if fromStdin {
		b, err := io.ReadAll(bufio.NewReader(os.Stdin))
		if err != nil {
			return "", fmt.Errorf("read stdin: %w", err)
		}
		return strings.TrimSpace(string(b)), nil
	}

	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return "", fmt.Errorf("open tty: %w", err)
	}
	defer tty.Close()

	fmt.Fprint(tty, "Enter Slite API key: ")
	b, err := term.ReadPassword(int(tty.Fd()))
	fmt.Fprintln(tty)
	if err != nil {
		return "", fmt.Errorf("read API key: %w", err)
	}
	return strings.TrimSpace(string(b)), nil
}

func verifyAPIKey(key string) error {
	client, err := slite.NewClient(slite.Config{
		APIKey:  key,
		BaseURL: baseURL,
		Timeout: timeout,
		Debug:   debug,
	})
	if err != nil {
		return err
	}

	if _, err := client.Me(context.Background()); err != nil {
		return fmt.Errorf("API key verification failed: %w", err)
	}
	return nil
}

func upsertAPIKeyInProfile(path, key string) error {
	line := apiKeyExportPrefix + shellSingleQuote(key)
	lines, err := readProfileLines(path)
	if err != nil {
		return err
	}

	out := make([]string, 0, len(lines)+1)
	replaced := false
	for _, l := range lines {
		if strings.HasPrefix(strings.TrimSpace(l), apiKeyExportPrefix) {
			if !replaced {
				out = append(out, line)
				replaced = true
			}
			continue
		}
		out = append(out, l)
	}
	if !replaced {
		out = append(out, line)
	}

	content := strings.Join(out, "\n")
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func readAPIKeyFromProfile(path string) (string, error) {
	lines, err := readProfileLines(path)
	if err != nil {
		return "", err
	}
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if strings.HasPrefix(trimmed, apiKeyExportPrefix) {
			value := strings.TrimPrefix(trimmed, apiKeyExportPrefix)
			return strings.Trim(strings.TrimSpace(value), "\"'"), nil
		}
	}
	return "", nil
}

func removeAPIKeyFromProfile(path string) error {
	lines, err := readProfileLines(path)
	if err != nil {
		return err
	}
	out := make([]string, 0, len(lines))
	for _, l := range lines {
		if strings.HasPrefix(strings.TrimSpace(l), apiKeyExportPrefix) {
			continue
		}
		out = append(out, l)
	}
	content := strings.Join(out, "\n")
	if content != "" && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func readProfileLines(path string) ([]string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	content := strings.ReplaceAll(string(b), "\r\n", "\n")
	content = strings.TrimSuffix(content, "\n")
	if content == "" {
		return []string{}, nil
	}
	return strings.Split(content, "\n"), nil
}

func shellSingleQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}

func maskKey(key string) string {
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "..." + key[len(key)-4:]
}
