// =====================================================================
// Plattform-Deployment: Azure Container Apps
//
// Erzeugt die komplette Plattform aus docker-compose.yml als
// Container Apps in einer einzigen Environment. Auch die Datenbanken
// laufen als Container Apps mit Azure Files Volume Mounts —
// keinerlei Managed Services.
//
// Reihenfolge:
//   1. Basis-Infrastruktur (ACR, Storage, Environment, Log Analytics)
//   2. 5x Storage-Mounts in der Environment registrieren
//   3. 5x Datenbank-Container-Apps (mit Volume Mount)
//   4. 7x Eigene Service-Container-Apps (Placeholder-Image)
//   5. Outputs
// =====================================================================

@description('Basisname für alle Resourcen')
param prefix string

@description('Eindeutiges Suffix (lowercase, max ~8 Zeichen) — global eindeutig für ACR/Storage')
param uniqueSuffix string

@description('Azure-Region')
param location string = resourceGroup().location

@description('JWT Private Key (PEM-Inhalt). Wird als Container Apps Secret unter /run/secrets/jwt_private.pem gemountet. Aus shared-secrets/jwt_private.pem lesen.')
@secure()
param jwtPrivateKey string

@description('JWT Public Key (PEM-Inhalt). Wird als Container Apps Secret unter /run/secrets/jwt_public.pem gemountet. Aus shared-secrets/jwt_public.pem lesen.')
@secure()
param jwtPublicKey string

// ---- Konstanten ------------------------------------------------------
// Placeholder-Image für die noch nicht in der ACR liegenden eigenen
// Services. CI/CD ersetzt das später, sobald die Images gebuildet sind.
var placeholderImage = 'mcr.microsoft.com/azuredocs/containerapps-helloworld:latest'

// ---- JWT Secret-Files Helper -----------------------------------------
// Wird an die 4 Apps weitergereicht, die JWT verifizieren/signieren.
// Pfade matchen exakt das Compose-Setup, sodass die JWT_*_FILE Env-Vars
// gleich bleiben können. Object-Form weil Bicep @secure nur auf
// string/object erlaubt, nicht auf array.
var jwtPublicOnly = {
  'jwt-public': {
    fileName: 'jwt_public.pem'
    value: jwtPublicKey
  }
}
var jwtBoth = {
  'jwt-public': {
    fileName: 'jwt_public.pem'
    value: jwtPublicKey
  }
  'jwt-private': {
    fileName: 'jwt_private.pem'
    value: jwtPrivateKey
  }
}

// =====================================================================
// 1. BASIS-INFRASTRUKTUR
// =====================================================================
module infra 'modules/infrastructure.bicep' = {
  name: 'infrastructure'
  params: {
    prefix: prefix
    uniqueSuffix: uniqueSuffix
    location: location
  }
}

// =====================================================================
// 2. STORAGE-MOUNTS (Azure File Shares in der Environment registrieren)
// =====================================================================
// Die Mounts müssen existieren BEVOR die DB-Container-Apps deployed
// werden. Bicep löst das normalerweise automatisch über die
// symbolic-references, hier explizit dependsOn weil wir auf zwei
// Sub-Resourcen aus dem infra-Modul gleichzeitig warten.

module rabbitmqStorageMount 'modules/env-storage.bicep' = {
  name: 'storage-mount-rabbitmq'
  dependsOn: [
    infra
  ]
  params: {
    environmentName: infra.outputs.environmentName
    storageMountName: 'rabbitmq-mount'
    storageAccountName: infra.outputs.storageAccountName
    storageAccountKey: infra.outputs.storageAccountKey
    fileShareName: 'rabbitmq-data'
  }
}

module redisStorageMount 'modules/env-storage.bicep' = {
  name: 'storage-mount-redis'
  dependsOn: [
    infra
  ]
  params: {
    environmentName: infra.outputs.environmentName
    storageMountName: 'redis-mount'
    storageAccountName: infra.outputs.storageAccountName
    storageAccountKey: infra.outputs.storageAccountKey
    fileShareName: 'redis-data'
  }
}

module mongoStorageMount 'modules/env-storage.bicep' = {
  name: 'storage-mount-mongo'
  dependsOn: [
    infra
  ]
  params: {
    environmentName: infra.outputs.environmentName
    storageMountName: 'mongo-mount'
    storageAccountName: infra.outputs.storageAccountName
    storageAccountKey: infra.outputs.storageAccountKey
    fileShareName: 'mongo-data'
  }
}

