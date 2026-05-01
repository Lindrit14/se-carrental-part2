#!/usr/bin/env bash
set -euo pipefail

RG="${RG:-platform-rg}"
LOCATION="${LOCATION:-westeurope}"
PARAMS_FILE="${PARAMS_FILE:-parameters.dev.bicepparam}"

cd "$(dirname "$0")"

# JWT-Schlüssel aus shared-secrets/ lesen — sind nicht im Repo, müssen
# lokal mit "make keys" erzeugt sein.
JWT_PRIVATE_PATH="../shared-secrets/jwt_private.pem"
JWT_PUBLIC_PATH="../shared-secrets/jwt_public.pem"

if [[ ! -f "$JWT_PRIVATE_PATH" || ! -f "$JWT_PUBLIC_PATH" ]]; then
  echo "✗ JWT-Schlüssel fehlen unter shared-secrets/" >&2
  echo "  Erwartet: $JWT_PRIVATE_PATH und $JWT_PUBLIC_PATH" >&2
  echo "  Lösung:   cd .. && make keys" >&2
  exit 1
fi

# Werden vom .bicepparam via readEnvironmentVariable() konsumiert.
export JWT_PRIVATE_KEY="$(cat "$JWT_PRIVATE_PATH")"
export JWT_PUBLIC_KEY="$(cat "$JWT_PUBLIC_PATH")"

echo "→ Resource Group $RG..."
az group create --name "$RG" --location "$LOCATION" --only-show-errors -o none

if [[ "${1:-}" == "--what-if" ]]; then
  az deployment group what-if \
    --resource-group "$RG" \
    --template-file main.bicep \
    --parameters "$PARAMS_FILE"
  exit 0
fi

echo "→ Deploye Bicep (kann beim ersten Mal ~10 Min dauern)..."
az deployment group create \
  --resource-group "$RG" \
  --template-file main.bicep \
  --parameters "$PARAMS_FILE" \
  --output table

echo ""
echo "✓ Outputs:"
az deployment group show \
  --resource-group "$RG" \
  --name main \
  --query properties.outputs \
  --output json
