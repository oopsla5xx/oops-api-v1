#!/bin/bash
# Run seed files for a given environment.
# Usage: ./scripts/seed.sh [development|test]
set -euo pipefail

ENV=${1:-development}
SEED_DIR="database/seeds/${ENV}"

if [ ! -d "$SEED_DIR" ]; then
  echo "Seed directory not found: $SEED_DIR"
  exit 1
fi

# Load env file
ENV_FILE=".env.${ENV}"
if [ -f "$ENV_FILE" ]; then
  # shellcheck disable=SC1090
  set -a && source "$ENV_FILE" && set +a
fi

if [ -z "${DATABASE_DSN:-}" ]; then
  echo "DATABASE_DSN is not set"
  exit 1
fi

echo "Running seeds for environment: ${ENV}"

for f in "${SEED_DIR}"/*.sql; do
  [ -f "$f" ] || continue
  echo "  → $(basename "$f")"
  psql "$DATABASE_DSN" -f "$f"
done

echo "Seeds completed."
