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

## Install (macOS)

From GitHub Releases:

```bash
curl -fsSL https://raw.githubusercontent.com/naude-candide/slite-cli/main/scripts/install.sh | bash
```

Optional env vars:

- `VERSION=v0.1.0` to install a specific release tag (default: latest)
- `INSTALL_DIR=$HOME/bin` to override install location
- `REPO=owner/repo` to install from a different repository
- `SKIP_API_KEY_PROMPT=1` to disable interactive API key setup

Example:

```bash
VERSION=v0.1.0 INSTALL_DIR=$HOME/bin \
curl -fsSL https://raw.githubusercontent.com/naude-candide/slite-cli/main/scripts/install.sh | bash
```

By default, the installer prompts for `SLITE_API_KEY` and can persist it to `~/.zshrc`.

Release assets expected by installer:
- `slite-darwin-arm64.tar.gz`
- `slite-darwin-amd64.tar.gz`

To build these release archives locally:

```bash
chmod +x scripts/build-release.sh
scripts/build-release.sh v0.1.0
```

This writes tarballs to `dist/` that you can upload to a GitHub release tag.

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
