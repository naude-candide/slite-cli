package slite

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const DefaultBaseURL = "https://api.slite.com"

type Config struct {
	APIKey  string
	BaseURL string
	Timeout time.Duration
	Debug   bool
}

type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
	debug      bool
}

func NewClient(cfg Config) (*Client, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("api key is required")
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = DefaultBaseURL
	}

	if cfg.Timeout <= 0 {
		cfg.Timeout = 15 * time.Second
	}

	return &Client{
		httpClient: &http.Client{Timeout: cfg.Timeout},
		baseURL:    strings.TrimRight(cfg.BaseURL, "/"),
		apiKey:     cfg.APIKey,
		debug:      cfg.Debug,
	}, nil
}

func (c *Client) Me(ctx context.Context) (*MeResponse, error) {
	var out MeResponse
	if err := c.doJSON(ctx, http.MethodGet, "/v1/me", nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ListNotes(ctx context.Context, owner, parentNoteID string, limit, offset int, cursor string) (*NotesResponse, error) {
	q := url.Values{}
	if owner != "" {
		q.Set("owner", owner)
	}
	if parentNoteID != "" {
		q.Set("parentNoteId", parentNoteID)
	}
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		q.Set("offset", strconv.Itoa(offset))
	}
	if cursor != "" {
		q.Set("cursor", cursor)
	}

	var out NotesResponse
	if err := c.doJSON(ctx, http.MethodGet, "/v1/notes", q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetNote(ctx context.Context, id string) (*NoteDetail, error) {
	var out map[string]any
	if err := c.doJSON(ctx, http.MethodGet, "/v1/notes/"+url.PathEscape(id), nil, nil, &out); err != nil {
		return nil, err
	}

	note, err := extractNoteDetail(out)
	if err != nil {
		return nil, err
	}
	return note, nil
}

func (c *Client) CreateNote(ctx context.Context, payload map[string]any) (*NoteDetail, error) {
	var out map[string]any
	if err := c.doJSON(ctx, http.MethodPost, "/v1/notes", nil, payload, &out); err != nil {
		return nil, err
	}
	return extractNoteDetail(out)
}

func (c *Client) UpdateNote(ctx context.Context, id string, payload map[string]any) (*NoteDetail, error) {
	var out map[string]any
	if err := c.doJSON(ctx, http.MethodPut, "/v1/notes/"+url.PathEscape(id), nil, payload, &out); err != nil {
		return nil, err
	}
	return extractNoteDetail(out)
}

func (c *Client) DeleteNote(ctx context.Context, id string) (*DeleteResponse, error) {
	var out map[string]any
	if err := c.doJSON(ctx, http.MethodDelete, "/v1/notes/"+url.PathEscape(id), nil, nil, &out); err != nil {
		return nil, err
	}
	return &DeleteResponse{
		ID:      firstString(out, "id"),
		Deleted: true,
		Status:  firstString(out, "status", "result"),
	}, nil
}

func (c *Client) SearchNotes(ctx context.Context, query string, limit, offset int, cursor string) (*SearchResponse, error) {
	q := url.Values{}
	q.Set("query", query)
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		q.Set("offset", strconv.Itoa(offset))
	}
	if cursor != "" {
		q.Set("cursor", cursor)
	}

	var out SearchResponse
	if err := c.doJSON(ctx, http.MethodGet, "/v1/search-notes", q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) doJSON(ctx context.Context, method, path string, query url.Values, in any, out any) error {
	endpoint := c.baseURL + path
	if len(query) > 0 {
		endpoint += "?" + query.Encode()
	}

	var reqBody []byte
	var err error
	if in != nil {
		reqBody, err = json.Marshal(in)
		if err != nil {
			return fmt.Errorf("encode request: %w", err)
		}
	}

	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		var bodyReader io.Reader
		if len(reqBody) > 0 {
			bodyReader = bytes.NewReader(reqBody)
		}

		req, err := http.NewRequestWithContext(ctx, method, endpoint, bodyReader)
		if err != nil {
			return err
		}

		req.Header.Set("x-slite-api-key", c.apiKey)
		req.Header.Set("Accept", "application/json")
		if len(reqBody) > 0 {
			req.Header.Set("Content-Type", "application/json")
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			if !sleepBackoff(ctx, attempt) {
				break
			}
			continue
		}

		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			return readErr
		}

		if c.debug {
			_, _ = fmt.Fprintf(os.Stderr, "[%d] %s %s\n", resp.StatusCode, method, endpoint)
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			if len(body) == 0 || out == nil {
				return nil
			}
			if err := json.Unmarshal(body, out); err != nil {
				return fmt.Errorf("decode response: %w", err)
			}
			return nil
		}

		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
			lastErr = fmt.Errorf("slite api temporary error: status=%d body=%s", resp.StatusCode, truncate(body, 300))
			if !sleepBackoff(ctx, attempt) {
				break
			}
			continue
		}

		return fmt.Errorf("slite api error: status=%d body=%s", resp.StatusCode, truncate(body, 300))
	}

	if lastErr != nil {
		return lastErr
	}
	return fmt.Errorf("request failed")
}

func extractNoteDetail(payload map[string]any) (*NoteDetail, error) {
	candidates := []map[string]any{payload}
	for _, key := range []string{"note", "item", "data", "hit"} {
		if v, ok := payload[key]; ok {
			if m, ok := v.(map[string]any); ok {
				candidates = append(candidates, m)
			}
		}
	}

	for _, candidate := range candidates {
		id := firstString(candidate, "id")
		if id == "" {
			continue
		}

		owner := firstString(candidate, "ownerId", "ownerID", "owner")
		if owner == "" {
			if ownerObj, ok := candidate["owner"].(map[string]any); ok {
				owner = firstString(ownerObj, "id", "name")
			}
		}

		return &NoteDetail{
			ID:        id,
			Title:     firstString(candidate, "title"),
			OwnerID:   owner,
			UpdatedAt: firstString(candidate, "updatedAt", "updated_at", "updated"),
			URL:       firstString(candidate, "url", "link"),
			Markdown:  firstString(candidate, "markdown", "content", "body"),
		}, nil
	}

	return nil, fmt.Errorf("unexpected note response shape")
}

func firstString(m map[string]any, keys ...string) string {
	for _, key := range keys {
		v, ok := m[key]
		if !ok || v == nil {
			continue
		}
		switch t := v.(type) {
		case string:
			return t
		case fmt.Stringer:
			return t.String()
		}
	}
	return ""
}

func sleepBackoff(ctx context.Context, attempt int) bool {
	wait := time.Duration(1<<attempt) * 300 * time.Millisecond
	t := time.NewTimer(wait)
	defer t.Stop()

	select {
	case <-ctx.Done():
		return false
	case <-t.C:
		return true
	}
}

func truncate(b []byte, n int) string {
	s := string(b)
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
