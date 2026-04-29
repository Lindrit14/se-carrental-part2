import type { Money } from '@/domain/money';

/** Formats a Money value using Intl.NumberFormat. */
export function formatMoney(money: Money, locale = 'en-US'): string {
  const amount = Number.parseFloat(money.amount);
  if (Number.isNaN(amount)) return `${money.amount} ${money.currency}`;

  try {
    return new Intl.NumberFormat(locale, {
      style: 'currency',
      currency: money.currency,
      maximumFractionDigits: 2,
    }).format(amount);
  } catch {
    // Unknown currency code — fall back to a plain rendering.
    return `${amount.toFixed(2)} ${money.currency}`;
  }
}

/**
 * Convert ``source`` into ``targetCurrency`` using an EUR-pivot rate map (the
 * shape returned by ``GET /api/v1/rates``: ``rate[X] = X per 1 EUR``, with
 * ``rate['EUR'] = 1``). Returns the source unchanged when rates are missing
 * for either side — the UI can still render something meaningful while rates
 * are loading or for an exotic currency we don't support.
 */
export function convertMoney(
  source: Money,
  targetCurrency: string,
  rates: Map<string, number>,
): Money {
  if (source.currency === targetCurrency) return source;
  const rateFrom = rates.get(source.currency);
  const rateTo = rates.get(targetCurrency);
  if (!rateFrom || !rateTo) return source;
  const sourceAmount = Number.parseFloat(source.amount);
  if (!Number.isFinite(sourceAmount)) return source;
  const converted = (sourceAmount / rateFrom) * rateTo;
  return { amount: converted.toFixed(2), currency: targetCurrency };
}
