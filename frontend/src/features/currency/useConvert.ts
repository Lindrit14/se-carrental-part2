import { useMutation } from '@tanstack/react-query';
import { currencyApi } from '@/lib/api/currencyApi';

export function useConvert() {
  return useMutation({
    mutationFn: ({ amount, from, to }: { amount: string; from: string; to: string }) =>
      currencyApi.convert(amount, from, to),
  });
}
