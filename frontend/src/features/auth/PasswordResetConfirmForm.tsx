import { useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { toast } from 'sonner';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { FormError, FormField, FormLabel } from '@/components/ui/Form';
import { authApi } from '@/lib/api/authApi';
import { ApiError } from '@/lib/api/errors';

const schema = z.object({ password: z.string().min(12, 'At least 12 characters') });
type FormValues = z.infer<typeof schema>;

export function PasswordResetConfirmForm() {
  const [params] = useSearchParams();
  const token = params.get('token') ?? '';
  const navigate = useNavigate();
  const [submitting, setSubmitting] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<FormValues>({ resolver: zodResolver(schema) });

  if (!token) {
    return (
      <div className="rounded-[var(--radius-card)] border border-zinc-200 bg-white p-6 text-sm text-zinc-700">
        This reset link is missing a token. Request a new one from the password reset page.
      </div>
    );
  }

  const onSubmit = handleSubmit(async (values) => {
    setSubmitting(true);
    try {
      await authApi.confirmPasswordReset(token, values.password);
      toast.success('Password updated. Please sign in.');
      navigate('/login', { replace: true });
    } catch (err) {
      const message =
        err instanceof ApiError && err.code === 'invalid_token'
          ? 'This reset link is invalid or expired.'
          : 'Could not update your password. Please try again.';
      toast.error(message);
    } finally {
      setSubmitting(false);
    }
  });

  return (
    <form onSubmit={onSubmit} className="flex flex-col gap-4" noValidate>
      <FormField>
        <FormLabel htmlFor="password">New password</FormLabel>
        <Input
          id="password"
          type="password"
          autoComplete="new-password"
          aria-invalid={!!errors.password}
          {...register('password')}
        />
        <FormError>{errors.password?.message}</FormError>
      </FormField>
      <Button type="submit" disabled={submitting}>
        {submitting ? 'Updating…' : 'Update password'}
      </Button>
    </form>
  );
}
