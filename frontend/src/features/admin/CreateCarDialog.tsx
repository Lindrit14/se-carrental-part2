import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { toast } from 'sonner';
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
import { FormError, FormField, FormLabel } from '@/components/ui/Form';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/Select';
import { useCreateCar } from './useAdmin';

const CURRENCIES = ['EUR', 'USD', 'GBP', 'CHF'];

const schema = z.object({
  brand: z.string().min(1, 'Required').max(80),
  model: z.string().min(1, 'Required').max(120),
  licensePlate: z.string().min(1, 'Required').max(40),
  dailyRateAmount: z
    .string()
    .regex(/^\d+(\.\d{1,2})?$/, 'Use a positive number, e.g. 79.00'),
  dailyRateCurrency: z.string().regex(/^[A-Z]{3}$/, 'Pick a currency'),
});

type FormValues = z.infer<typeof schema>;

interface CreateCarDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function CreateCarDialog({ open, onOpenChange }: CreateCarDialogProps) {
  const create = useCreateCar();
  const [serverError, setServerError] = useState<string | null>(null);

  const {
    register,
    handleSubmit,
    reset,
    setValue,
    watch,
    formState: { errors },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: {
      brand: '',
      model: '',
      licensePlate: '',
      dailyRateAmount: '',
      dailyRateCurrency: 'EUR',
    },
  });

  const currency = watch('dailyRateCurrency');

  const onSubmit = handleSubmit(async (values) => {
    setServerError(null);
    try {
      await create.mutateAsync(values);
      toast.success('Car added to fleet');
      reset();
      onOpenChange(false);
    } catch {
      setServerError('Could not create the car.');
    }
  });

  return (
    <Dialog
      open={open}
      onOpenChange={(o) => {
        onOpenChange(o);
        if (!o) {
          reset();
          setServerError(null);
        }
      }}
    >
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Add a car</DialogTitle>
          <DialogDescription>
            New cars are immediately bookable by all users.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={onSubmit} className="flex flex-col gap-4" noValidate>
          <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
            <FormField>
              <FormLabel htmlFor="brand">Brand</FormLabel>
              <Input id="brand" aria-invalid={!!errors.brand} {...register('brand')} />
              <FormError>{errors.brand?.message}</FormError>
            </FormField>
            <FormField>
              <FormLabel htmlFor="model">Model</FormLabel>
              <Input id="model" aria-invalid={!!errors.model} {...register('model')} />
              <FormError>{errors.model?.message}</FormError>
            </FormField>
          </div>

          <FormField>
            <FormLabel htmlFor="licensePlate">License plate</FormLabel>
            <Input
              id="licensePlate"
              aria-invalid={!!errors.licensePlate}
              {...register('licensePlate')}
            />
            <FormError>{errors.licensePlate?.message}</FormError>
          </FormField>

          <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
            <FormField>
              <FormLabel htmlFor="dailyRateAmount">Daily rate</FormLabel>
              <Input
                id="dailyRateAmount"
                inputMode="decimal"
                aria-invalid={!!errors.dailyRateAmount}
                placeholder="79.00"
                {...register('dailyRateAmount')}
              />
              <FormError>{errors.dailyRateAmount?.message}</FormError>
            </FormField>
            <FormField>
              <FormLabel>Currency</FormLabel>
              <Select
                value={currency}
                onValueChange={(v) =>
                  setValue('dailyRateCurrency', v, { shouldValidate: true })
                }
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
              <FormError>{errors.dailyRateCurrency?.message}</FormError>
            </FormField>
          </div>

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
              {create.isPending ? 'Saving…' : 'Add car'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
