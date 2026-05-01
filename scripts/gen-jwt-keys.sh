#!/usr/bin/env bash
# Generates an RSA-2048 keypair for JWT (RS256) into the given output dir
# (defaults to ./shared-secrets/). user-auth mounts this dir at /run/secrets
# and reads the PEM files via JWT_*_KEY_FILE env vars; the other services
# mount the same dir read-only and consume jwt_public.pem for verification.
#
# Usage: ./scripts/gen-jwt-keys.sh [out_dir]
set -euo pipefail

OUT_DIR="${1:-./shared-secrets}"
mkdir -p "$OUT_DIR"

PRIV="$OUT_DIR/jwt_private.pem"
PUB="$OUT_DIR/jwt_public.pem"

# Guard on the private key only — the public is regenerated from it,
# and the committed shared-secrets/jwt_public.pem always exists.
if [[ -f "$PRIV" ]]; then
  echo "Private key already exists at $PRIV — refusing to overwrite." >&2
  echo "Delete it first if you want to regenerate the keypair." >&2
  exit 1
fi

openssl genrsa -out "$PRIV" 2048 >/dev/null 2>&1
chmod 600 "$PRIV"
openssl rsa -in "$PRIV" -pubout -out "$PUB" >/dev/null 2>&1

echo "Generated:"
echo "  $PRIV"
echo "  $PUB"
