import { cn } from '@/lib/utils/cn';
import type { Money } from '@/domain/money';
import { formatMoney } from '@/lib/format/money';

interface PriceProps {
  /** The price to render. The backend has already converted into the
   *  requested target currency (or returned the source). */
  amount: Money;
  /** Optional source-currency value to render below as a muted hint —
   *  pass when `amount` is the converted price and you also want to
   *  show the original quote. */
  original?: Money;
  /** Optional suffix rendered inline (e.g. "/ day", "/ 3 days"). */
  suffix?: string;
  className?: string;
}

/**
 * Dumb money formatter. Frontend no longer does FX conversion — the backend
 * returns prices already converted into the user's target currency. This
 * component just renders.
 */
export function Price({ amount, original, suffix, className }: PriceProps) {
  const showOriginal = original && original.currency !== amount.currency;
  return (
    <span className={cn('inline-flex flex-col', className)}>
      <span className="inline-flex items-baseline gap-1.5">
        <span className="font-semibold text-zinc-900">{formatMoney(amount)}</span>
        {suffix && <span className="text-xs text-zinc-500">{suffix}</span>}
      </span>
      {showOriginal && (
        <span className="text-xs text-zinc-400">{formatMoney(original)}</span>
      )}
    </span>
  );
}
