import { Link, useLocation } from 'react-router-dom';
import { AuthLayout } from './AuthLayout';
import { LoginForm } from '@/features/auth/LoginForm';

export function LoginPage() {
  // Carry returnTo through the "Create an account" link so users who came
  // here from a "View deal" click don't lose their destination.
  const location = useLocation();
  const registerHref = `/register${location.search}`;

  return (
    <AuthLayout
      title="Welcome back"
      subtitle="Sign in to manage your bookings."
      footer={
        <>
          New here?{' '}
          <Link to={registerHref} className="font-medium text-zinc-900 hover:underline">
            Create an account
          </Link>
        </>
      }
    >
      <LoginForm />
    </AuthLayout>
  );
}