module carsDbStorageMount 'modules/env-storage.bicep' = {
  name: 'storage-mount-cars-db'
  dependsOn: [
    infra
  ]
  params: {
    environmentName: infra.outputs.environmentName
    storageMountName: 'cars-db-mount'
    storageAccountName: infra.outputs.storageAccountName
    storageAccountKey: infra.outputs.storageAccountKey
    fileShareName: 'cars-db-data'
  }
}

module bookingDbStorageMount 'modules/env-storage.bicep' = {
  name: 'storage-mount-booking-db'
  dependsOn: [
    infra
  ]
  params: {
    environmentName: infra.outputs.environmentName
    storageMountName: 'booking-db-mount'
    storageAccountName: infra.outputs.storageAccountName
    storageAccountKey: infra.outputs.storageAccountKey
    fileShareName: 'booking-db-data'
  }
}

// =====================================================================
// 3. DATENBANK-CONTAINER-APPS (mit persistentem Volume)
// =====================================================================
// Alle DBs nutzen ingress=internal, transport=tcp, exposedPort=targetPort.
// minReplicas=1 verhindert Scale-to-Zero (DBs sollen immer laufen).

// ---- RabbitMQ --------------------------------------------------------
module rabbitmq 'modules/container-app-with-volume.bicep' = {
  name: 'app-rabbitmq'
  params: {
    name: 'rabbitmq'
    location: location
    environmentId: infra.outputs.environmentId
    image: 'rabbitmq:3.13-management-alpine'
    targetPort: 5672
    ingressType: 'internal'
    transport: 'tcp'
    exposedPort: 5672
    acrLoginServer: infra.outputs.acrLoginServer
    acrUsername: infra.outputs.acrUsername
    acrPassword: infra.outputs.acrPassword
    cpu: '0.5'
    memory: '1Gi'
    minReplicas: 1
    maxReplicas: 1
    volumeStorageName: rabbitmqStorageMount.outputs.storageName
    volumeMountPath: '/var/lib/rabbitmq'
  }
}

// ---- Redis -----------------------------------------------------------
module redis 'modules/container-app-with-volume.bicep' = {
  name: 'app-redis'
  params: {
    name: 'redis'
    location: location
    environmentId: infra.outputs.environmentId
    image: 'redis:7-alpine'
    targetPort: 6379
    ingressType: 'internal'
    transport: 'tcp'
    exposedPort: 6379
    acrLoginServer: infra.outputs.acrLoginServer
    acrUsername: infra.outputs.acrUsername
    acrPassword: infra.outputs.acrPassword
    cpu: '0.25'
    memory: '0.5Gi'
    minReplicas: 1
    maxReplicas: 1
    volumeStorageName: redisStorageMount.outputs.storageName
    volumeMountPath: '/data'
  }
}

// ---- user-auth-mongo -------------------------------------------------
module userAuthMongo 'modules/container-app-with-volume.bicep' = {
  name: 'app-user-auth-mongo'
  params: {
    name: 'user-auth-mongo'
    location: location
    environmentId: infra.outputs.environmentId
    image: 'mongo:7'
    targetPort: 27017
    ingressType: 'internal'
    transport: 'tcp'
    exposedPort: 27017
    envVars: [
      {
        name: 'MONGO_INITDB_ROOT_USERNAME'
        value: 'root'
      }
      {
        name: 'MONGO_INITDB_ROOT_PASSWORD'
        value: 'example'
      }
    ]
    acrLoginServer: infra.outputs.acrLoginServer
    acrUsername: infra.outputs.acrUsername
    acrPassword: infra.outputs.acrPassword
    cpu: '0.5'
    memory: '1Gi'
    minReplicas: 1
    maxReplicas: 1
    volumeStorageName: mongoStorageMount.outputs.storageName
    volumeMountPath: '/data/db'
  }
}

