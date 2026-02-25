#!/usr/bin/env bash
set -euo pipefail

VERSION="${1:-dev}"
OUT_DIR="${OUT_DIR:-dist}"
BIN_NAME="slite"

mkdir -p "$OUT_DIR"

build_one() {
  local goarch="$1"
  local arch_suffix="$2"
  local work_dir
  work_dir="$(mktemp -d)"
  trap 'rm -rf "$work_dir"' RETURN

  local bin_path="${work_dir}/${BIN_NAME}"
  local archive_path="${OUT_DIR}/${BIN_NAME}-darwin-${arch_suffix}.tar.gz"

  GOOS=darwin GOARCH="$goarch" CGO_ENABLED=0 go build -ldflags "-s -w -X main.version=${VERSION}" -o "$bin_path" .
  tar -C "$work_dir" -czf "$archive_path" "$BIN_NAME"
  shasum -a 256 "$archive_path"
}

build_one "arm64" "arm64"
build_one "amd64" "amd64"

echo "Artifacts created in ${OUT_DIR}/"
