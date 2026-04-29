import { CalendarRange, AlertCircle } from 'lucide-react';
import { Badge } from '@/components/ui/Badge';
import { Button } from '@/components/ui/Button';
import { Skeleton } from '@/components/ui/Skeleton';
import { EmptyState } from '@/components/ui/EmptyState';
import type { BookingStatus } from '@/domain/booking';
import { useAdminBookings } from './useAdmin';

const STATUS_VARIANT: Record<BookingStatus, 'success' | 'warning' | 'destructive'> = {
  CONFIRMED: 'success',
  PENDING: 'warning',
  CANCELLED: 'destructive',
};

function fmtDate(iso: string): string {
  const d = new Date(iso);
  return Number.isNaN(d.getTime()) ? iso : d.toLocaleString();
}

export function AdminBookingTable() {
  const { data, isLoading, isError, refetch } = useAdminBookings();

  if (isLoading) {
    return (
      <div className="flex flex-col gap-2">
        {Array.from({ length: 6 }).map((_, i) => (
          <Skeleton key={i} className="h-14 w-full" />
        ))}
      </div>
    );
  }

  if (isError) {
    return (
      <EmptyState
        icon={<AlertCircle className="h-8 w-8" />}
        title="Could not load bookings"
        description="Something went wrong reaching the booking service."
        action={<Button onClick={() => refetch()}>Try again</Button>}
      />
    );
  }

  if (!data || data.length === 0) {
    return (
      <EmptyState
        icon={<CalendarRange className="h-8 w-8" />}
        title="No bookings yet"
        description="No customers have made a booking."
      />
    );
  }

  return (
    <div className="overflow-hidden rounded-lg border border-zinc-200 bg-white">
      <table className="w-full text-sm">
        <thead className="border-b border-zinc-200 bg-zinc-50 text-left text-zinc-500">
          <tr>
            <th className="px-4 py-3 font-medium">Booking</th>
            <th className="px-4 py-3 font-medium">Customer</th>
            <th className="px-4 py-3 font-medium">Car</th>
            <th className="px-4 py-3 font-medium">Dates</th>
            <th className="px-4 py-3 font-medium">Status</th>
            <th className="px-4 py-3 font-medium text-right">Total</th>
          </tr>
        </thead>
        <tbody>
          {data.map((b) => (
            <tr key={b.id} className="border-b border-zinc-100 last:border-0">
              <td className="px-4 py-3">
                <div className="font-medium text-zinc-900">{b.id.slice(0, 8)}…</div>
                <div className="text-xs text-zinc-500">{fmtDate(b.createdAt)}</div>
              </td>
              <td className="px-4 py-3 font-mono text-xs text-zinc-600">
                {b.customerId.slice(0, 8)}…
              </td>
              <td className="px-4 py-3 font-mono text-xs text-zinc-600">
                {b.carId.slice(0, 8)}…
              </td>
              <td className="px-4 py-3 text-zinc-700">
                {b.startDate} → {b.endDate}
              </td>
              <td className="px-4 py-3">
                <Badge variant={STATUS_VARIANT[b.status]}>{b.status.toLowerCase()}</Badge>
              </td>
              <td className="px-4 py-3 text-right text-zinc-900">
                {b.totalTarget.amount} {b.totalTarget.currency}
                <div className="text-xs text-zinc-500">
                  ({b.totalSource.amount} {b.totalSource.currency})
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
