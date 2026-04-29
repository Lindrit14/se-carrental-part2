export interface AuthTokens {
  accessToken: string;
  refreshToken: string;
  accessExpiresAt: string;
  refreshExpiresAt: string;
}

/**
 * Decoded session built from the access token JWT. Pure value, derived
 * client-side — never sent to the server.
 */
export interface AuthSession {
  userId: string;
  email?: string;
  roles: string[];
  expiresAt: number; // unix seconds
}
