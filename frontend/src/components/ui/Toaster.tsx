import { Toaster as Sonner } from 'sonner';

/**
 * App-wide toast container. Mounted once near the root of the tree.
 * Style is intentionally muted — no colored backgrounds, just borders.
 */
export function Toaster() {
  return (
    <Sonner
      position="top-right"
      toastOptions={{
        classNames: {
          toast:
            'group flex items-center gap-3 rounded-[var(--radius-card)] border border-zinc-200 bg-white text-sm text-zinc-900 shadow-md',
          title: 'text-sm font-medium',
          description: 'text-sm text-zinc-500',
          actionButton: 'text-zinc-900',
          cancelButton: 'text-zinc-500',
          error: 'border-red-200',
          success: 'border-emerald-200',
        },
      }}
    />
  );
}
