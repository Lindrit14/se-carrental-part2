import { forwardRef, type InputHTMLAttributes } from 'react';
import { cn } from '@/lib/utils/cn';

export type InputProps = InputHTMLAttributes<HTMLInputElement>;

export const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ className, type = 'text', ...props }, ref) => (
    <input
      ref={ref}
      type={type}
      className={cn(
        'flex h-10 w-full rounded-[var(--radius-control)] border border-zinc-200 bg-white px-3 py-2 text-sm placeholder:text-zinc-400',
        'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-zinc-900/10 focus-visible:border-zinc-300',
        'disabled:cursor-not-allowed disabled:bg-zinc-50 disabled:text-zinc-500',
        'aria-invalid:border-red-300 aria-invalid:focus-visible:ring-red-100',
        className,
      )}
      {...props}
    />
  ),
);
Input.displayName = 'Input';