// ---- cars-db ---------------------------------------------------------
module carsDb 'modules/container-app-with-volume.bicep' = {
  name: 'app-cars-db'
  params: {
    name: 'cars-db'
    location: location
    environmentId: infra.outputs.environmentId
    image: 'postgres:16-alpine'
    targetPort: 5432
    ingressType: 'internal'
    transport: 'tcp'
    exposedPort: 5432
    envVars: [
      {
        name: 'POSTGRES_USER'
        value: 'car'
      }
      {
        name: 'POSTGRES_PASSWORD'
        value: 'car'
      }
      {
        name: 'POSTGRES_DB'
        value: 'cars_db'
      }
    ]
    acrLoginServer: infra.outputs.acrLoginServer
    acrUsername: infra.outputs.acrUsername
    acrPassword: infra.outputs.acrPassword
    cpu: '0.5'
    memory: '1Gi'
    minReplicas: 1
    maxReplicas: 1
    volumeStorageName: carsDbStorageMount.outputs.storageName
    volumeMountPath: '/var/lib/postgresql/data'
  }
}

// ---- booking-db ------------------------------------------------------
module bookingDb 'modules/container-app-with-volume.bicep' = {
  name: 'app-booking-db'
  params: {
    name: 'booking-db'
    location: location
    environmentId: infra.outputs.environmentId
    image: 'postgres:16-alpine'
    targetPort: 5432
    ingressType: 'internal'
    transport: 'tcp'
    exposedPort: 5432
    envVars: [
      {
        name: 'POSTGRES_USER'
        value: 'booking'
      }
      {
        name: 'POSTGRES_PASSWORD'
        value: 'booking'
      }
      {
        name: 'POSTGRES_DB'
        value: 'booking'
      }
    ]
    acrLoginServer: infra.outputs.acrLoginServer
    acrUsername: infra.outputs.acrUsername
    acrPassword: infra.outputs.acrPassword
    cpu: '0.5'
    memory: '1Gi'
    minReplicas: 1
    maxReplicas: 1
    volumeStorageName: bookingDbStorageMount.outputs.storageName
    volumeMountPath: '/var/lib/postgresql/data'
  }
}

// =====================================================================
// 4. EIGENE SERVICES (Placeholder-Image, ersetzt durch CI/CD)
// =====================================================================
// HTTP-URLs ohne Port: Container Apps internal HTTP-Ingress routet auf
// Port 80 zum targetPort. Service Discovery via simplem Hostname.
//
// JWT-Keys werden als Container Apps Secrets + Secret-Volume gemountet.
// Pfade matchen das Compose-Setup (/run/secrets/jwt_*.pem), sodass die
// _FILE-Env-Vars unverändert übernommen werden können.

// ---- api-gateway (Spring Boot, public) -------------------------------
module apiGateway 'modules/container-app.bicep' = {
  name: 'app-api-gateway'
  params: {
    name: 'api-gateway'
    location: location
    environmentId: infra.outputs.environmentId
    image: placeholderImage
    targetPort: 8080
    ingressType: 'external'
    transport: 'auto'
    envVars: [
      {
        name: 'SPRING_PROFILES_ACTIVE'
        value: 'docker'
      }
      {
        name: 'HTTP_PORT'
        value: '8080'
      }
      {
        name: 'USER_AUTH_URL'
        value: 'http://user-auth'
      }
      {
        name: 'CAR_SERVICE_URL'
        value: 'http://car-service'
      }
      {
        name: 'BOOKING_SERVICE_URL'
        value: 'http://booking'
      }
      {
        name: 'REDIS_HOST'
        value: 'redis'
      }
      {
        name: 'REDIS_PORT'
        value: '6379'
      }
      {
        name: 'JWT_PUBLIC_KEY_FILE'
        value: '/run/secrets/jwt_public.pem'
      }
      {
        name: 'CORS_ALLOWED_ORIGINS'
        value: 'http://localhost:3000,http://localhost:5173'
      }
    ]
    acrLoginServer: infra.outputs.acrLoginServer
    acrUsername: infra.outputs.acrUsername
    acrPassword: infra.outputs.acrPassword
    cpu: '0.5'
    memory: '1Gi'
    minReplicas: 1
    maxReplicas: 3
    secretFiles: jwtPublicOnly
  }
}

