import { describe, expect, it } from 'vitest';
import { render, screen } from '@testing-library/react';
import { CurrencyContext } from '@/features/currency/CurrencyContext';
import { Price } from '@/components/ui/Price';

interface CtxOverrides {
  displayCurrency?: string;
  rates?: Map<string, number>;
}

function withCtx(node: React.ReactNode, ctx: CtxOverrides = {}) {
  return render(
    <CurrencyContext.Provider
      value={{
        displayCurrency: ctx.displayCurrency ?? 'EUR',
        setDisplayCurrency: () => {},
        rates: ctx.rates ?? new Map(),
        available: Array.from((ctx.rates ?? new Map()).keys()),
        ratesLoading: false,
        ratesError: false,
        rateDate: '2026-04-29',
      }}
    >
      {node}
    </CurrencyContext.Provider>,
  );
}

describe('Price', () => {
  it('renders source unchanged when display currency matches', () => {
    withCtx(
      <Price source={{ amount: '45.00', currency: 'EUR' }} suffix="/ day" />,
      { displayCurrency: 'EUR' },
    );
    expect(screen.getByText(/45.00/)).toBeInTheDocument();
    expect(screen.getByText('/ day')).toBeInTheDocument();
  });

  it('falls back to source unchanged when rates are not loaded', () => {
    withCtx(
      <Price source={{ amount: '100.00', currency: 'EUR' }} />,
      { displayCurrency: 'TRY', rates: new Map() },
    );
    // Empty rate map => unchanged amount/currency
    expect(screen.getByText(/100.00/)).toBeInTheDocument();
  });

  it('converts EUR -> USD via the rate map', () => {
    withCtx(
      <Price source={{ amount: '100.00', currency: 'EUR' }} />,
      {
        displayCurrency: 'USD',
        rates: new Map<string, number>([['EUR', 1], ['USD', 1.10]]),
      },
    );
    // 100 / 1 * 1.10 = 110.00
    expect(screen.getByText(/110\.00/)).toBeInTheDocument();
  });
});
