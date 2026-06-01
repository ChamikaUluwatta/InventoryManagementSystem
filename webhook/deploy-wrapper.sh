#!/usr/bin/env bash
set -euo pipefail
TAG="${1#refs/tags/}"
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
exec "$ROOT_DIR/deploy.sh" "$TAG"
