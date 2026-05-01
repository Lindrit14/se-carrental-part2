import { useQuery } from '@tanstack/react-query';
import { carsApi, type SearchCarsInput } from '@/lib/api/carsApi';
import { useDisplayCurrency } from '@/features/currency/useDisplayCurrency';

export function useCarSearch(input: SearchCarsInput) {
  const { displayCurrency } = useDisplayCurrency();
  const merged: SearchCarsInput = { ...input, targetCurrency: displayCurrency };
  return useQuery({
    queryKey: ['cars', 'search', merged],
    queryFn: () => carsApi.searchCars(merged),
    staleTime: 60 * 1000,
  });
}
