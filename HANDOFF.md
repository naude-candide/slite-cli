# Slite CLI Handoff

Last updated: 2026-02-25

## Current status

Go CLI MVP in this folder now includes:
- `slite me`
- `slite docs list`
- `slite docs get <id>`
- `slite search <query>`

Implemented components:
- Cobra command tree in `cmd/`
- API key loading from `SLITE_API_KEY`
- Slite HTTP client with:
  - `x-slite-api-key` auth header
  - retries/backoff for `429` and `5xx`
  - timeout/base URL/debug support
- list + search + single-doc fetch support
- table output and `--json` output mode

## Files updated in latest step

- `cmd/docs.go` (added `docs get <id>` command)
- `internal/slite/client.go` (added `GetNote`, response normalization helpers)
- `internal/slite/types.go` (added `NoteDetail`)
- `internal/output/output.go` (added `RenderNote`)
- `README.md` (documented new command)

## Known blockers in this environment

- `go` is not installed (`go: command not found`)
- `gofmt` is not installed (`gofmt: command not found`)
- This folder is not a git repository (`fatal: not a git repository`)

Because of the above, compile/test/format checks were not executed here.

## Resume checklist

1. Install Go (1.22+ recommended).
2. Run:

```bash
cd ~/apps/cli/slite-cli
go mod tidy
gofmt -w ./...
go build -o slite .
```

3. Set auth key:

```bash
export SLITE_API_KEY=your_key_here
```

4. Smoke test:

```bash
./slite me
./slite docs list --limit 5
./slite docs get <note_id>
./slite search "roadmap" --limit 5
```

## Suggested next improvements

- Add `docs create/update/delete`
- Add pagination helpers and cursor support if API returns cursors
- Add tests for `internal/slite` using mock HTTP server
- Add release flow (GitHub Actions + goreleaser)
