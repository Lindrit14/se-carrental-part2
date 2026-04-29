// Azure Container Registry. Deploy once per environment.
@description('Name of the registry — must be globally unique, lowercase, 5-50 chars.')
param registryName string

@description('Azure region.')
param location string = resourceGroup().location

@allowed(['Basic', 'Standard', 'Premium'])
param sku string = 'Basic'

resource acr 'Microsoft.ContainerRegistry/registries@2023-11-01-preview' = {
  name: registryName
  location: location
  sku: { name: sku }
  properties: {
    adminUserEnabled: true
  }
}

output loginServer string = acr.properties.loginServer
output registryId  string = acr.id
