import type { AuthTokens } from '@/domain/auth';
import type { UserProfile, Role } from '@/domain/user';
import { request } from './client';

/* ---------- Wire DTOs (snake_case from the Go backend) ---------- */

interface TokenResponseDto {
  access_token: string;
  access_expires_at: string;
  refresh_token: string;
  refresh_expires_at: string;
  token_type: 'Bearer';
}

interface RegisterResponseDto {
  user_id: string;
  email: string;
}

interface ProfileDto {
  id: string;
  email: string;
  roles: string[];
  verified: boolean;
  created_at: string;
}

interface AdminUserDto {
  id: string;
  email: string;
  roles: string[];
  verified: boolean;
  created_at: string;
  updated_at: string;
}

interface AdminUserListDto {
  items: AdminUserDto[];
  total: number;
  limit: number;
  offset: number;
}

export interface AdminUser {
  id: string;
  email: string;
  roles: Role[];
  verified: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface AdminUserListResult {
  items: AdminUser[];
  total: number;
  limit: number;
  offset: number;
}

const toAdminUser = (dto: AdminUserDto): AdminUser => ({
  id: dto.id,
  email: dto.email,
  roles: dto.roles as Role[],
  verified: dto.verified,
  createdAt: dto.created_at,
  updatedAt: dto.updated_at,
});

const toTokens = (dto: TokenResponseDto): AuthTokens => ({
  accessToken: dto.access_token,
  accessExpiresAt: dto.access_expires_at,
  refreshToken: dto.refresh_token,
  refreshExpiresAt: dto.refresh_expires_at,
});

const toProfile = (dto: ProfileDto): UserProfile => ({
  id: dto.id,
  email: dto.email,
  roles: dto.roles as Role[],
  verified: dto.verified,
  createdAt: dto.created_at,
});

/* ---------- Endpoints ---------- */

export const authApi = {
  async register(email: string, password: string): Promise<{ userId: string; email: string }> {
    const dto = await request<RegisterResponseDto>('auth', '/api/v1/auth/register', {
      method: 'POST',
      body: { email, password },
      auth: false,
    });
    return { userId: dto.user_id, email: dto.email };
  },

  async login(email: string, password: string): Promise<AuthTokens> {
    const dto = await request<TokenResponseDto>('auth', '/api/v1/auth/login', {
      method: 'POST',
      body: { email, password },
      auth: false,
    });
    return toTokens(dto);
  },

  async refresh(refreshToken: string): Promise<AuthTokens> {
    const dto = await request<TokenResponseDto>('auth', '/api/v1/auth/refresh', {
      method: 'POST',
      body: { refresh_token: refreshToken },
      auth: false,
    });
    return toTokens(dto);
  },

  async logout(refreshToken: string): Promise<void> {
    await request<void>('auth', '/api/v1/auth/logout', {
      method: 'POST',
      body: { refresh_token: refreshToken },
    });
  },

  async requestPasswordReset(email: string): Promise<void> {
    await request<void>('auth', '/api/v1/auth/password/reset', {
      method: 'POST',
      body: { email },
      auth: false,
    });
  },

  async confirmPasswordReset(resetToken: string, newPassword: string): Promise<void> {
    await request<void>('auth', '/api/v1/auth/password/confirm', {
      method: 'POST',
      body: { reset_token: resetToken, new_password: newPassword },
      auth: false,
    });
  },

  async me(): Promise<UserProfile> {
    const dto = await request<ProfileDto>('auth', '/api/v1/users/me');
    return toProfile(dto);
  },

  async updateMe(input: { email?: string }): Promise<UserProfile> {
    const dto = await request<ProfileDto>('auth', '/api/v1/users/me', {
      method: 'PATCH',
      body: input,
    });
    return toProfile(dto);
  },

  async deleteMe(): Promise<void> {
    await request<void>('auth', '/api/v1/users/me', { method: 'DELETE' });
  },

  /* ---------- Admin ---------- */

  async adminListUsers(opts?: { limit?: number; offset?: number }): Promise<AdminUserListResult> {
    const params = new URLSearchParams();
    if (opts?.limit != null) params.set('limit', String(opts.limit));
    if (opts?.offset != null) params.set('offset', String(opts.offset));
    const qs = params.toString();
    const path = qs ? `/api/v1/admin/users?${qs}` : '/api/v1/admin/users';
    const dto = await request<AdminUserListDto>('auth', path);
    return {
      items: dto.items.map(toAdminUser),
      total: dto.total,
      limit: dto.limit,
      offset: dto.offset,
    };
  },

  async adminGetUser(id: string): Promise<AdminUser> {
    const dto = await request<AdminUserDto>('auth', `/api/v1/admin/users/${id}`);
    return toAdminUser(dto);
  },

  async adminSetRoles(id: string, roles: Role[]): Promise<AdminUser> {
    const dto = await request<AdminUserDto>('auth', `/api/v1/admin/users/${id}/roles`, {
      method: 'PATCH',
      body: { roles },
    });
    return toAdminUser(dto);
  },

  async adminDeleteUser(id: string): Promise<void> {
    await request<void>('auth', `/api/v1/admin/users/${id}`, { method: 'DELETE' });
  },
};
