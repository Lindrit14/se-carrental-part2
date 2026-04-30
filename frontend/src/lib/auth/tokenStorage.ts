import type { AuthTokens } from '@/domain/auth';

const KEY = 'auth.tokens';

/**
 * Typed wrapper around localStorage for our auth tokens. localStorage is
 * exposed to XSS; production would use
 * HttpOnly cookies + CSRF (see README).
 */
export const tokenStorage = {
  read(): AuthTokens | null {
    try {
      const raw = localStorage.getItem(KEY);
      if (!raw) return null;
      return JSON.parse(raw) as AuthTokens;
    } catch {
      return null;
    }
  },
  write(tokens: AuthTokens): void {
    localStorage.setItem(KEY, JSON.stringify(tokens));
  },
  clear(): void {
    localStorage.removeItem(KEY);
  },
};
