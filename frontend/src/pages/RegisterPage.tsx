import { Link } from 'react-router-dom';
import { AuthLayout } from './AuthLayout';
import { RegisterForm } from '@/features/auth/RegisterForm';

export function RegisterPage() {
  return (
    <AuthLayout
      title="Create your account"
      subtitle="It takes less than a minute."
      footer={
        <>
          Already have an account?{' '}
          <Link to="/login" className="font-medium text-zinc-900 hover:underline">
            Sign in
          </Link>
        </>
      }
    >
      <RegisterForm />
    </AuthLayout>
  );
}
