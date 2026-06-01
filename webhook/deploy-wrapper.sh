#!/usr/bin/env bash
set -euo pipefail
TAG="${1#refs/tags/}"
exec /home/uc/InventoryManagementSystem/deploy.sh "$TAG"
