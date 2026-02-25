# slite-cli

Small Go CLI wrapper around the Slite API.

## Commands

- `slite me`
- `slite docs list [--owner <id>] [--limit 20] [--offset 0]`
- `slite docs get <id>`
- `slite search <query> [--limit 20] [--offset 0]`

## Auth

Set your API key in the environment:

```bash
export SLITE_API_KEY=your_api_key
```

## Build

```bash
go mod tidy
go build -o slite .
```

## Examples

```bash
./slite me
./slite docs list --limit 10
./slite docs get abc123
./slite search "product roadmap" --json
```

## Global flags

- `--json` output JSON
- `--debug` print status + URL to stderr
- `--base-url` override API base (default `https://api.slite.com`)
- `--timeout` request timeout (default `15s`)

See `HANDOFF.md` for project status and resume steps.
