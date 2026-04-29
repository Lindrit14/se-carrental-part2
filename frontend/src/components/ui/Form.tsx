import { type HTMLAttributes, type LabelHTMLAttributes } from 'react';
import { cn } from '@/lib/utils/cn';

/**
 * Tiny composition helpers for react-hook-form layouts. We don't try to
 * abstract react-hook-form itself — these are just consistent visual containers.
 */
export function FormField({ className, ...props }: HTMLAttributes<HTMLDivElement>) {
  return <div className={cn('flex flex-col gap-1.5', className)} {...props} />;
}

export function FormLabel({ className, ...props }: LabelHTMLAttributes<HTMLLabelElement>) {
  return (
    <label
      className={cn('text-sm font-medium text-zinc-700', className)}
      {...props}
    />
  );
}

export function FormHelp({ className, ...props }: HTMLAttributes<HTMLParagraphElement>) {
  return <p className={cn('text-xs text-zinc-500', className)} {...props} />;
}

export function FormError({ className, ...props }: HTMLAttributes<HTMLParagraphElement>) {
  if (!props.children) return null;
  return <p className={cn('text-xs text-red-600', className)} {...props} />;
}
