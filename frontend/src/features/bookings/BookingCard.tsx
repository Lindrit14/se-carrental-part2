import { Calendar } from 'lucide-react';
import type { Booking } from '@/domain/booking';
import { Card, CardContent, CardFooter } from '@/components/ui/Card';
import { Badge, type BadgeProps } from '@/components/ui/Badge';
import { Button } from '@/components/ui/Button';
import { formatMoney } from '@/lib/format/money';
import { formatDateRange, daysBetween } from '@/lib/format/date';

const STATUS_VARIANTS: Record<Booking['status'], BadgeProps['variant']> = {
  PENDING: 'warning',
  CONFIRMED: 'success',
  CANCELLED: 'muted',
};

const STATUS_LABEL: Record<Booking['status'], string> = {
  PENDING: 'Pending',
  CONFIRMED: 'Confirmed',
  CANCELLED: 'Cancelled',
};

interface BookingCardProps {
  booking: Booking;
  onCancel?: (id: string) => void;
  cancelling?: boolean;
}

export function BookingCard({ booking, onCancel, cancelling }: BookingCardProps) {
  const days = daysBetween(booking.startDate, booking.endDate);
  const sameCurrency = booking.totalSource.currency === booking.totalTarget.currency;

  return (
    <Card>
      <CardContent className="flex flex-col gap-4 p-6">
        <div className="flex items-start justify-between gap-2">
          <div>
            <p className="text-xs uppercase tracking-wide text-zinc-500">Booking</p>
            <p className="font-mono text-xs text-zinc-700">{booking.id.slice(0, 8)}…</p>
          </div>
          <Badge variant={STATUS_VARIANTS[booking.status]}>
            {STATUS_LABEL[booking.status]}
          </Badge>
        </div>

        <div className="flex items-center gap-2 text-sm text-zinc-700">
          <Calendar className="h-4 w-4 text-zinc-400" />
          <span>{formatDateRange(booking.startDate, booking.endDate)}</span>
          <span className="text-zinc-400">·</span>
          <span className="text-zinc-500">{days} day{days === 1 ? '' : 's'}</span>
        </div>

        <div className="flex flex-col gap-1 border-t border-zinc-100 pt-4">
          <div className="flex items-baseline justify-between text-sm">
            <span className="text-zinc-500">Total</span>
            <span className="font-semibold text-zinc-900">
              {formatMoney(booking.totalTarget)}
            </span>
          </div>
          {!sameCurrency && (
            <p className="text-xs text-zinc-500">
              Original: {formatMoney(booking.totalSource)}
            </p>
          )}
        </div>
      </CardContent>
      {booking.status !== 'CANCELLED' && onCancel && (
        <CardFooter className="border-t border-zinc-100 px-6 py-4">
          <Button
            variant="destructive"
            size="sm"
            onClick={() => onCancel(booking.id)}
            disabled={cancelling}
          >
            {cancelling ? 'Cancelling…' : 'Cancel booking'}
          </Button>
        </CardFooter>
      )}
    </Card>
  );
}
