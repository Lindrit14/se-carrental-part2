import { Navigate, Outlet, Route, Routes } from 'react-router-dom';
import { AppShell } from '@/components/layout/AppShell';
import { HomePage } from '@/pages/HomePage';
import { LoginPage } from '@/pages/LoginPage';
import { RegisterPage } from '@/pages/RegisterPage';
import { PasswordResetRequestPage } from '@/pages/PasswordResetRequestPage';
import { PasswordResetConfirmPage } from '@/pages/PasswordResetConfirmPage';
import { CarSearchResultsPage } from '@/pages/CarSearchResultsPage';
import { BookCarPage } from '@/pages/BookCarPage';
import { BookingsPage } from '@/pages/BookingsPage';
import { ProfilePage } from '@/pages/ProfilePage';
import { NotFoundPage } from '@/pages/NotFoundPage';
import { AdminLayout } from '@/pages/AdminLayout';
import { AdminUsersPage } from '@/pages/AdminUsersPage';
import { AdminCarsPage } from '@/pages/AdminCarsPage';
import { AdminBookingsPage } from '@/pages/AdminBookingsPage';
import { ProtectedRoute } from './ProtectedRoute';
import { AdminRoute } from './AdminRoute';

/**
 * Auth pages have their own centred layout (no Navbar). All other routes
 * render inside the AppShell.
 */
export function AppRouter() {
  return (
    <Routes>
      {/* Auth routes (own layout) */}
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<RegisterPage />} />
      <Route path="/password/reset" element={<PasswordResetRequestPage />} />
      <Route path="/password/confirm" element={<PasswordResetConfirmPage />} />

      {/* App routes (shared shell) */}
      <Route element={<AppShellRoute />}>
        <Route index element={<HomePage />} />
        <Route path="/cars" element={<CarSearchResultsPage />} />
        <Route
          path="/cars/:carId/book"
          element={
            <ProtectedRoute>
              <BookCarPage />
            </ProtectedRoute>
          }
        />
        <Route
          path="/bookings"
          element={
            <ProtectedRoute>
              <BookingsPage />
            </ProtectedRoute>
          }
        />
        <Route
          path="/profile"
          element={
            <ProtectedRoute>
              <ProfilePage />
            </ProtectedRoute>
          }
        />
        <Route
          path="/admin"
          element={
            <AdminRoute>
              <AdminLayout />
            </AdminRoute>
          }
        >
          <Route index element={<Navigate to="users" replace />} />
          <Route path="users" element={<AdminUsersPage />} />
          <Route path="cars" element={<AdminCarsPage />} />
          <Route path="bookings" element={<AdminBookingsPage />} />
        </Route>
      </Route>

      <Route path="*" element={<NotFoundPage />} />
    </Routes>
  );
}

function AppShellRoute() {
  return (
    <AppShell>
      <Outlet />
    </AppShell>
  );
}
