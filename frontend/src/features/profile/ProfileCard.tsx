import { useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { toast } from 'sonner';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { Skeleton } from '@/components/ui/Skeleton';
import { Badge } from '@/components/ui/Badge';
import { FormError, FormField, FormLabel } from '@/components/ui/Form';
import { useMyProfile, useUpdateProfile } from './useProfile';
import { ApiError } from '@/lib/api/errors';

const schema = z.object({ email: z.string().email() });
type FormValues = z.infer<typeof schema>;

export function ProfileCard() {
  const { data, isLoading } = useMyProfile();
  const update = useUpdateProfile();

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isDirty },
  } = useForm<FormValues>({ resolver: zodResolver(schema), defaultValues: { email: '' } });

  useEffect(() => {
    if (data) reset({ email: data.email });
  }, [data, reset]);

  if (isLoading || !data) {
    return <Skeleton className="h-72 w-full" />;
  }

  const onSubmit = handleSubmit(async (values) => {
    try {
      await update.mutateAsync({ email: values.email });
      toast.success('Profile updated');
    } catch (err) {
      const message =
        err instanceof ApiError && err.code === 'email_taken'
          ? 'That email is already taken.'
          : 'Could not update profile.';
      toast.error(message);
    }
  });

  return (
    <Card>
      <CardHeader>
        <CardTitle>Account details</CardTitle>
        <CardDescription>Update your contact information.</CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={onSubmit} className="flex flex-col gap-4">
          <div className="flex flex-wrap items-center gap-2">
            {data.roles.map((r) => (
              <Badge key={r} variant="muted">
                {r}
              </Badge>
            ))}
            <Badge variant={data.verified ? 'success' : 'warning'}>
              {data.verified ? 'Verified' : 'Unverified'}
            </Badge>
          </div>

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

          <div className="flex justify-end">
            <Button type="submit" disabled={!isDirty || update.isPending}>
              {update.isPending ? 'Saving…' : 'Save changes'}
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  );
}
