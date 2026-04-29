import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { bookingApi, type CreateBookingInput } from '@/lib/api/bookingApi';

export const bookingsKeys = {
  all: ['bookings'] as const,
  mine: () => [...bookingsKeys.all, 'mine'] as const,
};

export function useMyBookings() {
  return useQuery({
    queryKey: bookingsKeys.mine(),
    queryFn: () => bookingApi.listMyBookings(),
  });
}

export function useCreateBooking() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (input: CreateBookingInput) => bookingApi.createBooking(input),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: bookingsKeys.mine() });
    },
  });
}

export function useCancelBooking() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => bookingApi.cancelBooking(id),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: bookingsKeys.mine() });
    },
  });
}
