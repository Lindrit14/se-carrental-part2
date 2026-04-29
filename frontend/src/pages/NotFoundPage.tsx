import { Link } from 'react-router-dom';
import { Button } from '@/components/ui/Button';

export function NotFoundPage() {
  return (
    <div className="flex min-h-dvh items-center justify-center bg-zinc-50 px-6">
      <div className="flex max-w-md flex-col items-center gap-4 text-center">
        <span className="font-mono text-sm text-zinc-500">404</span>
        <h1 className="text-2xl font-semibold tracking-tight text-zinc-900">Page not found</h1>
        <p className="text-sm text-zinc-500">
          The page you're looking for doesn't exist or has been moved.
        </p>
        <Button asChild className="mt-2">
          <Link to="/">Back home</Link>
        </Button>
      </div>
    </div>
  );
}
