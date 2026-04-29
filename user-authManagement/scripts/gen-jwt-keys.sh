#!/usr/bin/env bash
# Generates an RSA-2048 keypair for JWT (RS256) into ./secrets/.
# docker-compose mounts this directory into the container at /run/secrets,
# and the service reads the PEM files directly via JWT_*_KEY_FILE env vars.
#
# Usage: ./scripts/gen-jwt-keys.sh [out_dir]
set -euo pipefail

OUT_DIR="${1:-./secrets}"
mkdir -p "$OUT_DIR"

PRIV="$OUT_DIR/jwt_private.pem"
PUB="$OUT_DIR/jwt_public.pem"

if [[ -f "$PRIV" || -f "$PUB" ]]; then
  echo "Key files already exist in $OUT_DIR — refusing to overwrite." >&2
  echo "Delete them first if you want to regenerate." >&2
  exit 1
fi

openssl genrsa -out "$PRIV" 2048 >/dev/null 2>&1
chmod 600 "$PRIV"
openssl rsa -in "$PRIV" -pubout -out "$PUB" >/dev/null 2>&1

echo "Generated:"
echo "  $PRIV"
echo "  $PUB"
echo
echo "docker-compose mounts ./secrets/ into the container automatically."
echo "No further .env edits are required for local development."
