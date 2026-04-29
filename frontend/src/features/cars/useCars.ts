import { useQuery } from '@tanstack/react-query';
import { bookingApi } from '@/lib/api/bookingApi';

export const carsKeys = {
  all: ['cars'] as const,
  list: () => [...carsKeys.all, 'list'] as const,
  detail: (id: string) => [...carsKeys.all, 'detail', id] as const,
};

export function useCars() {
  return useQuery({
    queryKey: carsKeys.list(),
    queryFn: () => bookingApi.listCars(),
  });
}

export function useCar(id: string) {
  return useQuery({
    queryKey: carsKeys.detail(id),
    queryFn: () => bookingApi.getCar(id),
    enabled: Boolean(id),
  });
}
