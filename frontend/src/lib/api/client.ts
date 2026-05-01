import type { AuthTokens } from '@/domain/auth';
import { tokenStorage } from '@/lib/auth/tokenStorage';
import { ApiError, type ApiErrorBody } from './errors';

export type Service = 'auth' | 'cars' | 'bookings';

// Single gateway URL — all routing handled server-side by the API gateway.
// Service param kept for caller readability; resolves to the same base URL.
const API_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080';

const BASE_URLS: Record<Service, string> = {
  auth: API_URL,
  cars: API_URL,
  bookings: API_URL,
};

interface RequestOptions {
  method?: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';
  body?: unknown;
  auth?: boolean; // include Authorization header (default: true)
  signal?: AbortSignal;
}

type RefreshFn = (refreshToken: string) => Promise<AuthTokens>;
type LogoutFn = () => void;

/* The auth feature wires these in at startup; keeps client.ts framework-free. */
let refreshFn: RefreshFn | null = null;
let logoutFn: LogoutFn | null = null;

export function configureClient(opts: { refresh: RefreshFn; logout: LogoutFn }): void {
  refreshFn = opts.refresh;
  logoutFn = opts.logout;
}

export async function request<T>(
  service: Service,
  path: string,
  opts: RequestOptions = {},
): Promise<T> {
  const { method = 'GET', body, auth = true, signal } = opts;
  const url = `${BASE_URLS[service]}${path}`;

  const exec = async (token?: string): Promise<Response> => {
    const headers: Record<string, string> = { Accept: 'application/json' };
    if (body !== undefined) headers['Content-Type'] = 'application/json';
    if (auth && token) headers['Authorization'] = `Bearer ${token}`;
    return fetch(url, {
      method,
      headers,
      body: body !== undefined ? JSON.stringify(body) : undefined,
      signal,
    });
  };

  const tokens = auth ? tokenStorage.read() : null;
  let response = await exec(tokens?.accessToken);

  // Try refresh once on 401 if we had a token.
  if (auth && response.status === 401 && tokens?.refreshToken && refreshFn) {
    try {
      const fresh = await refreshFn(tokens.refreshToken);
      tokenStorage.write(fresh);
      response = await exec(fresh.accessToken);
    } catch {
      tokenStorage.clear();
      logoutFn?.();
      throw await toApiError(response);
    }
  }

  if (!response.ok) {
    throw await toApiError(response);
  }

  if (response.status === 204) return undefined as T;
  return (await response.json()) as T;
}

async function toApiError(response: Response): Promise<ApiError> {
  let body: ApiErrorBody = { code: 'unknown', message: response.statusText };
  try {
    body = (await response.json()) as ApiErrorBody;
  } catch {
    /* ignore */
  }
  return new ApiError(response.status, body);
}
