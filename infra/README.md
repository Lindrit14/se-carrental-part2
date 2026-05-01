# Infra: Azure Container Apps Deployment

Deployt die komplette Plattform aus `../docker-compose.yml` als Azure Container Apps. Eine einzige `main.bicep`-Datei orchestriert alles.

## Was wird deployed?

- **1× Log Analytics Workspace** — sammelt Logs aller Container Apps
- **1× Azure Container Registry (ACR, Basic SKU)** — hostet später die Images der eigenen Services
- **1× Storage Account + 5× File Shares** — für die DB-Volumes
- **1× Container Apps Environment** — Hosting-Plane mit internem DNS
- **5× DB-Container-Apps** mit Volume-Mount: `rabbitmq`, `redis`, `user-auth-mongo`, `cars-db`, `booking-db`
- **7× App-Container-Apps**: `api-gateway` (public), `user-auth`, `car-service`, `booking`, `currency-converter` (gRPC/HTTP2), `notification`, `frontend` (public)

**Keine** Managed Services (kein Azure DB, kein Cosmos, kein Redis Cache).

## Voraussetzungen

```bash
az --version          # >= 2.50
az bicep version      # neueste Version
az login              # falls nicht eingeloggt
az account set --subscription <YOUR-SUBSCRIPTION-ID>
```

**JWT-Schlüssel müssen lokal vorhanden sein:**

```bash
ls ../shared-secrets/jwt_private.pem ../shared-secrets/jwt_public.pem
# Falls nicht vorhanden:
cd .. && make keys
```

Die zwei `.pem`-Files werden vom Deploy-Skript gelesen und als Container Apps Secrets unter `/run/secrets/jwt_*.pem` in api-gateway, user-auth, car-service und booking gemountet — analog zum lokalen Compose-Setup. Sie verlassen deinen lokalen Rechner als Bicep-`@secure()` Parameter; das Repo enthält nur den Public Key (Private ist gitignored).

## 1. Parameter editieren

`parameters.dev.bicepparam`:

```bicep
using 'main.bicep'
param prefix = 'platform'
param uniqueSuffix = 'CHANGEME'  // ← anpassen!
```

`uniqueSuffix` fließt in den ACR- und Storage-Account-Namen ein. Die müssen global eindeutig sein → Initialen + Zahl, lowercase, max 8 Zeichen (z.B. `lp042`).

## 2. Deployen

```bash
./deploy.sh --what-if   # Dry-Run, zeigt was erstellt wird
./deploy.sh             # Echter Deploy (~10 Min beim ersten Mal)
```

Outputs am Ende:
- `frontendUrl` — öffentliche Frontend-URL
- `apiGatewayUrl` — öffentliche API-Gateway-URL
- `acrLoginServer`, `acrName` — für CI/CD

## 3. Wieder löschen

```bash
az group delete --name platform-rg --yes --no-wait
```

## Wichtige Hinweise

- **Placeholder-Image:** Die 7 eigenen Services starten initial mit `mcr.microsoft.com/azuredocs/containerapps-helloworld:latest`. CI/CD baut die echten Images, pusht sie in die ACR und updated die Apps. Bis dahin zeigen `frontendUrl` und `apiGatewayUrl` einen Hello-World-Screen.
- **Frontend `VITE_API_URL`:** Vite-Apps backen die API-URL beim Build ins Image. Der CI/CD-Build des Frontends muss `VITE_API_URL=https://<api-gateway-fqdn>` setzen, sonst spricht das Frontend mit `http://localhost:8080` (so wie aktuell in der `compose.yml`).
- **JWT-Keys:** Werden aus `../shared-secrets/{jwt_private,jwt_public}.pem` gelesen und als Container Apps Secrets in 4 Apps gemountet. Wenn die Files fehlen → `make keys` im Repo-Root ausführen.
- **CI/CD nimmt nach erstem Deploy übers Image-Lifecycle:** Sobald GitHub Actions die Apps mit echten Images updatet, würde ein erneuter `./deploy.sh`-Lauf die Apps zurück auf das Placeholder-Image setzen. Bicep ist Tag-1-only — für Image-Updates GitHub Actions nehmen.
- **gRPC-Adresse:** `currency-converter` läuft mit `transport=http2` und `allowInsecure=true`. Spring-Clients sprechen ihn unter `currency-converter:80` an. Falls TLS-gRPC bevorzugt wird, im Bicep `allowInsecure=false` und im Spring-Client `:443` verwenden.
- **Persistenz:** DBs überleben Neustart und Re-Deploy, weil sie auf Azure File Shares schreiben. Wenn das Storage Account aber gelöscht wird, sind die Daten weg.
- **Idempotenz:** Wiederholte `./deploy.sh`-Calls sind unschädlich, sofern sich `main.bicep` und `parameters.dev.bicepparam` nicht ändern.