// ---- user-auth (Go, internal) ----------------------------------------
module userAuth 'modules/container-app.bicep' = {
  name: 'app-user-auth'
  params: {
    name: 'user-auth'
    location: location
    environmentId: infra.outputs.environmentId
    image: placeholderImage
    targetPort: 8081
    ingressType: 'internal'
    transport: 'auto'
    envVars: [
      {
        name: 'HTTP_PORT'
        value: '8081'
      }
      {
        name: 'LOG_LEVEL'
        value: 'info'
      }
      {
        name: 'LOG_FORMAT'
        value: 'json'
      }
      {
        name: 'MONGO_URI'
        value: 'mongodb://root:example@user-auth-mongo:27017/?authSource=admin'
      }
      {
        name: 'MONGO_DATABASE'
        value: 'auth'
      }
      {
        name: 'RABBITMQ_URL'
        value: 'amqp://guest:guest@rabbitmq:5672/'
      }
      {
        name: 'RABBITMQ_EXCHANGE'
        value: 'user.events'
      }
      {
        name: 'JWT_PRIVATE_KEY_FILE'
        value: '/run/secrets/jwt_private.pem'
      }
      {
        name: 'JWT_PUBLIC_KEY_FILE'
        value: '/run/secrets/jwt_public.pem'
      }
      {
        name: 'JWT_ISSUER'
        value: 'user-auth'
      }
      {
        name: 'ADMIN_BOOTSTRAP_EMAIL'
        value: 'admin@uni.de'
      }
      {
        name: 'ADMIN_BOOTSTRAP_PASSWORD'
        value: 'ChangeMeChangeMe'
      }
    ]
    acrLoginServer: infra.outputs.acrLoginServer
    acrUsername: infra.outputs.acrUsername
    acrPassword: infra.outputs.acrPassword
    cpu: '0.25'
    memory: '0.5Gi'
    minReplicas: 1
    maxReplicas: 3
    secretFiles: jwtBoth
  }
}

// ---- car-service (Spring Boot, internal) -----------------------------
module carService 'modules/container-app.bicep' = {
  name: 'app-car-service'
  params: {
    name: 'car-service'
    location: location
    environmentId: infra.outputs.environmentId
    image: placeholderImage
    targetPort: 8082
    ingressType: 'internal'
    transport: 'auto'
    envVars: [
      {
        name: 'SPRING_PROFILES_ACTIVE'
        value: 'docker'
      }
      {
        name: 'SPRING_DATASOURCE_URL'
        value: 'jdbc:postgresql://cars-db:5432/cars_db'
      }
      {
        name: 'SPRING_DATASOURCE_USERNAME'
        value: 'car'
      }
      {
        name: 'SPRING_DATASOURCE_PASSWORD'
        value: 'car'
      }
      {
        name: 'JWT_PUBLIC_KEY_FILE'
        value: '/run/secrets/jwt_public.pem'
      }
      {
        name: 'CURRENCY_CONVERTER_GRPC_ADDRESS'
        // currency-converter läuft mit transport=http2 und allowInsecure;
        // internal HTTP/2-Ingress terminiert auf Port 80 zum targetPort 9000.
        value: 'currency-converter:80'
      }
      {
        name: 'CURRENCY_CONVERTER_GRPC_DEADLINE_MS'
        value: '5000'
      }
    ]
    acrLoginServer: infra.outputs.acrLoginServer
    acrUsername: infra.outputs.acrUsername
    acrPassword: infra.outputs.acrPassword
    cpu: '0.5'
    memory: '1Gi'
    minReplicas: 1
    maxReplicas: 3
    secretFiles: jwtPublicOnly
  }
}

