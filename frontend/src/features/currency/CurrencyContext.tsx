import { createContext, useCallback, useEffect, useMemo, useState, type ReactNode } from 'react';
import { useQuery } from '@tanstack/react-query';
import { currencyApi } from '@/lib/api/currencyApi';

/**
 * Global display-currency state. Holds the user's chosen currency (persisted
 * to localStorage) and the latest ECB rate map from the currency-converter
 * service, so any component can render a Money value in the user's currency
 * without doing its own RPC.
 */

const STORAGE_KEY = 'display.currency';
const DEFAULT_CURRENCY = 'EUR';

interface CurrencyContextValue {
  displayCurrency: string;
  setDisplayCurrency: (c: string) => void;
  /** Map keyed by ISO code, including the EUR base (1.0). Empty until loaded. */
  rates: Map<string, number>;
  /** All currencies the user can pick (sorted, includes EUR). */
  available: string[];
  ratesLoading: boolean;
  ratesError: boolean;
  rateDate: string | null;
}

export const CurrencyContext = createContext<CurrencyContextValue | null>(null);

function readStored(): string {
  try {
    const v = localStorage.getItem(STORAGE_KEY);
    if (v && /^[A-Z]{3}$/.test(v)) return v;
  } catch {
    /* SSR / private mode — fall through */
  }
  return DEFAULT_CURRENCY;
}

export function CurrencyProvider({ children }: { children: ReactNode }) {
  const [displayCurrency, setDisplayCurrencyState] = useState<string>(readStored);

  const { data, isLoading, isError } = useQuery({
    queryKey: ['currency', 'rates'],
    queryFn: () => currencyApi.rates(),
    staleTime: 60 * 60 * 1000,
    gcTime: 24 * 60 * 60 * 1000,
  });

  const rates = useMemo<Map<string, number>>(() => {
    const m = new Map<string, number>();
    if (data) {
      m.set(data.base, 1);
      for (const r of data.rates) {
        const n = Number.parseFloat(r.rate);
        if (Number.isFinite(n)) m.set(r.currency, n);
      }
    }
    return m;
  }, [data]);

  const available = useMemo(
    () => Array.from(rates.keys()).sort((a, b) => a.localeCompare(b)),
    [rates],
  );

  const setDisplayCurrency = useCallback((c: string) => {
    setDisplayCurrencyState(c);
    try {
      localStorage.setItem(STORAGE_KEY, c);
    } catch {
      /* ignore */
    }
  }, []);

  // If the persisted choice isn't in the loaded rate set (typo, retired code),
  // fall back to EUR rather than rendering broken prices forever.
  useEffect(() => {
    if (rates.size > 0 && !rates.has(displayCurrency)) {
      setDisplayCurrency(DEFAULT_CURRENCY);
    }
  }, [rates, displayCurrency, setDisplayCurrency]);

  const value = useMemo<CurrencyContextValue>(
    () => ({
      displayCurrency,
      setDisplayCurrency,
      rates,
      available,
      ratesLoading: isLoading,
      ratesError: isError,
      rateDate: data?.rateDate ?? null,
    }),
    [displayCurrency, setDisplayCurrency, rates, available, isLoading, isError, data?.rateDate],
  );

  return <CurrencyContext.Provider value={value}>{children}</CurrencyContext.Provider>;
}
