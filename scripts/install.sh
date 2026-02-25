#!/usr/bin/env bash
set -euo pipefail

REPO="${REPO:-naude-candide/slite-cli}"
VERSION="${VERSION:-latest}"
BIN_NAME="slite"
ARCHIVE_NAME=""

if [[ "$(uname -s)" != "Darwin" ]]; then
  echo "This installer currently supports macOS only." >&2
  exit 1
fi

case "$(uname -m)" in
  arm64) ARCH_SUFFIX="arm64" ;;
  x86_64) ARCH_SUFFIX="amd64" ;;
  *)
    echo "Unsupported CPU architecture: $(uname -m)" >&2
    exit 1
    ;;
esac

if [[ "$VERSION" == "latest" ]]; then
  DOWNLOAD_BASE="https://github.com/${REPO}/releases/latest/download"
else
  DOWNLOAD_BASE="https://github.com/${REPO}/releases/download/${VERSION}"
fi

ARCHIVE_NAME="${BIN_NAME}-darwin-${ARCH_SUFFIX}.tar.gz"
DOWNLOAD_URL="${DOWNLOAD_BASE}/${ARCHIVE_NAME}"

if [[ -n "${INSTALL_DIR:-}" ]]; then
  TARGET_DIR="$INSTALL_DIR"
elif [[ -w "/usr/local/bin" ]]; then
  TARGET_DIR="/usr/local/bin"
else
  TARGET_DIR="${HOME}/.local/bin"
fi

mkdir -p "$TARGET_DIR"

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

ARCHIVE_PATH="${TMP_DIR}/${ARCHIVE_NAME}"
BIN_PATH="${TMP_DIR}/${BIN_NAME}"

echo "Downloading ${DOWNLOAD_URL}"
curl -fL --retry 3 --connect-timeout 10 -o "$ARCHIVE_PATH" "$DOWNLOAD_URL"

tar -xzf "$ARCHIVE_PATH" -C "$TMP_DIR"

if [[ ! -f "$BIN_PATH" ]]; then
  echo "Archive did not contain '${BIN_NAME}' binary." >&2
  exit 1
fi

chmod +x "$BIN_PATH"
mv "$BIN_PATH" "${TARGET_DIR}/${BIN_NAME}"

echo "Installed ${BIN_NAME} to ${TARGET_DIR}/${BIN_NAME}"
echo
echo "Run: ${BIN_NAME} --help"

if [[ ":$PATH:" != *":${TARGET_DIR}:"* ]]; then
  echo
  echo "Note: ${TARGET_DIR} is not currently in your PATH."
  echo "Add this to your shell profile (~/.zshrc):"
  echo "  export PATH=\"${TARGET_DIR}:\$PATH\""
fi
