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
import type { Car } from '@/domain/car';
import { useDisplayCurrency } from '@/features/currency/useDisplayCurrency';
import { formatMoney } from '@/lib/format/money';
import { useCreateBooking } from './useBookings';
import { ApiError } from '@/lib/api/errors';

const schema = z
  .object({
    startDate: z.string().min(1, 'Start date is required'),
    endDate: z.string().min(1, 'End date is required'),
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
  defaultStartDate?: string;
  defaultEndDate?: string;
}

export function CreateBookingDialog({
  car,
  open,
  onOpenChange,
  defaultStartDate,
  defaultEndDate,
}: CreateBookingDialogProps) {
  const create = useCreateBooking();
  const navigate = useNavigate();
  const { displayCurrency } = useDisplayCurrency();
  const [serverError, setServerError] = useState<string | null>(null);

  const today = new Date().toISOString().slice(0, 10);
  const initialStart = defaultStartDate || today;
  const initialEnd = defaultEndDate || '';

  const {
    register,
    handleSubmit,
    reset,
    watch,
    formState: { errors },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { startDate: initialStart, endDate: initialEnd },
  });

  // Reset form whenever the target car or default dates change.
  useEffect(() => {
    if (car) {
      reset({ startDate: initialStart, endDate: initialEnd });
      setServerError(null);
    }
  }, [car, reset, initialStart, initialEnd]);

  const startDate = watch('startDate');
  const endDate = watch('endDate');
  const days =
    startDate && endDate && endDate > startDate
      ? Math.round(
          (new Date(endDate).getTime() - new Date(startDate).getTime()) / (1000 * 60 * 60 * 24),
        )
      : 0;

  const onSubmit = handleSubmit(async (values) => {
    if (!car) return;
    setServerError(null);
    try {
      await create.mutateAsync({
        carId: car.id,
        startDate: values.startDate,
        endDate: values.endDate,
        targetCurrency: displayCurrency,
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

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Book {car.brand} {car.model}</DialogTitle>
          <DialogDescription>
            Daily rate: {car.dailyRate.amount} {car.dailyRate.currency}. Total will be charged
            in {displayCurrency} (set in the navbar).
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

          {days > 0 && (
            <div className="flex items-baseline justify-between rounded-md border border-zinc-200 bg-zinc-50 px-3 py-2.5">
              <span className="text-sm text-zinc-600">
                Estimated total ({days} day{days === 1 ? '' : 's'})
              </span>
              <span className="font-semibold text-zinc-900">
                {formatMoney({
                  amount: (Number(car.dailyRate.amount) * days).toFixed(2),
                  currency: car.dailyRate.currency,
                })}
              </span>
            </div>
          )}

          <FormHelp>
            The total will be locked in at the current ECB reference rate when the booking is
            confirmed and shown in {displayCurrency} on your bookings page.
          </FormHelp>

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
