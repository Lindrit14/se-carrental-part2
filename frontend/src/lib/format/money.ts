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
