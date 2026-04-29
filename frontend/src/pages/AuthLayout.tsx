import type { ReactNode } from 'react';
import { Link } from 'react-router-dom';

interface AuthLayoutProps {
  title: string;
  subtitle?: string;
  footer?: ReactNode;
  children: ReactNode;
}

/**
 * Centered card layout used by all auth screens. Visually consistent and
 * keeps the marketing chrome out of the auth flow.
 */
export function AuthLayout({ title, subtitle, footer, children }: AuthLayoutProps) {
  return (
    <div className="flex min-h-dvh flex-col bg-zinc-50">
      <div className="mx-auto flex w-full max-w-md flex-1 flex-col justify-center px-6 py-12">
        <Link to="/" className="mb-8 flex items-center gap-2 self-start font-semibold tracking-tight">
          <span className="grid h-7 w-7 place-items-center rounded-md bg-zinc-900 text-zinc-50 text-sm">
            D
          </span>
          <span>Drive</span>
        </Link>
        <div className="rounded-[var(--radius-card)] border border-zinc-200 bg-white p-8 shadow-sm">
          <div className="mb-6 flex flex-col gap-1">
            <h1 className="text-xl font-semibold tracking-tight text-zinc-900">{title}</h1>
            {subtitle && <p className="text-sm text-zinc-500">{subtitle}</p>}
          </div>
          {children}
        </div>
        {footer && <div className="mt-6 text-center text-sm text-zinc-500">{footer}</div>}
      </div>
    </div>
  );
}
