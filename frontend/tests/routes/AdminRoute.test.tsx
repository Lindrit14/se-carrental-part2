import { describe, expect, it, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { AuthContext } from '@/features/auth/AuthProvider';
import { AdminRoute } from '@/routes/AdminRoute';
import type { AuthSession } from '@/domain/auth';

function renderAt(initialPath: string, session: AuthSession | null) {
  return render(
    <MemoryRouter initialEntries={[initialPath]}>
      <AuthContext.Provider
        value={{
          session,
          initialized: true,
          login: vi.fn(),
          register: vi.fn(),
          logout: vi.fn().mockResolvedValue(undefined),
        }}
      >
        <Routes>
          <Route path="/" element={<div>HOME</div>} />
          <Route path="/login" element={<div>LOGIN</div>} />
          <Route
            path="/admin"
            element={
              <AdminRoute>
                <div>ADMIN_PANEL</div>
              </AdminRoute>
            }
          />
        </Routes>
      </AuthContext.Provider>
    </MemoryRouter>,
  );
}

describe('AdminRoute', () => {
  it('redirects anonymous visitors to /login', () => {
    renderAt('/admin', null);
    expect(screen.getByText('LOGIN')).toBeInTheDocument();
  });

  it('redirects non-admin users to /', () => {
    renderAt('/admin', {
      userId: 'u1',
      email: 'u@example.com',
      roles: ['user'],
      expiresAt: 9999999999,
    });
    expect(screen.getByText('HOME')).toBeInTheDocument();
  });

  it('lets admins through', () => {
    renderAt('/admin', {
      userId: 'u2',
      email: 'a@example.com',
      roles: ['user', 'admin'],
      expiresAt: 9999999999,
    });
    expect(screen.getByText('ADMIN_PANEL')).toBeInTheDocument();
  });
});
