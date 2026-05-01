// Parameter-Werte für die Dev-Umgebung.
// uniqueSuffix MUSS angepasst werden — er fließt in die Namen von ACR
// und Storage Account ein, die global eindeutig sein müssen.
// Empfehlung: deine Initialen + 3-stellige Zahl, lowercase, max 8 Zeichen.

using 'main.bicep'

param prefix = 'platform'
param uniqueSuffix = 'lp13'

// JWT-Schlüssel werden aus Env-Vars gelesen — gesetzt vom deploy.sh.
// Default '' erlaubt `az bicep build` ohne gesetzte Vars (für CI/Linting).
param jwtPrivateKey = readEnvironmentVariable('JWT_PRIVATE_KEY', '')
param jwtPublicKey = readEnvironmentVariable('JWT_PUBLIC_KEY', '')
