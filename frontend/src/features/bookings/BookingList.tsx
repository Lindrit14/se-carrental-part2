import { Calendar } from 'lucide-react';
import { Link } from 'react-router-dom';
import { toast } from 'sonner';
import { Button } from '@/components/ui/Button';
import { Skeleton } from '@/components/ui/Skeleton';
import { EmptyState } from '@/components/ui/EmptyState';
import { BookingCard } from './BookingCard';
import { useCancelBooking, useMyBookings } from './useBookings';

export function BookingList() {
  const { data, isLoading, isError, refetch } = useMyBookings();
  const cancel = useCancelBooking();

  const onCancel = (id: string) => {
    cancel.mutate(id, {
      onSuccess: () => toast.success('Booking cancelled'),
      onError: () => toast.error('Could not cancel booking'),
    });
  };

  if (isLoading) {
    return (
      <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
        {Array.from({ length: 4 }).map((_, i) => (
          <Skeleton key={i} className="h-44 w-full" />
        ))}
      </div>
    );
  }

  if (isError) {
    return (
      <EmptyState
        icon={<Calendar className="h-8 w-8" />}
        title="Could not load your bookings"
        action={<Button onClick={() => refetch()}>Try again</Button>}
      />
    );
  }

  if (!data || data.length === 0) {
    return (
      <EmptyState
        icon={<Calendar className="h-8 w-8" />}
        title="No bookings yet"
        description="Browse the fleet and reserve your first car."
        action={
          <Button asChild>
            <Link to="/cars">Browse cars</Link>
          </Button>
        }
      />
    );
  }

  return (
    <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
      {data.map((b) => (
        <BookingCard
          key={b.id}
          booking={b}
          onCancel={onCancel}
          cancelling={cancel.isPending && cancel.variables === b.id}
        />
      ))}
    </div>
  );
}
