import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/Select';
import { useDisplayCurrency } from './useDisplayCurrency';

/**
 * Compact currency picker for the navbar. Lists every currency reported by
 * the rates feed (plus the EUR base). Selecting a value updates the global
 * display currency, so every <Price> on the page re-renders.
 */
export function CurrencySelector() {
  const { displayCurrency, setDisplayCurrency, available, ratesLoading } = useDisplayCurrency();

  return (
    <Select value={displayCurrency} onValueChange={setDisplayCurrency}>
      <SelectTrigger
        className="h-9 w-[5.5rem]"
        aria-label="Display currency"
        disabled={ratesLoading && available.length === 0}
      >
        <SelectValue placeholder={displayCurrency} />
      </SelectTrigger>
      <SelectContent>
        {(available.length > 0 ? available : [displayCurrency]).map((c) => (
          <SelectItem key={c} value={c}>
            {c}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}