// ---- booking (Spring Boot, internal) ---------------------------------
module booking 'modules/container-app.bicep' = {
  name: 'app-booking'
  params: {
    name: 'booking'
    location: location
    environmentId: infra.outputs.environmentId
    image: placeholderImage
    targetPort: 8080
    ingressType: 'internal'
    transport: 'auto'
    envVars: [
      {
        name: 'SPRING_PROFILES_ACTIVE'
        value: 'docker'
      }
      {
        name: 'SPRING_DATASOURCE_URL'
        value: 'jdbc:postgresql://booking-db:5432/booking'
      }
      {
        name: 'SPRING_DATASOURCE_USERNAME'
        value: 'booking'
      }
      {
        name: 'SPRING_DATASOURCE_PASSWORD'
        value: 'booking'
      }
      {
        name: 'SPRING_RABBITMQ_HOST'
        value: 'rabbitmq'
      }
      {
        name: 'SPRING_RABBITMQ_PORT'
        value: '5672'
      }
      {
        name: 'SPRING_RABBITMQ_USERNAME'
        value: 'guest'
      }
      {
        name: 'SPRING_RABBITMQ_PASSWORD'
        value: 'guest'
      }
      {
        name: 'JWT_PUBLIC_KEY_FILE'
        value: '/run/secrets/jwt_public.pem'
      }
      {
        name: 'CURRENCY_CONVERTER_GRPC_ADDRESS'
        value: 'currency-converter:80'
      }
      {
        name: 'CURRENCY_CONVERTER_GRPC_DEADLINE_MS'
        value: '5000'
      }
      {
        name: 'CAR_SERVICE_URL'
        value: 'http://car-service'
      }
    ]
    acrLoginServer: infra.outputs.acrLoginServer
    acrUsername: infra.outputs.acrUsername
    acrPassword: infra.outputs.acrPassword
    cpu: '0.5'
    memory: '1Gi'
    minReplicas: 1
    maxReplicas: 3
    secretFiles: jwtPublicOnly
  }
}

// ---- currency-converter (gRPC, internal, http2) ----------------------
module currencyConverter 'modules/container-app.bicep' = {
  name: 'app-currency-converter'
  params: {
    name: 'currency-converter'
    location: location
    environmentId: infra.outputs.environmentId
    image: placeholderImage
    targetPort: 9000
    ingressType: 'internal'
    transport: 'http2'
    envVars: [
      {
        name: 'ECB_FEED_URL'
        value: 'https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml'
      }
      {
        name: 'ECB_REFRESH_INTERVAL'
        value: '3600'
      }
      {
        name: 'GRPC_PORT'
        value: '9000'
      }
      {
        name: 'LOG_LEVEL'
        value: 'info'
      }
    ]
    acrLoginServer: infra.outputs.acrLoginServer
    acrUsername: infra.outputs.acrUsername
    acrPassword: infra.outputs.acrPassword
    cpu: '0.25'
    memory: '0.5Gi'
    minReplicas: 1
    maxReplicas: 3
  }
}

// ---- notification (internal) -----------------------------------------
module notification 'modules/container-app.bicep' = {
  name: 'app-notification'
  params: {
    name: 'notification'
    location: location
    environmentId: infra.outputs.environmentId
    image: placeholderImage
    targetPort: 8084
    ingressType: 'internal'
    transport: 'auto'
    envVars: [
      {
        name: 'RABBITMQ_URL'
        value: 'amqp://guest:guest@rabbitmq:5672/'
      }
      {
        name: 'HTTP_PORT'
        value: '8084'
      }
      {
        name: 'NOTIFIER_TYPE'
        value: 'mock'
      }
    ]
    acrLoginServer: infra.outputs.acrLoginServer
    acrUsername: infra.outputs.acrUsername
    acrPassword: infra.outputs.acrPassword
    cpu: '0.25'
    memory: '0.5Gi'
    minReplicas: 1
    maxReplicas: 2
  }
}

// ---- frontend (Vite/Nginx, public) -----------------------------------
// Vite ist build-time: VITE_API_URL muss beim Image-Build gesetzt sein,
// nicht zur Laufzeit. CI/CD baut das Image mit der echten api-gateway-URL.
module frontend 'modules/container-app.bicep' = {
  name: 'app-frontend'
  params: {
    name: 'frontend'
    location: location
    environmentId: infra.outputs.environmentId
    image: placeholderImage
    targetPort: 80
    ingressType: 'external'
    transport: 'auto'
    envVars: []
    acrLoginServer: infra.outputs.acrLoginServer
    acrUsername: infra.outputs.acrUsername
    acrPassword: infra.outputs.acrPassword
    cpu: '0.25'
    memory: '0.5Gi'
    minReplicas: 1
    maxReplicas: 3
  }
}

// =====================================================================
// 5. OUTPUTS
// =====================================================================
output acrLoginServer string = infra.outputs.acrLoginServer
output acrName string = infra.outputs.acrName
output frontendUrl string = 'https://${frontend.outputs.fqdn}'
output apiGatewayUrl string = 'https://${apiGateway.outputs.fqdn}'
