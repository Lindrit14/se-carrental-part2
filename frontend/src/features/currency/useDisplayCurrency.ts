import { useContext } from 'react';
import { CurrencyContext } from './CurrencyContext';

export function useDisplayCurrency() {
  const ctx = useContext(CurrencyContext);
  if (!ctx) {
    throw new Error('useDisplayCurrency must be used within a CurrencyProvider');
  }
  return ctx;
}
