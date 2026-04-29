# Azure Deployment

Provisions the user-auth microservice on **Azure Container Apps** with a
self-hosted RabbitMQ side-car and image storage in **Azure Container Registry**.

## Prerequisites

- `az` CLI logged in (`az login`)
- Subscription set: `az account set --subscription <id>`
- Resource group: `az group create -n rg-user-auth -l westeurope`
- For Mongo in production: provision **Azure Cosmos DB for MongoDB API**
  separately and obtain its connection string.

## 1. Container Registry

```bash
az deployment group create \
  -g rg-user-auth \
  -f acr.bicep \
  -p registryName=userauthacr$RANDOM
```

Note the `loginServer` output (e.g. `userauthacr1234.azurecr.io`).

## 2. Container Apps Environment

A Container Apps Environment must exist before deploying the app:

```bash
az containerapp env create \
  -n cae-user-auth \
  -g rg-user-auth \
  -l westeurope
```

Capture its resource id:

```bash
ENV_ID=$(az containerapp env show -n cae-user-auth -g rg-user-auth --query id -o tsv)
```

## 3. Build & push the image

```bash
TAG=$(git rev-parse --short HEAD)
ACR=userauthacr1234.azurecr.io      # from step 1
az acr login -n ${ACR%.azurecr.io}

docker build -f ../docker/Dockerfile -t $ACR/user-auth:$TAG ../..
docker push $ACR/user-auth:$TAG
```

## 4. Generate JWT keys (once per environment)

```bash
../../scripts/gen-jwt-keys.sh ./secrets
PRIV=$(cat ./secrets/jwt_private.pem)
PUB=$(cat ./secrets/jwt_public.pem)
```

Treat these as secrets — store them in **Azure Key Vault** for prod and
reference them via `secretRef` instead of inline parameters.

## 5. Deploy the Container Apps

```bash
ACR_USER=$(az acr credential show -n ${ACR%.azurecr.io} --query username -o tsv)
ACR_PASS=$(az acr credential show -n ${ACR%.azurecr.io} --query 'passwords[0].value' -o tsv)
MONGO_URI="mongodb+srv://...azure.com/auth?tls=true"   # from Cosmos DB

az deployment group create \
  -g rg-user-auth \
  -f containerapp.bicep \
  -p environmentId="$ENV_ID" \
     registryServer="$ACR" \
     registryUsername="$ACR_USER" \
     registryPassword="$ACR_PASS" \
     imageTag="$TAG" \
     mongoUri="$MONGO_URI" \
     jwtPrivateKey="$PRIV" \
     jwtPublicKey="$PUB"
```

The deployment outputs `authFqdn` — the public URL of the service.

## 6. Smoke test

```bash
FQDN=$(az containerapp show -n auth-service -g rg-user-auth --query properties.configuration.ingress.fqdn -o tsv)
curl https://$FQDN/healthz
```

## Updates

Subsequent rollouts are a single command:

```bash
TAG=$(git rev-parse --short HEAD)
docker build -f ../docker/Dockerfile -t $ACR/user-auth:$TAG ../..
docker push $ACR/user-auth:$TAG
az containerapp update -n auth-service -g rg-user-auth --image $ACR/user-auth:$TAG
```

## Notes

- **Mongo**: For Prod use Azure Cosmos DB for MongoDB API — managed,
  geo-redundant, automatic backups. The self-hosted Mongo from
  `docker-compose.yml` is for local development only.
- **RabbitMQ**: The Bicep deploys a single replica with no persistence.
  For Prod consider CloudAMQP (Managed) or attach Azure Files to
  `/var/lib/rabbitmq` via a `volumeMounts` block.
- **Secrets**: For real environments, replace inline secret parameters
  with Key Vault references (`keyVaultUrl`).
