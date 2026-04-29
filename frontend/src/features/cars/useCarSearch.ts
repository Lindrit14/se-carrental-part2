import { useQuery } from '@tanstack/react-query';
import { bookingApi, type SearchCarsInput } from '@/lib/api/bookingApi';

export function useCarSearch(input: SearchCarsInput) {
  return useQuery({
    queryKey: ['cars', 'search', input],
    queryFn: () => bookingApi.searchCars(input),
    staleTime: 60 * 1000,
  });
}
