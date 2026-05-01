import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/Select';
import { useDisplayCurrency } from './useDisplayCurrency';

/**
 * The set of currencies the ECB publishes daily reference rates for, plus EUR
 * as the base. Stable enough to hardcode — if/when ECB adds a code, this list
 * is updated alongside it.
 *
 * Source: https://www.ecb.europa.eu/stats/policy_and_exchange_rates/euro_reference_exchange_rates/
 */
const SUPPORTED_CURRENCIES = [
  'EUR',
  'AUD',
  'BRL',
  'CAD',
  'CHF',
  'CNY',
  'CZK',
  'DKK',
  'GBP',
  'HKD',
  'HUF',
  'IDR',
  'ILS',
  'INR',
  'ISK',
  'JPY',
  'KRW',
  'MXN',
  'MYR',
  'NOK',
  'NZD',
  'PHP',
  'PLN',
  'RON',
  'SEK',
  'SGD',
  'THB',
  'TRY',
  'USD',
  'ZAR',
] as const;

/**
 * Compact currency picker for the navbar. Selecting a value updates the global
 * display currency, which is sent as ?targetCurrency=... on subsequent backend
 * calls so prices come back already converted.
 */
export function CurrencySelector() {
  const { displayCurrency, setDisplayCurrency } = useDisplayCurrency();

  return (
    <Select value={displayCurrency} onValueChange={setDisplayCurrency}>
      <SelectTrigger className="h-9 w-[5.5rem]" aria-label="Display currency">
        <SelectValue placeholder={displayCurrency} />
      </SelectTrigger>
      <SelectContent>
        {SUPPORTED_CURRENCIES.map((c) => (
          <SelectItem key={c} value={c}>
            {c}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}
