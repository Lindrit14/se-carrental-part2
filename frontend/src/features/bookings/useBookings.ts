import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { bookingsApi, type CreateBookingInput } from '@/lib/api/bookingsApi';
import { useDisplayCurrency } from '@/features/currency/useDisplayCurrency';

export const bookingsKeys = {
  all: ['bookings'] as const,
  mine: (targetCurrency: string) => [...bookingsKeys.all, 'mine', targetCurrency] as const,
};

export function useMyBookings() {
  const { displayCurrency } = useDisplayCurrency();
  return useQuery({
    queryKey: bookingsKeys.mine(displayCurrency),
    queryFn: () => bookingsApi.listMyBookings(displayCurrency),
  });
}

export function useCreateBooking() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (input: CreateBookingInput) => bookingsApi.createBooking(input),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: bookingsKeys.all });
    },
  });
}

export function useCancelBooking() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => bookingsApi.cancelBooking(id),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: bookingsKeys.all });
    },
  });
}
