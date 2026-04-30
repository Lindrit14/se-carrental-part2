import { useState } from 'react';
import { Car as CarIcon, Pencil, Trash2 } from 'lucide-react';
import { toast } from 'sonner';
import { Button } from '@/components/ui/Button';
import { Skeleton } from '@/components/ui/Skeleton';
import { EmptyState } from '@/components/ui/EmptyState';
import { CarCard } from './CarCard';
import { useCars } from './useCars';
import { useDeleteCar } from '@/features/admin/useAdmin';
import { EditCarDialog } from '@/features/admin/EditCarDialog';
import { ApiError } from '@/lib/api/errors';
import type { Car } from '@/domain/car';

/**
 * Admin-only list of every car. Public car browsing happens through
 * {@code CarSearchResultsPage} + {@code CarResultsList}.
 */
export function CarList() {
  const { data, isLoading, isError, refetch } = useCars();
  const deleteCar = useDeleteCar();
  const [editing, setEditing] = useState<Car | null>(null);

  const onDelete = (car: Car) => {
    if (!window.confirm(`Delete ${car.brand} ${car.model} (${car.licensePlate})?`)) return;
    deleteCar.mutate(car.id, {
      onSuccess: () => toast.success('Car deleted'),
      onError: (err) => {
        if (err instanceof ApiError && err.code === 'car_has_bookings') {
          toast.error('Cannot delete: car has existing bookings.');
        } else {
          toast.error('Could not delete the car.');
        }
      },
    });
  };

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
        title="No cars in the fleet"
        description="Add a car to get started."
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
            showLicensePlate
            action={
              <div className="flex w-full justify-end gap-2">
                <Button variant="secondary" onClick={() => setEditing(car)}>
                  <Pencil className="h-4 w-4" /> Edit
                </Button>
                <Button
                  variant="secondary"
                  onClick={() => onDelete(car)}
                  disabled={deleteCar.isPending}
                >
                  <Trash2 className="h-4 w-4" /> Delete
                </Button>
              </div>
            }
          />
        ))}
      </div>
      <EditCarDialog
        car={editing}
        open={editing !== null}
        onOpenChange={(o) => {
          if (!o) setEditing(null);
        }}
      />
    </>
  );
}
