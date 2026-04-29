import { useState } from 'react';
import { Car as CarIcon } from 'lucide-react';
import { Button } from '@/components/ui/Button';
import { Skeleton } from '@/components/ui/Skeleton';
import { EmptyState } from '@/components/ui/EmptyState';
import { CarCard } from './CarCard';
import { useCars } from './useCars';
import type { Car } from '@/domain/car';
import { CreateBookingDialog } from '@/features/bookings/CreateBookingDialog';

export function CarList() {
  const { data, isLoading, isError, refetch } = useCars();
  const [bookingCar, setBookingCar] = useState<Car | null>(null);

  if (isLoading) {
    return (
      <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
        {Array.from({ length: 6 }).map((_, i) => (
          <Skeleton key={i} className="h-48 w-full" />
        ))}
      </div>
    );
  }

  if (isError) {
    return (
      <EmptyState
        icon={<CarIcon className="h-8 w-8" />}
        title="Could not load cars"
        description="Something went wrong while reaching the booking service."
        action={<Button onClick={() => refetch()}>Try again</Button>}
      />
    );
  }

  if (!data || data.length === 0) {
    return (
      <EmptyState
        icon={<CarIcon className="h-8 w-8" />}
        title="No cars available"
        description="Check back later — our fleet is being restocked."
      />
    );
  }

  return (
    <>
      <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
        {data.map((car) => (
          <CarCard
            key={car.id}
            car={car}
            action={
              <Button size="sm" className="w-full" onClick={() => setBookingCar(car)}>
                Book this car
              </Button>
            }
          />
        ))}
      </div>
      <CreateBookingDialog
        car={bookingCar}
        open={Boolean(bookingCar)}
        onOpenChange={(open) => !open && setBookingCar(null)}
      />
    </>
  );
}
