import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { toast } from 'sonner';
import { useNavigate } from 'react-router-dom';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/Dialog';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { FormError, FormField, FormHelp, FormLabel } from '@/components/ui/Form';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/Select';
import type { Car } from '@/domain/car';
import { useCreateBooking } from './useBookings';
import { ApiError } from '@/lib/api/errors';

const CURRENCIES = ['EUR', 'USD', 'GBP', 'CHF'];

const schema = z
  .object({
    startDate: z.string().min(1, 'Start date is required'),
    endDate: z.string().min(1, 'End date is required'),
    targetCurrency: z.string().regex(/^[A-Z]{3}$/, 'Pick a currency'),
  })
  .refine((v) => v.endDate > v.startDate, {
    message: 'End date must be after start date',
    path: ['endDate'],
  });

type FormValues = z.infer<typeof schema>;

interface CreateBookingDialogProps {
  car: Car | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function CreateBookingDialog({ car, open, onOpenChange }: CreateBookingDialogProps) {
  const create = useCreateBooking();
  const navigate = useNavigate();
  const [serverError, setServerError] = useState<string | null>(null);

  const today = new Date().toISOString().slice(0, 10);

  const {
    register,
    handleSubmit,
    reset,
    setValue,
    watch,
    formState: { errors },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { startDate: today, endDate: '', targetCurrency: car?.dailyRate.currency ?? 'EUR' },
  });

  // Reset form whenever the target car changes.
  useEffect(() => {
    if (car) {
      reset({ startDate: today, endDate: '', targetCurrency: car.dailyRate.currency });
      setServerError(null);
    }
  }, [car, reset, today]);

  const onSubmit = handleSubmit(async (values) => {
    if (!car) return;
    setServerError(null);
    try {
      await create.mutateAsync({
        carId: car.id,
        startDate: values.startDate,
        endDate: values.endDate,
        targetCurrency: values.targetCurrency,
      });
      toast.success('Booking confirmed');
      onOpenChange(false);
      navigate('/bookings');
    } catch (err) {
      const message =
        err instanceof ApiError && err.code === 'currency_conversion_failed'
          ? 'The currency service is unavailable. Please try a different currency.'
          : 'Could not create the booking.';
      setServerError(message);
    }
  });

  if (!car) return null;

  const targetCurrency = watch('targetCurrency');

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Book {car.brand} {car.model}</DialogTitle>
          <DialogDescription>
            Daily rate: {car.dailyRate.amount} {car.dailyRate.currency}. We'll convert the
            total to your preferred currency.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={onSubmit} className="flex flex-col gap-4" noValidate>
          <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
            <FormField>
              <FormLabel htmlFor="startDate">Start date</FormLabel>
              <Input
                id="startDate"
                type="date"
                min={today}
                aria-invalid={!!errors.startDate}
                {...register('startDate')}
              />
              <FormError>{errors.startDate?.message}</FormError>
            </FormField>

            <FormField>
              <FormLabel htmlFor="endDate">End date</FormLabel>
              <Input
                id="endDate"
                type="date"
                min={today}
                aria-invalid={!!errors.endDate}
                {...register('endDate')}
              />
              <FormError>{errors.endDate?.message}</FormError>
            </FormField>
          </div>

          <FormField>
            <FormLabel>Total currency</FormLabel>
            <Select
              value={targetCurrency}
              onValueChange={(v) => setValue('targetCurrency', v, { shouldValidate: true })}
            >
              <SelectTrigger>
                <SelectValue placeholder="Pick a currency" />
              </SelectTrigger>
              <SelectContent>
                {CURRENCIES.map((c) => (
                  <SelectItem key={c} value={c}>
                    {c}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <FormHelp>
              The booking total will be converted via the currency-converter service.
            </FormHelp>
            <FormError>{errors.targetCurrency?.message}</FormError>
          </FormField>

          {serverError && (
            <p className="rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-700">
              {serverError}
            </p>
          )}

          <DialogFooter>
            <Button type="button" variant="secondary" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={create.isPending}>
              {create.isPending ? 'Booking…' : 'Confirm booking'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
