import { AuthLayout } from './AuthLayout';
import { PasswordResetConfirmForm } from '@/features/auth/PasswordResetConfirmForm';

export function PasswordResetConfirmPage() {
  return (
    <AuthLayout title="Set a new password" subtitle="Pick something at least 12 characters long.">
      <PasswordResetConfirmForm />
    </AuthLayout>
  );
}
