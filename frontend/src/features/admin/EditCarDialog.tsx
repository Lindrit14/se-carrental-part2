import { useEffect, useState } from 'react';
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
import { CAR_CATEGORIES, type Car, type CarCategory } from '@/domain/car';
import { useUpdateCar } from './useAdmin';



const CATEGORY_LABEL: Record<CarCategory, string> = {
  SMALL: 'Small car',
  MEDIUM: 'Medium car',
  LARGE: 'Large car',
  SUV: 'SUV',
  PEOPLE_CARRIER: 'People carrier',
  PREMIUM: 'Premium car',
};

const schema = z.object({
  brand: z.string().min(1, 'Required').max(80),
  model: z.string().min(1, 'Required').max(120),
  licensePlate: z.string().min(1, 'Required').max(40),
  dailyRateAmount: z
    .string()
    .regex(/^\d+(\.\d{1,2})?$/, 'Use a positive number, e.g. 79.00'),
  dailyRateCurrency: z.string().regex(/^[A-Z]{3}$/, 'Pick a currency'),
  location: z.string().min(1, 'Required').max(160),
  category: z.enum(['SMALL', 'MEDIUM', 'LARGE', 'SUV', 'PEOPLE_CARRIER', 'PREMIUM']),
});

type FormValues = z.infer<typeof schema>;

interface EditCarDialogProps {
  car: Car | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function EditCarDialog({ car, open, onOpenChange }: EditCarDialogProps) {
  const update = useUpdateCar();
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
      location: '',
      category: 'MEDIUM',
    },
  });

  // Hydrate form whenever the dialog opens with a different car.
  useEffect(() => {
    if (open && car) {
      reset({
        brand: car.brand,
        model: car.model,
        licensePlate: car.licensePlate,
        dailyRateAmount: String(car.dailyRate.amount),
        dailyRateCurrency: car.dailyRate.currency,
        location: car.location,
        category: car.category,
      });
      setServerError(null);
    }
  }, [open, car, reset]);

  const category = watch('category');

  const onSubmit = handleSubmit(async (values) => {
    if (!car) return;
    setServerError(null);
    try {
      await update.mutateAsync({ id: car.id, input: values });
      toast.success('Car updated');
      onOpenChange(false);
    } catch {
      setServerError('Could not update the car.');
    }
  });

  return (
    <Dialog
      open={open}
      onOpenChange={(o) => {
        onOpenChange(o);
        if (!o) setServerError(null);
      }}
    >
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Edit car</DialogTitle>
          <DialogDescription>Changes apply immediately.</DialogDescription>
        </DialogHeader>

        <form onSubmit={onSubmit} className="flex flex-col gap-4" noValidate>
          <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
            <FormField>
              <FormLabel htmlFor="edit-brand">Brand</FormLabel>
              <Input id="edit-brand" aria-invalid={!!errors.brand} {...register('brand')} />
              <FormError>{errors.brand?.message}</FormError>
            </FormField>
            <FormField>
              <FormLabel htmlFor="edit-model">Model</FormLabel>
              <Input id="edit-model" aria-invalid={!!errors.model} {...register('model')} />
              <FormError>{errors.model?.message}</FormError>
            </FormField>
          </div>

          <FormField>
            <FormLabel htmlFor="edit-licensePlate">License plate</FormLabel>
            <Input
              id="edit-licensePlate"
              aria-invalid={!!errors.licensePlate}
              {...register('licensePlate')}
            />
            <FormError>{errors.licensePlate?.message}</FormError>
          </FormField>

          <FormField>
            <FormLabel htmlFor="edit-location">Location</FormLabel>
            <Input
              id="edit-location"
              aria-invalid={!!errors.location}
              {...register('location')}
            />
            <FormError>{errors.location?.message}</FormError>
          </FormField>

          <FormField>
            <FormLabel>Category</FormLabel>
            <Select
              value={category}
              onValueChange={(v) =>
                setValue('category', v as CarCategory, { shouldValidate: true })
              }
            >
              <SelectTrigger>
                <SelectValue placeholder="Pick a category" />
              </SelectTrigger>
              <SelectContent>
                {CAR_CATEGORIES.map((c) => (
                  <SelectItem key={c} value={c}>
                    {CATEGORY_LABEL[c]}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <FormError>{errors.category?.message}</FormError>
          </FormField>

          <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
            <FormField>
              <FormLabel htmlFor="edit-dailyRateAmount">Daily rate</FormLabel>
              <Input
                id="edit-dailyRateAmount"
                inputMode="decimal"
                aria-invalid={!!errors.dailyRateAmount}
                {...register('dailyRateAmount')}
              />
              <FormError>{errors.dailyRateAmount?.message}</FormError>
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
            <Button type="submit" disabled={update.isPending}>
              {update.isPending ? 'Saving…' : 'Save changes'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
