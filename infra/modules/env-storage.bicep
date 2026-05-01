// =====================================================================
// Storage-Mount-Definition in der Container Apps Environment
//
// Registriert einen Azure File Share als logischen Mount, den einzelne
// Container Apps anschließend per "storageName" referenzieren können.
// =====================================================================

@description('Name der bereits existierenden Container Apps Environment')
param environmentName string

@description('Logischer Name des Mounts, später in der App referenziert')
param storageMountName string

@description('Name des Storage Accounts, der den File Share hostet')
param storageAccountName string

@description('Storage Account Key (wird als Secret im Environment gespeichert)')
@secure()
param storageAccountKey string

@description('Name des File Shares im Storage Account')
param fileShareName string

@description('Mount-Modus: ReadOnly oder ReadWrite')
@allowed([
  'ReadOnly'
  'ReadWrite'
])
param accessMode string = 'ReadWrite'

// ---- Existing Environment als parent ---------------------------------
resource environment 'Microsoft.App/managedEnvironments@2024-03-01' existing = {
  name: environmentName
}

// ---- Storage Mount ---------------------------------------------------
resource storageMount 'Microsoft.App/managedEnvironments/storages@2024-03-01' = {
  parent: environment
  name: storageMountName
  properties: {
    azureFile: {
      accountName: storageAccountName
      accountKey: storageAccountKey
      shareName: fileShareName
      accessMode: accessMode
    }
  }
}

output storageName string = storageMount.name
