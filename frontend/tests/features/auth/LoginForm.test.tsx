import { describe, expect, it, vi } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router-dom';
import { LoginForm } from '@/features/auth/LoginForm';
import { AuthContext } from '@/features/auth/AuthProvider';

function renderWithAuth(login: (email: string, password: string) => Promise<void>) {
  return render(
    <MemoryRouter>
      <AuthContext.Provider
        value={{
          session: null,
          initialized: true,
          login,
          register: vi.fn(),
          logout: vi.fn().mockResolvedValue(undefined),
        }}
      >
        <LoginForm />
      </AuthContext.Provider>
    </MemoryRouter>,
  );
}

describe('LoginForm', () => {
  it('shows validation errors on empty submit', async () => {
    const login = vi.fn();
    renderWithAuth(login);
    await userEvent.click(screen.getByRole('button', { name: /sign in/i }));
    expect(await screen.findByText(/valid email/i)).toBeInTheDocument();
    expect(login).not.toHaveBeenCalled();
  });

  it('calls login with form values on submit', async () => {
    const login = vi.fn().mockResolvedValue(undefined);
    renderWithAuth(login);

    await userEvent.type(screen.getByLabelText(/email/i), 'alice@example.com');
    await userEvent.type(screen.getByLabelText(/password/i), 'hunter2hunter2');
    await userEvent.click(screen.getByRole('button', { name: /sign in/i }));

    await waitFor(() => {
      expect(login).toHaveBeenCalledWith('alice@example.com', 'hunter2hunter2');
    });
  });
});
