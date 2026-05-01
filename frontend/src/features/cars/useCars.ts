import { useQuery } from '@tanstack/react-query';
import { carsApi } from '@/lib/api/carsApi';

export const carsKeys = {
  all: ['cars'] as const,
  list: () => [...carsKeys.all, 'list'] as const,
  detail: (id: string) => [...carsKeys.all, 'detail', id] as const,
};

export function useCars() {
  return useQuery({
    queryKey: carsKeys.list(),
    queryFn: () => carsApi.listCars(),
  });
}

export function useCar(id: string) {
  return useQuery({
    queryKey: carsKeys.detail(id),
    queryFn: () => carsApi.getCar(id),
    enabled: Boolean(id),
  });
}
