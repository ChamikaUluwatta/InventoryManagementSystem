#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "$0")"

TAG="${1:?Usage: $0 <tag>}"

echo "[$(date -Iseconds)] Deploying tag $TAG"

git fetch --tags

git checkout -f "$TAG"

docker compose -f docker-compose.prod.yml run --rm migrator

docker compose -f compose.yaml up -d --build --remove-orphans

echo "[$(date -Iseconds)] Deploy complete"
