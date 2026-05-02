// =====================================================================
// Wiederverwendbares Container-App-Modul MIT Azure Files Volume Mount
//
// Identisch zu container-app.bicep, ergänzt um:
//   - volumes-Block auf Template-Ebene (verweist auf Storage-Mount in Env)
//   - volumeMounts-Block auf Container-Ebene (mountPath im Container)
//
// Wird für die persistenten DBs verwendet:
//   rabbitmq, redis, user-auth-mongo, cars-db, booking-db
// =====================================================================

@description('Container App Name (= Service-Name aus docker-compose)')
param name string

@description('Azure Region')
param location string

@description('Resource-ID der Container Apps Environment')
param environmentId string

@description('Voll qualifizierter Image-Name')
param image string

@description('Container-Port')
param targetPort int

@description('Ingress-Typ')
@allowed([
  'external'
  'internal'
  'none'
])
param ingressType string = 'internal'

@description('Transport-Protokoll')
@allowed([
  'auto'
  'http'
  'http2'
  'tcp'
])
param transport string = 'auto'

@description('Externer Port für TCP-Ingress (0 = nicht setzen)')
param exposedPort int = 0

@description('Environment-Variablen')
param envVars array = []

@description('ACR Login Server')
param acrLoginServer string

@description('ACR Username')
param acrUsername string

@description('ACR Password')
@secure()
param acrPassword string

@description('CPU als String')
param cpu string = '0.25'

@description('Memory')
param memory string = '0.5Gi'

@description('Mindest-Replicas (DBs sollten nicht auf 0 skalieren!)')
param minReplicas int = 1

@description('Maximal-Replicas (DBs typischerweise 1, da kein Sharding)')
param maxReplicas int = 1

// ---- Volume-spezifische Parameter ------------------------------------

@description('Logischer Name des Storage-Mounts in der Environment (siehe env-storage.bicep)')
param volumeStorageName string

@description('Mount-Pfad im Container, z.B. /var/lib/postgresql/data')
param volumeMountPath string

@description('Interner Bicep-Name des Volumes (egal welcher, muss konsistent sein)')
param volumeName string = 'data-volume'

@description('System-Assigned Managed Identity aktivieren? Wird für Azure-Plane-Auth (z.B. ACS Email) gebraucht.')
param assignSystemIdentity bool = false

// ---- Ingress-Block ---------------------------------------------------
var ingressBase = {
  external: ingressType == 'external'
  targetPort: targetPort
  transport: transport
  allowInsecure: transport == 'http2'
}
var ingressWithPort = exposedPort > 0 && transport == 'tcp' ? union(ingressBase, {
  exposedPort: exposedPort
}) : ingressBase
var ingress = ingressType == 'none' ? null : ingressWithPort

// ---- Container App mit Volume ----------------------------------------
resource app 'Microsoft.App/containerApps@2024-03-01' = {
  name: name
  location: location
  identity: {
    type: assignSystemIdentity ? 'SystemAssigned' : 'None'
  }
  properties: {
    managedEnvironmentId: environmentId
    configuration: {
      registries: [
        {
          server: acrLoginServer
          username: acrUsername
          passwordSecretRef: 'registry-password'
        }
      ]
      secrets: [
        {
          name: 'registry-password'
          value: acrPassword
        }
      ]
      ingress: ingress
    }
    template: {
      containers: [
        {
          name: name
          image: image
          resources: {
            cpu: json(cpu)
            memory: memory
          }
          env: envVars
          volumeMounts: [
            {
              volumeName: volumeName
              mountPath: volumeMountPath
            }
          ]
        }
      ]
      // Volume verweist auf den vorher in der Environment registrierten
      // Storage-Mount (azureFile mit AccountKey).
      volumes: [
        {
          name: volumeName
          storageType: 'AzureFile'
          storageName: volumeStorageName
        }
      ]
      scale: {
        minReplicas: minReplicas
        maxReplicas: maxReplicas
      }
    }
  }
}

// ---- Outputs ---------------------------------------------------------
output appName string = app.name
output fqdn string = ingressType == 'external' ? app.properties.configuration.ingress.fqdn : ''
// principalId der System-Identity — leer falls nicht aktiviert. Wird für
// Role-Assignments gegen ACS, Key Vault, etc. gebraucht.
output principalId string = assignSystemIdentity ? app.identity.principalId : ''
