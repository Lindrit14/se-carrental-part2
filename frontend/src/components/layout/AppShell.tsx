import type { ReactNode } from 'react';
import { Navbar } from './Navbar';

export function AppShell({ children }: { children: ReactNode }) {
  return (
    <div className="flex min-h-dvh flex-col bg-zinc-50">
      <Navbar />
      <main className="flex-1">
        <div className="mx-auto w-full max-w-7xl px-6 py-10 lg:px-8">{children}</div>
      </main>
      <footer className="border-t border-zinc-200 bg-white">
        <div className="mx-auto flex max-w-7xl items-center justify-between px-6 py-6 text-xs text-zinc-500 lg:px-8">
          <span>© {new Date().getFullYear()} Drive</span>
          <span>v0.1.0</span>
        </div>
      </footer>
    </div>
  );
}
