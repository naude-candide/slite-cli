package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/naude-candide/slite-cli/internal/output"
)

var (
	docsOwner  string
	docsLimit  int
	docsOffset int
	docsCursor string

	docsCreateTitle    string
	docsCreateMarkdown string
	docsCreateParent   string
	docsCreateBodyJSON string

	docsUpdateTitle    string
	docsUpdateMarkdown string
	docsUpdateParent   string
	docsUpdateBodyJSON string
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

		notes, err := client.ListNotes(context.Background(), docsOwner, docsLimit, docsOffset, docsCursor)
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

var docsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a doc",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		payload, err := buildNotePayload(docsCreateTitle, docsCreateMarkdown, docsCreateParent, docsCreateBodyJSON)
		if err != nil {
			return err
		}

		note, err := client.CreateNote(context.Background(), payload)
		if err != nil {
			return err
		}

		return output.RenderNote(note, jsonOutput)
	},
}

var docsUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a doc",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		payload, err := buildNotePayload(docsUpdateTitle, docsUpdateMarkdown, docsUpdateParent, docsUpdateBodyJSON)
		if err != nil {
			return err
		}

		note, err := client.UpdateNote(context.Background(), args[0], payload)
		if err != nil {
			return err
		}

		return output.RenderNote(note, jsonOutput)
	},
}

var docsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a doc",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		result, err := client.DeleteNote(context.Background(), args[0])
		if err != nil {
			return err
		}
		if result.ID == "" {
			result.ID = args[0]
		}

		return output.RenderDelete(result, jsonOutput)
	},
}

func init() {
	docsListCmd.Flags().StringVar(&docsOwner, "owner", "", "Owner user ID")
	docsListCmd.Flags().IntVar(&docsLimit, "limit", 20, "Page size")
	docsListCmd.Flags().IntVar(&docsOffset, "offset", 0, "Offset")
	docsListCmd.Flags().StringVar(&docsCursor, "cursor", "", "Pagination cursor")

	docsCreateCmd.Flags().StringVar(&docsCreateTitle, "title", "", "Doc title")
	docsCreateCmd.Flags().StringVar(&docsCreateMarkdown, "markdown", "", "Doc markdown content")
	docsCreateCmd.Flags().StringVar(&docsCreateParent, "parent", "", "Parent doc ID")
	docsCreateCmd.Flags().StringVar(&docsCreateBodyJSON, "body-json", "", "Raw JSON request body")

	docsUpdateCmd.Flags().StringVar(&docsUpdateTitle, "title", "", "Doc title")
	docsUpdateCmd.Flags().StringVar(&docsUpdateMarkdown, "markdown", "", "Doc markdown content")
	docsUpdateCmd.Flags().StringVar(&docsUpdateParent, "parent", "", "Parent doc ID")
	docsUpdateCmd.Flags().StringVar(&docsUpdateBodyJSON, "body-json", "", "Raw JSON request body")

	docsCmd.AddCommand(docsListCmd)
	docsCmd.AddCommand(docsGetCmd)
	docsCmd.AddCommand(docsCreateCmd)
	docsCmd.AddCommand(docsUpdateCmd)
	docsCmd.AddCommand(docsDeleteCmd)
	rootCmd.AddCommand(docsCmd)
}

func buildNotePayload(title, markdown, parent, bodyJSON string) (map[string]any, error) {
	payload := map[string]any{}
	if strings.TrimSpace(bodyJSON) != "" {
		if err := json.Unmarshal([]byte(bodyJSON), &payload); err != nil {
			return nil, fmt.Errorf("invalid --body-json: %w", err)
		}
	}

	if title != "" {
		payload["title"] = title
	}
	if markdown != "" {
		payload["markdown"] = markdown
	}
	if parent != "" {
		payload["parentNoteId"] = parent
	}

	if len(payload) == 0 {
		return nil, fmt.Errorf("provide at least one field via flags or --body-json")
	}

	return payload, nil
}
