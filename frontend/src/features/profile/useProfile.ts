import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { authApi } from '@/lib/api/authApi';

export const profileKeys = {
  all: ['profile'] as const,
  me: () => [...profileKeys.all, 'me'] as const,
};

export function useMyProfile() {
  return useQuery({
    queryKey: profileKeys.me(),
    queryFn: () => authApi.me(),
  });
}

export function useUpdateProfile() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (input: { email?: string }) => authApi.updateMe(input),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: profileKeys.me() });
    },
  });
}

export function useDeleteAccount() {
  return useMutation({
    mutationFn: () => authApi.deleteMe(),
  });
}
