package slite

import "encoding/json"

type MeResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Organization struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"organization"`
}

type Note struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	OwnerID   string `json:"ownerId"`
	ParentID  string `json:"parentNoteId"`
	UpdatedAt string `json:"updatedAt"`
}

type NoteDetail struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	OwnerID   string `json:"ownerId"`
	UpdatedAt string `json:"updatedAt"`
	URL       string `json:"url"`
}

type NotesResponse struct {
	Hits       []Note `json:"hits,omitempty"`
	Cursor     string `json:"cursor,omitempty"`
	NextCursor string `json:"nextCursor,omitempty"`
	HasNext    bool   `json:"hasNextPage,omitempty"`
	Total      int    `json:"total,omitempty"`
}

type SearchResponse struct {
	Hits       []Note `json:"hits,omitempty"`
	Cursor     string `json:"cursor,omitempty"`
	NextCursor string `json:"nextCursor,omitempty"`
	HasNext    bool   `json:"hasNextPage,omitempty"`
	Total      int    `json:"total,omitempty"`
}

type DeleteResponse struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
	Status  string `json:"status,omitempty"`
}

type notesResponseAlias struct {
	Hits       []Note `json:"hits"`
	Notes      []Note `json:"notes"`
	Cursor     string `json:"cursor"`
	NextCursor string `json:"nextCursor"`
	HasNext    bool   `json:"hasNextPage"`
	Total      int    `json:"total"`
}

func (r *NotesResponse) UnmarshalJSON(data []byte) error {
	var raw notesResponseAlias
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	r.Cursor = raw.Cursor
	r.NextCursor = raw.NextCursor
	r.HasNext = raw.HasNext
	r.Total = raw.Total
	if len(raw.Notes) > 0 {
		r.Hits = raw.Notes
	} else {
		r.Hits = raw.Hits
	}
	return nil
}

func (r *SearchResponse) UnmarshalJSON(data []byte) error {
	var raw notesResponseAlias
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	r.Cursor = raw.Cursor
	r.NextCursor = raw.NextCursor
	r.HasNext = raw.HasNext
	r.Total = raw.Total
	if len(raw.Notes) > 0 {
		r.Hits = raw.Notes
	} else {
		r.Hits = raw.Hits
	}
	return nil
}
