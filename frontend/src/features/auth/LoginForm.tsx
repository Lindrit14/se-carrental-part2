import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { toast } from 'sonner';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { FormError, FormField, FormLabel } from '@/components/ui/Form';
import { ApiError } from '@/lib/api/errors';
import { useAuth } from './useAuth';

const schema = z.object({
  email: z.string().email('Enter a valid email'),
  password: z.string().min(1, 'Password is required'),
});

type FormValues = z.infer<typeof schema>;

export function LoginForm() {
  const { login } = useAuth();
  const navigate = useNavigate();
  const [submitting, setSubmitting] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<FormValues>({ resolver: zodResolver(schema) });

  const onSubmit = handleSubmit(async (values) => {
    setSubmitting(true);
    try {
      await login(values.email, values.password);
      toast.success('Welcome back');
      navigate('/cars', { replace: true });
    } catch (err) {
      const message =
        err instanceof ApiError && err.code === 'invalid_credentials'
          ? 'Email or password is incorrect.'
          : 'Could not sign you in. Please try again.';
      toast.error(message);
    } finally {
      setSubmitting(false);
    }
  });

  return (
    <form onSubmit={onSubmit} className="flex flex-col gap-4" noValidate>
      <FormField>
        <FormLabel htmlFor="email">Email</FormLabel>
        <Input
          id="email"
          type="email"
          autoComplete="email"
          aria-invalid={!!errors.email}
          {...register('email')}
        />
        <FormError>{errors.email?.message}</FormError>
      </FormField>

      <FormField>
        <div className="flex items-center justify-between">
          <FormLabel htmlFor="password">Password</FormLabel>
          <Link
            to="/password/reset"
            className="text-xs font-medium text-zinc-500 hover:text-zinc-900"
          >
            Forgot password?
          </Link>
        </div>
        <Input
          id="password"
          type="password"
          autoComplete="current-password"
          aria-invalid={!!errors.password}
          {...register('password')}
        />
        <FormError>{errors.password?.message}</FormError>
      </FormField>

      <Button type="submit" disabled={submitting} className="mt-2">
        {submitting ? 'Signing in…' : 'Sign in'}
      </Button>
    </form>
  );
}
