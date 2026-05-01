import { useQuery } from '@tanstack/react-query';
import { carsApi, type SearchCarsInput } from '@/lib/api/carsApi';

export function useCarSearch(input: SearchCarsInput) {
  return useQuery({
    queryKey: ['cars', 'search', input],
    queryFn: () => carsApi.searchCars(input),
    staleTime: 60 * 1000,
  });
}
