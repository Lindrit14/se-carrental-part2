import { describe, expect, it } from 'vitest';
import { convertMoney } from '@/lib/format/money';

describe('convertMoney', () => {
  it('returns source unchanged when source = target', () => {
    const m = { amount: '50.00', currency: 'EUR' };
    expect(convertMoney(m, 'EUR', new Map([['EUR', 1]]))).toBe(m);
  });

  it('returns source unchanged when rate is missing', () => {
    const m = { amount: '50.00', currency: 'EUR' };
    expect(convertMoney(m, 'XYZ', new Map([['EUR', 1]]))).toBe(m);
  });

  it('converts EUR -> USD via the rate map', () => {
    const result = convertMoney(
      { amount: '100.00', currency: 'EUR' },
      'USD',
      new Map([['EUR', 1], ['USD', 1.10]]),
    );
    expect(result).toEqual({ amount: '110.00', currency: 'USD' });
  });

  it('converts USD -> GBP using the EUR pivot', () => {
    // amount in EUR = 100 / 1.10 = 90.909...
    // amount in GBP = 90.909... * 0.85 = 77.272...
    const result = convertMoney(
      { amount: '100.00', currency: 'USD' },
      'GBP',
      new Map([['EUR', 1], ['USD', 1.10], ['GBP', 0.85]]),
    );
    expect(result.currency).toBe('GBP');
    const n = Number.parseFloat(result.amount);
    expect(n).toBeGreaterThan(77);
    expect(n).toBeLessThan(78);
  });
});
