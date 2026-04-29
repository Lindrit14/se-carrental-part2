import type { ReactNode } from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { useAuth } from '@/features/auth/useAuth';

export function AdminRoute({ children }: { children: ReactNode }) {
  const { session, initialized } = useAuth();
  const location = useLocation();

  if (!initialized) return null;
  if (!session) {
    return <Navigate to="/login" replace state={{ from: location.pathname }} />;
  }
  if (!session.roles.includes('admin')) {
    return <Navigate to="/" replace />;
  }
  return <>{children}</>;
}
