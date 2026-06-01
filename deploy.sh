#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "$0")"

TAG="${1:?Usage: $0 <tag>}"

echo "[$(date -Iseconds)] Deploying tag $TAG"

git fetch --tags

git checkout "$TAG"

docker compose -f docker-compose.prod.yml build

docker compose -f docker-compose.prod.yml up -d

docker image prune -f

EXITED=$(docker compose -f docker-compose.prod.yml ps --format json \
  | jq -r 'select(.State == "exited") | .Name' || true)

if [ -n "$EXITED" ]; then

  echo "WARN: the following containers exited:"

  echo "$EXITED"
  for c in $EXITED; do

    echo "--- logs for $c ---"

    docker logs --tail 50 "$c" || true
  done
fi

echo "[$(date -Iseconds)] Deploy complete"
