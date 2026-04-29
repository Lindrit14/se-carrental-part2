import { useMemo, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Car as CarIcon } from 'lucide-react';
import { Button } from '@/components/ui/Button';
import { Skeleton } from '@/components/ui/Skeleton';
import { EmptyState } from '@/components/ui/EmptyState';
import type { Car } from '@/domain/car';
import { useAuth } from '@/features/auth/useAuth';
import { CarCard } from './CarCard';

type SortKey = 'recommended' | 'price-asc' | 'price-desc';

interface CarResultsListProps {
  cars: Car[] | undefined;
  isLoading: boolean;
  isError: boolean;
  refetch: () => void;
  /** When set, "View deal" forwards these dates into the booking page URL. */
  bookingDates?: { from: string; to: string };
}

export function CarResultsList({ cars, isLoading, isError, refetch, bookingDates }: CarResultsListProps) {
  const [sort, setSort] = useState<SortKey>('recommended');
  const { session } = useAuth();
  const navigate = useNavigate();

  const sorted = useMemo(() => {
    if (!cars) return [];
    if (sort === 'recommended') return cars;
    const copy = [...cars];
    const num = (c: Car) => Number.parseFloat(c.dailyRate.amount) || 0;
    copy.sort((a, b) => (sort === 'price-asc' ? num(a) - num(b) : num(b) - num(a)));
    return copy;
  }, [cars, sort]);

  if (isLoading) {
    return (
      <div className="grid grid-cols-1 gap-4">
        {Array.from({ length: 4 }).map((_, i) => (
          <Skeleton key={i} className="h-40 w-full" />
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
        action={<Button onClick={refetch}>Try again</Button>}
      />
    );
  }

  if (!cars || cars.length === 0) {
    return (
      <EmptyState
        icon={<CarIcon className="h-8 w-8" />}
        title="No cars match this search"
        description="Try a different location, dates, or category."
      />
    );
  }

  const handleViewDeal = (car: Car) => {
    const params = new URLSearchParams();
    if (bookingDates?.from) params.set('from', bookingDates.from);
    if (bookingDates?.to) params.set('to', bookingDates.to);
    const target = `/cars/${car.id}/book${params.size ? `?${params.toString()}` : ''}`;
    if (session) {
      navigate(target);
    } else {
      navigate(`/login?returnTo=${encodeURIComponent(target)}`);
    }
  };

  return (
    <div className="flex flex-col gap-4">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold text-zinc-900">{cars.length} cars available</h2>
        <label className="flex items-center gap-2 text-sm text-zinc-600">
          <span>Sort by:</span>
          <select
            className="rounded-md border border-zinc-200 bg-white px-2 py-1 text-sm"
            value={sort}
            onChange={(e) => setSort(e.target.value as SortKey)}
          >
            <option value="recommended">Recommended</option>
            <option value="price-asc">Price: low to high</option>
            <option value="price-desc">Price: high to low</option>
          </select>
        </label>
      </div>

      <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
        {sorted.map((car) => (
          <CarCard
            key={car.id}
            car={car}
            action={
              <Button size="sm" className="w-full bg-emerald-600 hover:bg-emerald-700" onClick={() => handleViewDeal(car)}>
                View deal
              </Button>
            }
          />
        ))}
      </div>
    </div>
  );
}
