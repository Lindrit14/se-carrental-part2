// =====================================================================
// Wiederverwendbares Container-App-Modul OHNE Volume
//
// Deckt alle drei Ingress-Typen ab:
//   - external: öffentlich erreichbar (z.B. api-gateway, frontend)
//   - internal: nur innerhalb der Environment erreichbar (Defaultfall)
//   - none:     keine Ingress (z.B. reine Worker)
// Unterstützt http / http2 (gRPC) / tcp / auto Transport.
// Optional: Secret-Files via secretFiles[] mounten unter /run/secrets/.
// =====================================================================

@description('Container App Name (= Service-Name aus docker-compose, wichtig für Service Discovery)')
param name string

@description('Azure Region')
param location string

@description('Resource-ID der Container Apps Environment')
param environmentId string

@description('Voll qualifizierter Image-Name, z.B. acr.azurecr.io/foo:tag')
param image string

@description('Container-Port, auf dem die App lauscht')
param targetPort int

@description('Ingress-Typ')
@allowed([
  'external'
  'internal'
  'none'
])
param ingressType string = 'internal'

@description('Transport-Protokoll für das Ingress')
@allowed([
  'auto'
  'http'
  'http2'
  'tcp'
])
param transport string = 'auto'

@description('Externer Port für TCP-Ingress (optional, 0 = nicht setzen)')
param exposedPort int = 0

@description('Environment-Variablen für den Container, Form: [{ name, value }] oder [{ name, secretRef }]')
param envVars array = []

@description('ACR Login Server, z.B. myacr.azurecr.io')
param acrLoginServer string

@description('ACR Username')
param acrUsername string

@description('ACR Password (wird als Secret hinterlegt)')
@secure()
param acrPassword string

@description('CPU als String, z.B. "0.25" oder "0.5" — wird per json() in Number gewandelt')
param cpu string = '0.25'

@description('Memory, z.B. "0.5Gi" oder "1Gi"')
param memory string = '0.5Gi'

@description('Mindest-Replicas (0 = Scale-to-Zero erlaubt)')
param minReplicas int = 0

@description('Maximal-Replicas')
param maxReplicas int = 1

@description('Secret-Files die als Volume unter /run/secrets/ gemountet werden. Form: { "<secretName>": { fileName, value } }. Leeres Objekt = kein Mount.')
@secure()
param secretFiles object = {}

@description('System-Assigned Managed Identity aktivieren? Wird für Azure-Plane-Auth (z.B. ACS Email) gebraucht.')
param assignSystemIdentity bool = false

// ---- Ingress-Block je nach Typ aufbauen ------------------------------
// 'none' → kein Ingress; sonst external/internal mit ggf. exposedPort.
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

// ---- Secrets + Volumes für Secret-Files ------------------------------
// secretFiles → 1 Container Apps Secret pro File + 1 Secret-Volume das
// alle Files unter /run/secrets/<fileName> mountet (kompatibel mit den
// _FILE-Env-Vars die der lokale Compose-Stack nutzt).
//
// Form: { '<secretName>': { fileName: '<x>.pem', value: '<PEM-Inhalt>' } }
// items() macht aus dem Object [{ key, value }] für [for]-Iteration.
var hasSecretFiles = !empty(secretFiles)

var fileSecrets = [for f in items(secretFiles): {
  name: f.key
  value: f.value.value
}]

var allSecrets = concat([
  {
    name: 'registry-password'
    value: acrPassword
  }
], fileSecrets)

var secretVolumeRefs = [for f in items(secretFiles): {
  secretRef: f.key
  path: f.value.fileName
}]

// ---- Container App ---------------------------------------------------
resource app 'Microsoft.App/containerApps@2024-03-01' = {
  name: name
  location: location
  identity: {
    type: assignSystemIdentity ? 'SystemAssigned' : 'None'
  }
  properties: {
    managedEnvironmentId: environmentId
    configuration: {
      // Registry-Credentials immer eintragen, damit CI/CD Images aus
      // der ACR ziehen kann sobald das Placeholder-Image ersetzt wird.
      registries: [
        {
          server: acrLoginServer
          username: acrUsername
          passwordSecretRef: 'registry-password'
        }
      ]
      secrets: allSecrets
      ingress: ingress
    }
    template: {
      containers: [
        {
          name: name
          image: image
          resources: {
            // Bicep-Quirk: cpu muss eine Number sein, daher json()-Wrap.
            cpu: json(cpu)
            memory: memory
          }
          env: envVars
          // Secret-Volume unter /run/secrets/ mounten — nur wenn welche
          // konfiguriert sind. Ohne Files bleibt volumeMounts leer.
          volumeMounts: hasSecretFiles ? [
            {
              volumeName: 'secrets-volume'
              mountPath: '/run/secrets'
            }
          ] : []
        }
      ]
      volumes: hasSecretFiles ? [
        {
          name: 'secrets-volume'
          storageType: 'Secret'
          secrets: secretVolumeRefs
        }
      ] : []
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
