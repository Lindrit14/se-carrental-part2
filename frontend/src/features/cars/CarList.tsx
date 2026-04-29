import { Car as CarIcon } from 'lucide-react';
import { Button } from '@/components/ui/Button';
import { Skeleton } from '@/components/ui/Skeleton';
import { EmptyState } from '@/components/ui/EmptyState';
import { CarCard } from './CarCard';
import { useCars } from './useCars';

/**
 * Admin-only list of every car. Public car browsing happens through
 * {@code CarSearchResultsPage} + {@code CarResultsList}.
 */
export function CarList() {
  const { data, isLoading, isError, refetch } = useCars();

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
    <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
      {data.map((car) => (
        <CarCard key={car.id} car={car} showLicensePlate />
      ))}
    </div>
  );
}
