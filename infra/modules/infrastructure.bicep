// =====================================================================
// Basis-Infrastruktur für die Container-Apps-Plattform
//
// Erstellt:
//   - Log Analytics Workspace (Logs/Metrics für die Environment)
//   - Azure Container Registry (Basic SKU, Admin-User aktiviert)
//   - Storage Account + 5 File Shares (für persistente DB-Volumes)
//   - Container Apps Managed Environment (verbunden mit Log Analytics)
// =====================================================================

@description('Basisname für alle Resourcen, z.B. "platform"')
param prefix string

@description('Eindeutiges Suffix (lowercase, max ~8 Zeichen) für globale Namen')
param uniqueSuffix string

@description('Azure-Region für alle Resourcen')
param location string

// ---- Naming -----------------------------------------------------------
// ACR und Storage Account dürfen keine Bindestriche enthalten und müssen
// global eindeutig sein. Daher Suffix als zusätzliche Eindeutigkeit.
var acrName = toLower('${prefix}acr${uniqueSuffix}')
var storageAccountName = toLower('${prefix}st${uniqueSuffix}')
var logAnalyticsName = '${prefix}-logs-${uniqueSuffix}'
var environmentName = '${prefix}-env-${uniqueSuffix}'

// ---- Log Analytics Workspace -----------------------------------------
resource logAnalytics 'Microsoft.OperationalInsights/workspaces@2023-09-01' = {
  name: logAnalyticsName
  location: location
  properties: {
    sku: {
      name: 'PerGB2018'
    }
    retentionInDays: 30
  }
}

// ---- Azure Container Registry ----------------------------------------
// Basic SKU reicht für Dev. Admin-User ist nötig damit Container Apps
// per Username/Password die Images ziehen können.
resource acr 'Microsoft.ContainerRegistry/registries@2023-11-01-preview' = {
  name: acrName
  location: location
  sku: {
    name: 'Basic'
  }
  properties: {
    adminUserEnabled: true
  }
}

// ---- Storage Account für Azure Files ---------------------------------
// largeFileSharesState=Enabled ist Voraussetzung für Quotas > 5 TiB.
resource storageAccount 'Microsoft.Storage/storageAccounts@2023-05-01' = {
  name: storageAccountName
  location: location
  sku: {
    name: 'Standard_LRS'
  }
  kind: 'StorageV2'
  properties: {
    largeFileSharesState: 'Enabled'
    minimumTlsVersion: 'TLS1_2'
    allowBlobPublicAccess: false
  }
}

// File Service ist immer "default" und wird implizit durch das Storage
// Account erstellt. Wir referenzieren ihn als parent für die Shares.
resource fileService 'Microsoft.Storage/storageAccounts/fileServices@2023-05-01' existing = {
  parent: storageAccount
  name: 'default'
}

// ---- 5 File Shares (je ein Volume pro persistentem DB-Service) -------
resource rabbitmqShare 'Microsoft.Storage/storageAccounts/fileServices/shares@2023-05-01' = {
  parent: fileService
  name: 'rabbitmq-data'
  properties: {
    shareQuota: 10
    enabledProtocols: 'SMB'
  }
}

resource redisShare 'Microsoft.Storage/storageAccounts/fileServices/shares@2023-05-01' = {
  parent: fileService
  name: 'redis-data'
  properties: {
    shareQuota: 5
    enabledProtocols: 'SMB'
  }
}

resource mongoShare 'Microsoft.Storage/storageAccounts/fileServices/shares@2023-05-01' = {
  parent: fileService
  name: 'mongo-data'
  properties: {
    shareQuota: 10
    enabledProtocols: 'SMB'
  }
}

resource carsDbShare 'Microsoft.Storage/storageAccounts/fileServices/shares@2023-05-01' = {
  parent: fileService
  name: 'cars-db-data'
  properties: {
    shareQuota: 10
    enabledProtocols: 'SMB'
  }
}

resource bookingDbShare 'Microsoft.Storage/storageAccounts/fileServices/shares@2023-05-01' = {
  parent: fileService
  name: 'booking-db-data'
  properties: {
    shareQuota: 10
    enabledProtocols: 'SMB'
  }
}

// ---- Container Apps Managed Environment ------------------------------
// Sammelt alle Apps, hostet das interne DNS und schickt Logs/Metrics
// an Log Analytics.
resource environment 'Microsoft.App/managedEnvironments@2024-03-01' = {
  name: environmentName
  location: location
  properties: {
    appLogsConfiguration: {
      destination: 'log-analytics'
      logAnalyticsConfiguration: {
        customerId: logAnalytics.properties.customerId
        // listKeys() inline aufrufen, nicht in Variable speichern,
        // damit der Key nicht in der Deployment-History landet.
        sharedKey: logAnalytics.listKeys().primarySharedKey
      }
    }
  }
}

// ---- Outputs ---------------------------------------------------------
output acrLoginServer string = acr.properties.loginServer
output acrName string = acr.name
output acrUsername string = acr.listCredentials().username

@secure()
output acrPassword string = acr.listCredentials().passwords[0].value

output environmentId string = environment.id
output environmentName string = environment.name

// Default-Domain wird für Service-Discovery zwischen Container Apps gebraucht:
//   - Internal-FQDN:  <app>.internal.<defaultDomain>
//   - External-FQDN:  <app>.<defaultDomain>
// Spring Cloud Gateway o.ä. brauchen den vollen FQDN, weil der interne LB
// anhand des Host-Headers routet — Short-Hostname (http://user-auth) bringt 404.
output environmentDefaultDomain string = environment.properties.defaultDomain

output storageAccountName string = storageAccount.name

@secure()
output storageAccountKey string = storageAccount.listKeys().keys[0].value
