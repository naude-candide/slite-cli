package output

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/naude-candide/slite-cli/internal/slite"
)

func RenderMe(me *slite.MeResponse, asJSON bool) error {
	if asJSON {
		return writeJSON(me)
	}

	tw := tabwriter.NewWriter(os.Stdout, 2, 4, 2, ' ', 0)
	fmt.Fprintln(tw, "ID\tNAME\tEMAIL\tORG")
	fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", me.ID, me.Name, me.Email, me.Organization.Name)
	return tw.Flush()
}

func RenderNotes(notes *slite.NotesResponse, asJSON bool) error {
	if asJSON {
		return writeJSON(notes)
	}

	tw := tabwriter.NewWriter(os.Stdout, 2, 4, 2, ' ', 0)
	fmt.Fprintln(tw, "ID\tTITLE\tOWNER\tUPDATED")
	for _, n := range notes.Hits {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", n.ID, n.Title, n.OwnerID, n.UpdatedAt)
	}
	return tw.Flush()
}

func RenderNote(note *slite.NoteDetail, asJSON bool) error {
	if asJSON {
		return writeJSON(note)
	}

	tw := tabwriter.NewWriter(os.Stdout, 2, 4, 2, ' ', 0)
	fmt.Fprintln(tw, "ID\tTITLE\tOWNER\tUPDATED\tURL")
	fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", note.ID, note.Title, note.OwnerID, note.UpdatedAt, note.URL)
	tw.Flush()

	if note.Markdown != "" {
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, note.Markdown)
	}

	return nil
}

func RenderSearch(result *slite.SearchResponse, asJSON bool) error {
	if asJSON {
		return writeJSON(result)
	}

	tw := tabwriter.NewWriter(os.Stdout, 2, 4, 2, ' ', 0)
	fmt.Fprintln(tw, "ID\tTITLE\tOWNER\tUPDATED")
	for _, n := range result.Hits {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", n.ID, n.Title, n.OwnerID, n.UpdatedAt)
	}
	return tw.Flush()
}

func RenderDelete(result *slite.DeleteResponse, asJSON bool) error {
	if asJSON {
		return writeJSON(result)
	}

	tw := tabwriter.NewWriter(os.Stdout, 2, 4, 2, ' ', 0)
	fmt.Fprintln(tw, "ID\tDELETED\tSTATUS")
	fmt.Fprintf(tw, "%s\t%t\t%s\n", result.ID, result.Deleted, result.Status)
	return tw.Flush()
}

func writeJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
