import { cn } from '@/lib/utils/cn';
import type { Money } from '@/domain/money';
import { convertMoney, formatMoney } from '@/lib/format/money';
import { useDisplayCurrency } from '@/features/currency/useDisplayCurrency';

interface PriceProps {
  /** The price as quoted in its source currency (e.g. car.dailyRate). */
  source: Money;
  /** Optional suffix rendered inline (e.g. ``/ day``, ``/ 3 days``). */
  suffix?: string;
  className?: string;
  /** Override the display currency (rarely needed — defaults to navbar choice). */
  displayCurrency?: string;
  /** When true, hide the muted source-currency hint. */
  hideOriginal?: boolean;
}

/**
 * Render a Money in the user's currently selected display currency. Falls
 * back to the source currency if rates aren't loaded yet — the price still
 * renders, just unconverted, so the page never has empty cells.
 */
export function Price({
  source,
  suffix,
  className,
  displayCurrency,
}: PriceProps) {
  const { displayCurrency: ctxCurrency, rates } = useDisplayCurrency();
  const target = displayCurrency ?? ctxCurrency;
  const converted = convertMoney(source, target, rates);

  return (
    <span className={cn('inline-flex flex-col', className)}>
      <span className="inline-flex items-baseline gap-1.5">
        <span className="font-semibold text-zinc-900">{formatMoney(converted)}</span>
        {suffix && <span className="text-xs text-zinc-500">{suffix}</span>}
      </span>
    </span>
  );
}
