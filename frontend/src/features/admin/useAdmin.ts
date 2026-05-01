import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { authApi, type AdminUser } from '@/lib/api/authApi';
import { bookingsApi } from '@/lib/api/bookingsApi';
import { carsApi, type CreateCarInput } from '@/lib/api/carsApi';
import type { Role } from '@/domain/user';
import { carsKeys } from '@/features/cars/useCars';

export const adminKeys = {
  all: ['admin'] as const,
  users: () => [...adminKeys.all, 'users'] as const,
  bookings: () => [...adminKeys.all, 'bookings'] as const,
};

export function useAdminUsers(opts?: { limit?: number; offset?: number }) {
  return useQuery({
    queryKey: [...adminKeys.users(), opts?.limit ?? 50, opts?.offset ?? 0],
    queryFn: () => authApi.adminListUsers(opts),
  });
}

export function useAdminSetRoles() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ userId, roles }: { userId: string; roles: Role[] }) =>
      authApi.adminSetRoles(userId, roles),
    onSuccess: (updated: AdminUser) => {
      void qc.invalidateQueries({ queryKey: adminKeys.users() });
      return updated;
    },
  });
}

export function useAdminDeleteUser() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (userId: string) => authApi.adminDeleteUser(userId),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: adminKeys.users() });
    },
  });
}

export function useAdminBookings() {
  return useQuery({
    queryKey: adminKeys.bookings(),
    queryFn: () => bookingsApi.adminListBookings(),
  });
}

export function useCreateCar() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (input: CreateCarInput) => carsApi.createCar(input),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: carsKeys.list() });
    },
  });
}

export function useUpdateCar() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, input }: { id: string; input: CreateCarInput }) =>
      carsApi.updateCar(id, input),
    onSuccess: (_data, { id }) => {
      void qc.invalidateQueries({ queryKey: carsKeys.list() });
      void qc.invalidateQueries({ queryKey: carsKeys.detail(id) });
    },
  });
}

export function useDeleteCar() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => carsApi.deleteCar(id),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: carsKeys.list() });
    },
  });
}
