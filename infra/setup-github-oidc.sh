#!/usr/bin/env bash
set -euo pipefail

# =====================================================================
# Einmaliges Azure-Setup für GitHub Actions OIDC
#
# Was passiert:
#   1. Azure AD App + Service Principal anlegen
#   2. Contributor-Rolle auf der Resource Group + AcrPush auf der ACR
#   3. Federated Credential für GitHub-Repo registrieren
#         (push auf main + manueller workflow_dispatch)
#
# Output: 3 Werte für die GitHub-Repo-Secrets
#   - AZURE_CLIENT_ID
#   - AZURE_TENANT_ID
#   - AZURE_SUBSCRIPTION_ID
#
# Voraussetzungen: az login, ausreichende Rechte für AD-App-Erstellung
# (User Access Administrator + Application Developer reichen).
# =====================================================================

GITHUB_OWNER="${GITHUB_OWNER:-Lindrit14}"
GITHUB_REPO="${GITHUB_REPO:-se-carrental-part2}"
RG="${RG:-platform-rg}"
ACR_NAME="${ACR_NAME:-platformacrlp13}"
APP_NAME="${APP_NAME:-github-actions-microservices}"

echo "→ Setup für $GITHUB_OWNER/$GITHUB_REPO → $RG / $ACR_NAME"

# ---- 1. Azure AD App + Service Principal -----------------------------
EXISTING_APP_ID=$(az ad app list --display-name "$APP_NAME" --query "[0].appId" -o tsv 2>/dev/null || true)

if [[ -n "$EXISTING_APP_ID" ]]; then
  echo "  AD App '$APP_NAME' existiert bereits ($EXISTING_APP_ID), nutze sie weiter."
  APP_ID="$EXISTING_APP_ID"
else
  echo "→ Erstelle AD App '$APP_NAME'..."
  APP_ID=$(az ad app create --display-name "$APP_NAME" --query appId -o tsv)
  echo "  appId: $APP_ID"
fi

EXISTING_SP_ID=$(az ad sp list --filter "appId eq '$APP_ID'" --query "[0].id" -o tsv 2>/dev/null || true)
if [[ -z "$EXISTING_SP_ID" ]]; then
  echo "→ Erstelle Service Principal..."
  az ad sp create --id "$APP_ID" --query id -o tsv > /dev/null
fi

# ---- 2. Rollen-Zuweisungen -------------------------------------------
RG_ID=$(az group show -n "$RG" --query id -o tsv)
ACR_ID=$(az acr show -n "$ACR_NAME" --query id -o tsv)

echo "→ Contributor auf Resource Group..."
az role assignment create \
  --assignee "$APP_ID" \
  --role Contributor \
  --scope "$RG_ID" \
  --only-show-errors -o none 2>/dev/null || echo "  (bereits zugewiesen)"

echo "→ AcrPush auf ACR..."
az role assignment create \
  --assignee "$APP_ID" \
  --role AcrPush \
  --scope "$ACR_ID" \
  --only-show-errors -o none 2>/dev/null || echo "  (bereits zugewiesen)"

# ---- 3. Federated Credentials ----------------------------------------
# Eines für push-to-main, eines für workflow_dispatch (separater subject).
echo "→ Federated Credential 'github-main' (push auf main)..."
az ad app federated-credential create --id "$APP_ID" --parameters "{
  \"name\": \"github-main\",
  \"issuer\": \"https://token.actions.githubusercontent.com\",
  \"subject\": \"repo:${GITHUB_OWNER}/${GITHUB_REPO}:ref:refs/heads/main\",
  \"audiences\": [\"api://AzureADTokenExchange\"]
}" --only-show-errors -o none 2>/dev/null || echo "  (existiert bereits)"

echo "→ Federated Credential 'github-dispatch' (manueller Trigger)..."
# workflow_dispatch fired ohne ref-binding → wir nutzen das pull_request-Pattern
# nicht, sondern lassen für jetzt nur main gelten. Für Dispatch ohne ref kannst
# du später ein zweites FC mit subject "repo:.../...:environment:production"
# hinzufügen und im Workflow environment: production setzen.

# ---- Output ----------------------------------------------------------
echo ""
echo "════════════════════════════════════════════════════════════════"
echo "  Werte für GitHub Repo Secrets (Settings → Secrets → Actions)"
echo "════════════════════════════════════════════════════════════════"
echo ""
echo "  AZURE_CLIENT_ID       = $APP_ID"
echo "  AZURE_TENANT_ID       = $(az account show --query tenantId -o tsv)"
echo "  AZURE_SUBSCRIPTION_ID = $(az account show --query id -o tsv)"
echo ""
echo "  Zusätzlich noch ein Secret für den Frontend-Build:"
echo "  VITE_GOOGLE_MAPS_API_KEY = (aus deiner lokalen .env)"
echo ""
echo "════════════════════════════════════════════════════════════════"
