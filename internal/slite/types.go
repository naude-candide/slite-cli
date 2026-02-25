package slite

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
	Hits       []Note `json:"hits"`
	Cursor     string `json:"cursor,omitempty"`
	NextCursor string `json:"nextCursor,omitempty"`
}

type SearchResponse struct {
	Hits       []Note `json:"hits"`
	Cursor     string `json:"cursor,omitempty"`
	NextCursor string `json:"nextCursor,omitempty"`
}

type DeleteResponse struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
	Status  string `json:"status,omitempty"`
}
