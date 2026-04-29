// Azure Container Apps deployment for the user-auth service.
// Uses an existing Container Apps Environment + ACR; this template
// deploys the auth-service Container App and a self-hosted RabbitMQ
// Container App. MongoDB is expected to be Azure Cosmos DB for MongoDB
// (recommended for prod) — pass MONGO_URI via secret.

@description('Container Apps Environment resource id (existing).')
param environmentId string

@description('ACR login server, e.g. myreg.azurecr.io')
param registryServer string

@description('ACR admin username.')
param registryUsername string

@description('ACR admin password (set via secure parameter).')
@secure()
param registryPassword string

@description('Image tag (e.g. git SHA) for the auth service.')
param imageTag string = 'latest'

@description('Mongo connection URI.')
@secure()
param mongoUri string

@description('JWT private key (PEM).')
@secure()
param jwtPrivateKey string

@description('JWT public key (PEM).')
@secure()
param jwtPublicKey string

param location string = resourceGroup().location

// --- RabbitMQ (self-hosted, internal-only) ----------------------------
resource rabbitmq 'Microsoft.App/containerApps@2024-03-01' = {
  name: 'rabbitmq'
  location: location
  properties: {
    managedEnvironmentId: environmentId
    configuration: {
      ingress: {
        external: false
        targetPort: 5672
        transport: 'tcp'
        exposedPort: 5672
      }
    }
    template: {
      containers: [
        {
          name: 'rabbitmq'
          image: 'rabbitmq:3.13-management-alpine'
          resources: {
            cpu: json('0.5')
            memory: '1Gi'
          }
          probes: [
            {
              type: 'Liveness'
              tcpSocket: { port: 5672 }
              initialDelaySeconds: 30
              periodSeconds: 30
            }
          ]
        }
      ]
      scale: { minReplicas: 1, maxReplicas: 1 }
    }
  }
}

// --- auth-service ------------------------------------------------------
resource auth 'Microsoft.App/containerApps@2024-03-01' = {
  name: 'auth-service'
  location: location
  properties: {
    managedEnvironmentId: environmentId
    configuration: {
      ingress: {
        external: true
        targetPort: 8080
        transport: 'http'
        allowInsecure: false
      }
      registries: [
        {
          server: registryServer
          username: registryUsername
          passwordSecretRef: 'registry-password'
        }
      ]
      secrets: [
        { name: 'registry-password', value: registryPassword }
        { name: 'mongo-uri',         value: mongoUri }
        { name: 'jwt-private-key',   value: jwtPrivateKey }
        { name: 'jwt-public-key',    value: jwtPublicKey }
      ]
    }
    template: {
      containers: [
        {
          name: 'auth'
          image: '${registryServer}/user-auth:${imageTag}'
          resources: {
            cpu: json('0.5')
            memory: '1Gi'
          }
          env: [
            { name: 'HTTP_PORT',     value: '8080' }
            { name: 'LOG_LEVEL',     value: 'info' }
            { name: 'LOG_FORMAT',    value: 'json' }
            { name: 'MONGO_URI',     secretRef: 'mongo-uri' }
            { name: 'MONGO_DATABASE', value: 'auth' }
            { name: 'RABBITMQ_URL',  value: 'amqp://guest:guest@rabbitmq:5672/' }
            { name: 'RABBITMQ_EXCHANGE', value: 'user.events' }
            { name: 'JWT_PRIVATE_KEY', secretRef: 'jwt-private-key' }
            { name: 'JWT_PUBLIC_KEY',  secretRef: 'jwt-public-key' }
            { name: 'JWT_ISSUER',    value: 'user-auth' }
          ]
          probes: [
            {
              type: 'Liveness'
              httpGet: { path: '/healthz', port: 8080 }
              initialDelaySeconds: 5
              periodSeconds: 10
            }
            {
              type: 'Readiness'
              httpGet: { path: '/readyz', port: 8080 }
              initialDelaySeconds: 5
              periodSeconds: 10
            }
          ]
        }
      ]
      scale: {
        minReplicas: 1
        maxReplicas: 5
        rules: [
          {
            name: 'http-concurrency'
            http: { metadata: { concurrentRequests: '50' } }
          }
        ]
      }
    }
  }
  dependsOn: [ rabbitmq ]
}

output authFqdn string = auth.properties.configuration.ingress.fqdn
