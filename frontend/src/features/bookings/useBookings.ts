import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { bookingsApi, type CreateBookingInput } from '@/lib/api/bookingsApi';

export const bookingsKeys = {
  all: ['bookings'] as const,
  mine: () => [...bookingsKeys.all, 'mine'] as const,
};

export function useMyBookings() {
  return useQuery({
    queryKey: bookingsKeys.mine(),
    queryFn: () => bookingsApi.listMyBookings(),
  });
}

export function useCreateBooking() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (input: CreateBookingInput) => bookingsApi.createBooking(input),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: bookingsKeys.mine() });
    },
  });
}

export function useCancelBooking() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => bookingsApi.cancelBooking(id),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: bookingsKeys.mine() });
    },
  });
}
