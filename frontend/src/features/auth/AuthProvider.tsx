import { createContext, useCallback, useEffect, useMemo, useState, type ReactNode } from 'react';
import type { AuthSession, AuthTokens } from '@/domain/auth';
import { authApi } from '@/lib/api/authApi';
import { configureClient } from '@/lib/api/client';
import { tokenStorage } from '@/lib/auth/tokenStorage';

interface AuthContextValue {
  session: AuthSession | null;
  initialized: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
}

export const AuthContext = createContext<AuthContextValue | null>(null);

/* ---------- helpers ---------- */

interface JwtPayload {
  sub: string;
  exp: number;
  email?: string;
  roles?: string[];
}

function decodeJwt(token: string): JwtPayload | null {
  try {
    const part = token.split('.')[1];
    if (!part) return null;
    const padded = part.replace(/-/g, '+').replace(/_/g, '/');
    return JSON.parse(atob(padded)) as JwtPayload;
  } catch {
    return null;
  }
}

function sessionFromTokens(tokens: AuthTokens): AuthSession | null {
  const claims = decodeJwt(tokens.accessToken);
  if (!claims) return null;
  return {
    userId: claims.sub,
    email: claims.email,
    roles: claims.roles ?? [],
    expiresAt: claims.exp,
  };
}

/* ---------- provider ---------- */

export function AuthProvider({ children }: { children: ReactNode }) {
  const [session, setSession] = useState<AuthSession | null>(null);
  const [initialized, setInitialized] = useState(false);

  // Wire the API client to our refresh + logout once on mount.
  useEffect(() => {
    configureClient({
      refresh: async (refreshToken) => authApi.refresh(refreshToken),
      logout: () => {
        tokenStorage.clear();
        setSession(null);
      },
    });
    const stored = tokenStorage.read();
    if (stored) setSession(sessionFromTokens(stored));
    setInitialized(true);
  }, []);

  const login = useCallback(async (email: string, password: string) => {
    const tokens = await authApi.login(email, password);
    tokenStorage.write(tokens);
    setSession(sessionFromTokens(tokens));
  }, []);

  const register = useCallback(async (email: string, password: string) => {
    await authApi.register(email, password);
    // After register we log the user in automatically.
    await login(email, password);
  }, [login]);

  const logout = useCallback(async () => {
    const tokens = tokenStorage.read();
    if (tokens) {
      try {
        await authApi.logout(tokens.refreshToken);
      } catch {
        // ignore — the local state is what matters for the UI
      }
    }
    tokenStorage.clear();
    setSession(null);
  }, []);

  const value = useMemo<AuthContextValue>(
    () => ({ session, initialized, login, register, logout }),
    [session, initialized, login, register, logout],
  );

  return (
    <AuthContext.Provider value={value}>
      {initialized ? children : null}
    </AuthContext.Provider>
  );
}
