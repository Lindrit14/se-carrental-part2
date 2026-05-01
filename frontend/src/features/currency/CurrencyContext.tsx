import { createContext, useCallback, useMemo, useState, type ReactNode } from 'react';

/**
 * Holds the user's chosen display currency (persisted to localStorage). The
 * backend now does FX conversion server-side — this context no longer fetches
 * or stores rates. Components send the chosen currency to the backend and
 * render whatever it returns.
 */

const STORAGE_KEY = 'display.currency';
const DEFAULT_CURRENCY = 'EUR';

interface CurrencyContextValue {
  displayCurrency: string;
  setDisplayCurrency: (c: string) => void;
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

  const setDisplayCurrency = useCallback((c: string) => {
    setDisplayCurrencyState(c);
    try {
      localStorage.setItem(STORAGE_KEY, c);
    } catch {
      /* ignore */
    }
  }, []);

  const value = useMemo<CurrencyContextValue>(
    () => ({ displayCurrency, setDisplayCurrency }),
    [displayCurrency, setDisplayCurrency],
  );

  return <CurrencyContext.Provider value={value}>{children}</CurrencyContext.Provider>;
}
