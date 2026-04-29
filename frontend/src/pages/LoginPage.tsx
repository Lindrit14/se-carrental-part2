import { Link } from 'react-router-dom';
import { AuthLayout } from './AuthLayout';
import { LoginForm } from '@/features/auth/LoginForm';

export function LoginPage() {
  return (
    <AuthLayout
      title="Welcome back"
      subtitle="Sign in to manage your bookings."
      footer={
        <>
          New here?{' '}
          <Link to="/register" className="font-medium text-zinc-900 hover:underline">
            Create an account
          </Link>
        </>
      }
    >
      <LoginForm />
    </AuthLayout>
  );
}
