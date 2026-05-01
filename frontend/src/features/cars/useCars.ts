import { useQuery } from '@tanstack/react-query';
import { carsApi } from '@/lib/api/carsApi';
import { useDisplayCurrency } from '@/features/currency/useDisplayCurrency';

export const carsKeys = {
  all: ['cars'] as const,
  list: (targetCurrency: string) => [...carsKeys.all, 'list', targetCurrency] as const,
  detail: (id: string, targetCurrency: string) =>
    [...carsKeys.all, 'detail', id, targetCurrency] as const,
};

export function useCars() {
  const { displayCurrency } = useDisplayCurrency();
  return useQuery({
    queryKey: carsKeys.list(displayCurrency),
    queryFn: () => carsApi.listCars(displayCurrency),
  });
}

export function useCar(id: string) {
  const { displayCurrency } = useDisplayCurrency();
  return useQuery({
    queryKey: carsKeys.detail(id, displayCurrency),
    queryFn: () => carsApi.getCar(id, displayCurrency),
    enabled: Boolean(id),
  });
}
