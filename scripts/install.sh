#!/usr/bin/env bash
set -euo pipefail

REPO="${REPO:-naude-candide/slite-cli}"
VERSION="${VERSION:-latest}"
BIN_NAME="slite"
ARCHIVE_NAME=""

case "$(uname -s)" in
  Darwin) OS="darwin" ;;
  Linux)  OS="linux"  ;;
  *)
    echo "Unsupported OS: $(uname -s)" >&2
    exit 1
    ;;
esac

case "$(uname -m)" in
  arm64|aarch64) ARCH_SUFFIX="arm64" ;;
  x86_64)        ARCH_SUFFIX="amd64" ;;
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

ARCHIVE_NAME="${BIN_NAME}-${OS}-${ARCH_SUFFIX}.tar.gz"
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
if ! curl -fL --retry 3 --connect-timeout 10 -o "$ARCHIVE_PATH" "$DOWNLOAD_URL"; then
  if command -v gh >/dev/null 2>&1; then
    echo "Direct download failed; trying authenticated GitHub CLI download."
    if [[ "$VERSION" == "latest" ]]; then
      gh release download --repo "$REPO" --pattern "$ARCHIVE_NAME" --dir "$TMP_DIR"
    else
      gh release download "$VERSION" --repo "$REPO" --pattern "$ARCHIVE_NAME" --dir "$TMP_DIR"
    fi
    if [[ ! -f "$ARCHIVE_PATH" ]]; then
      echo "Failed to download ${ARCHIVE_NAME} from GitHub release." >&2
      exit 1
    fi
  else
    echo "Download failed and GitHub CLI (gh) is not available for authenticated fallback." >&2
    exit 1
  fi
fi

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

prompt_for_api_key() {
  local tty_device="/dev/tty"
  if [[ "${SKIP_API_KEY_PROMPT:-0}" == "1" ]]; then
    return
  fi
  if [[ -n "${SLITE_API_KEY:-}" ]]; then
    echo "SLITE_API_KEY is already set in this shell."
    return
  fi
  if [[ ! -r "$tty_device" ]]; then
    echo "Non-interactive shell detected; skipping API key prompt."
    return
  fi

  echo
  echo "Set up authentication now."
  read -r -s -p "Enter your Slite API key (leave blank to skip): " entered_key < "$tty_device"
  echo
  if [[ -z "${entered_key}" ]]; then
    echo "Skipped API key setup."
    return
  fi

  export SLITE_API_KEY="${entered_key}"
  echo "SLITE_API_KEY exported for current shell."

  read -r -p "Persist to ~/.zshrc for future shells? [Y/n]: " persist_answer < "$tty_device"
  case "${persist_answer:-Y}" in
    y|Y|yes|YES|"")
      if grep -q '^export SLITE_API_KEY=' "${HOME}/.zshrc" 2>/dev/null; then
        tmp_file="$(mktemp)"
        awk -v key="${entered_key}" '
          BEGIN { done=0 }
          /^export SLITE_API_KEY=/ {
            if (!done) {
              print "export SLITE_API_KEY=" key
              done=1
            }
            next
          }
          { print }
          END {
            if (!done) print "export SLITE_API_KEY=" key
          }
        ' "${HOME}/.zshrc" > "${tmp_file}"
        mv "${tmp_file}" "${HOME}/.zshrc"
      else
        printf '\nexport SLITE_API_KEY=%s\n' "${entered_key}" >> "${HOME}/.zshrc"
      fi
      echo "Saved to ~/.zshrc"
      ;;
    *)
      echo "Not saved to ~/.zshrc"
      ;;
  esac
}

prompt_for_api_key

if [[ ":$PATH:" != *":${TARGET_DIR}:"* ]]; then
  echo
  echo "Note: ${TARGET_DIR} is not currently in your PATH."
  echo "Add this to your shell profile (~/.zshrc):"
  echo "  export PATH=\"${TARGET_DIR}:\$PATH\""
fi
