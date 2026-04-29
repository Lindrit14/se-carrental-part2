import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { toast } from 'sonner';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { FormError, FormField, FormHelp, FormLabel } from '@/components/ui/Form';
import { ApiError } from '@/lib/api/errors';
import { useAuth } from './useAuth';

const schema = z.object({
  email: z.string().email('Enter a valid email'),
  password: z.string().min(12, 'At least 12 characters'),
});

type FormValues = z.infer<typeof schema>;

export function RegisterForm() {
  const { register: registerUser } = useAuth();
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
      await registerUser(values.email, values.password);
      toast.success('Account created');
      navigate('/cars', { replace: true });
    } catch (err) {
      const message =
        err instanceof ApiError && err.code === 'email_taken'
          ? 'That email is already registered.'
          : 'Sign-up failed. Please try again.';
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
        <FormLabel htmlFor="password">Password</FormLabel>
        <Input
          id="password"
          type="password"
          autoComplete="new-password"
          aria-invalid={!!errors.password}
          {...register('password')}
        />
        <FormHelp>Use at least 12 characters.</FormHelp>
        <FormError>{errors.password?.message}</FormError>
      </FormField>

      <Button type="submit" disabled={submitting} className="mt-2">
        {submitting ? 'Creating account…' : 'Create account'}
      </Button>
    </form>
  );
}
