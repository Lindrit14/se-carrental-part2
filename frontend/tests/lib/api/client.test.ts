import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { configureClient, request } from '@/lib/api/client';
import { tokenStorage } from '@/lib/auth/tokenStorage';
import { ApiError } from '@/lib/api/errors';

describe('api/client', () => {
  beforeEach(() => {
    tokenStorage.clear();
  });
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('attaches the bearer token from storage and parses JSON responses', async () => {
    tokenStorage.write({
      accessToken: 'access-1',
      refreshToken: 'refresh-1',
      accessExpiresAt: '',
      refreshExpiresAt: '',
    });
    const fetchMock = vi.fn().mockResolvedValue(
      new Response(JSON.stringify({ id: '42' }), {
        status: 200,
        headers: { 'content-type': 'application/json' },
      }),
    );
    vi.stubGlobal('fetch', fetchMock);

    const result = await request<{ id: string }>('booking', '/api/v1/cars/42');

    expect(result).toEqual({ id: '42' });
    const [, init] = fetchMock.mock.calls[0]!;
    expect(init.headers.Authorization).toBe('Bearer access-1');
  });

  it('refreshes the token on 401 and retries the request', async () => {
    tokenStorage.write({
      accessToken: 'expired',
      refreshToken: 'refresh-1',
      accessExpiresAt: '',
      refreshExpiresAt: '',
    });

    const responses = [
      new Response('{"code":"invalid_token","message":"expired"}', { status: 401 }),
      new Response(JSON.stringify({ ok: true }), {
        status: 200,
        headers: { 'content-type': 'application/json' },
      }),
    ];
    const fetchMock = vi.fn().mockImplementation(() => Promise.resolve(responses.shift()!));
    vi.stubGlobal('fetch', fetchMock);

    const refresh = vi.fn().mockResolvedValue({
      accessToken: 'fresh',
      refreshToken: 'refresh-2',
      accessExpiresAt: '',
      refreshExpiresAt: '',
    });
    const logout = vi.fn();
    configureClient({ refresh, logout });

    const result = await request<{ ok: boolean }>('booking', '/api/v1/bookings/me');

    expect(result).toEqual({ ok: true });
    expect(refresh).toHaveBeenCalledOnce();
    expect(fetchMock).toHaveBeenCalledTimes(2);
    expect(logout).not.toHaveBeenCalled();
    expect(tokenStorage.read()?.accessToken).toBe('fresh');
  });

  it('throws an ApiError when the response is not ok', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue(
        new Response(JSON.stringify({ code: 'invalid_credentials', message: 'nope' }), {
          status: 401,
          headers: { 'content-type': 'application/json' },
        }),
      ),
    );
    // No tokens → no refresh attempt → error bubbles up.
    await expect(
      request('auth', '/api/v1/auth/login', {
        method: 'POST',
        body: { email: 'a@b', password: 'x' },
        auth: false,
      }),
    ).rejects.toBeInstanceOf(ApiError);
  });
});
