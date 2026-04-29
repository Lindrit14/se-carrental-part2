import { useEffect, useState } from 'react';
import { useNavigate, useParams, useSearchParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { ArrowLeft, MapPin } from 'lucide-react';
import { bookingApi } from '@/lib/api/bookingApi';
import { Button } from '@/components/ui/Button';
import { Skeleton } from '@/components/ui/Skeleton';
import { EmptyState } from '@/components/ui/EmptyState';
import { Badge } from '@/components/ui/Badge';
import { Price } from '@/components/ui/Price';
import { CreateBookingDialog } from '@/features/bookings/CreateBookingDialog';
import type { CarCategory } from '@/domain/car';

const CATEGORY_LABEL: Record<CarCategory, string> = {
  SMALL: 'Small car',
  MEDIUM: 'Medium car',
  LARGE: 'Large car',
  SUV: 'SUV',
  PEOPLE_CARRIER: 'People carrier',
  PREMIUM: 'Premium car',
};

export function BookCarPage() {
  const { carId = '' } = useParams<{ carId: string }>();
  const [params] = useSearchParams();
  const navigate = useNavigate();
  const from = params.get('from') ?? '';
  const to = params.get('to') ?? '';

  const { data: car, isLoading, isError } = useQuery({
    queryKey: ['car', carId],
    queryFn: () => bookingApi.getCar(carId),
    enabled: !!carId,
  });

  // Auto-open the booking dialog as soon as the car is loaded — the page
  // exists to take the booking, no point in an extra click.
  const [open, setOpen] = useState(false);
  useEffect(() => {
    if (car) setOpen(true);
  }, [car]);

  if (isLoading) return <Skeleton className="h-48 w-full" />;
  if (isError || !car) {
    return (
      <EmptyState
        title="Car not found"
        description="This car is no longer available."
        action={<Button onClick={() => navigate('/cars')}>Back to search</Button>}
      />
    );
  }

  return (
    <div className="flex flex-col gap-6">
      <Button variant="ghost" size="sm" className="self-start" onClick={() => navigate(-1)}>
        <ArrowLeft className="h-4 w-4" /> Back
      </Button>

      <div className="rounded-lg border border-zinc-200 bg-white p-6">
        <div className="flex items-start justify-between gap-4">
          <div>
            <h1 className="text-2xl font-semibold tracking-tight text-zinc-900">
              {car.brand} {car.model}
            </h1>
            {car.location && (
              <p className="mt-1 inline-flex items-center gap-1 text-sm text-zinc-500">
                <MapPin className="h-4 w-4" /> {car.location}
              </p>
            )}
          </div>
          <Badge variant="muted">{CATEGORY_LABEL[car.category]}</Badge>
        </div>
        <div className="mt-4">
          <Price source={car.dailyRate} suffix="/ day" />
        </div>
      </div>

      <CreateBookingDialog
        car={car}
        open={open}
        onOpenChange={setOpen}
        defaultStartDate={from || undefined}
        defaultEndDate={to || undefined}
      />
    </div>
  );
}
