import { Link } from 'react-router-dom';
import { Button } from '@/components/ui/Button';
import { PageHeader } from '@/components/layout/PageHeader';
import { BookingList } from '@/features/bookings/BookingList';

export function BookingsPage() {
  return (
    <div className="flex flex-col gap-8">
      <PageHeader
        title="My bookings"
        description="Past and upcoming reservations."
        action={
          <Button asChild>
            <Link to="/cars">Browse cars</Link>
          </Button>
        }
      />
      <BookingList />
    </div>
  );
}
