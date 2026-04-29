import { Link } from 'react-router-dom';
import { AuthLayout } from './AuthLayout';
import { PasswordResetRequestForm } from '@/features/auth/PasswordResetRequestForm';

export function PasswordResetRequestPage() {
  return (
    <AuthLayout
      title="Reset password"
      subtitle="We'll email you a link to set a new password."
      footer={
        <Link to="/login" className="font-medium text-zinc-900 hover:underline">
          Back to sign in
        </Link>
      }
    >
      <PasswordResetRequestForm />
    </AuthLayout>
  );
}
