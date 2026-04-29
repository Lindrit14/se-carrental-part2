import type { ReactNode } from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { useAuth } from '@/features/auth/useAuth';

export function ProtectedRoute({ children }: { children: ReactNode }) {
  const { session, initialized } = useAuth();
  const location = useLocation();

  if (!initialized) {
    // First-tick: tokenStorage hasn't been read yet. Render nothing to avoid a flash.
    return null;
  }
  if (!session) {
    const returnTo = location.pathname + location.search;
    return <Navigate to={`/login?returnTo=${encodeURIComponent(returnTo)}`} replace />;
  }
  return <>{children}</>;
}
