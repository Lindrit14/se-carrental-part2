import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { toast } from 'sonner';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { FormError, FormField, FormLabel } from '@/components/ui/Form';
import { authApi } from '@/lib/api/authApi';

const schema = z.object({ email: z.string().email('Enter a valid email') });
type FormValues = z.infer<typeof schema>;

export function PasswordResetRequestForm() {
  const [submitted, setSubmitted] = useState(false);
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<FormValues>({ resolver: zodResolver(schema) });

  const onSubmit = handleSubmit(async (values) => {
    try {
      await authApi.requestPasswordReset(values.email);
    } catch {
      // Backend returns 202 even for unknown emails — silent failure here is fine.
    }
    setSubmitted(true);
    toast.success('If that account exists, we sent reset instructions.');
  });

  if (submitted) {
    return (
      <div className="rounded-[var(--radius-card)] border border-zinc-200 bg-white p-6 text-sm text-zinc-700">
        Check your inbox for a reset link. The link expires in 30 minutes.
      </div>
    );
  }

  return (
    <form onSubmit={onSubmit} className="flex flex-col gap-4" noValidate>
      <FormField>
        <FormLabel htmlFor="email">Email</FormLabel>
        <Input id="email" type="email" autoComplete="email" {...register('email')} />
        <FormError>{errors.email?.message}</FormError>
      </FormField>
      <Button type="submit" disabled={isSubmitting}>
        {isSubmitting ? 'Sending…' : 'Send reset link'}
      </Button>
    </form>
  );
}
