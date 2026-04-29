import { Link, NavLink, useNavigate } from 'react-router-dom';
import { LogOut, User as UserIcon, Menu, ShieldCheck } from 'lucide-react';
import { Button } from '@/components/ui/Button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/DropdownMenu';
import { useAuth } from '@/features/auth/useAuth';
import { CurrencySelector } from '@/features/currency/CurrencySelector';
import { cn } from '@/lib/utils/cn';

const PUBLIC_LINKS = [{ to: '/cars', label: 'Cars' }];
const AUTHED_LINKS = [{ to: '/bookings', label: 'My Bookings' }];

export function Navbar() {
  const { session, logout } = useAuth();
  const navigate = useNavigate();
  const isAdmin = session?.roles.includes('admin') ?? false;

  return (
    <header className="sticky top-0 z-30 border-b border-zinc-200 bg-white/80 backdrop-blur-md">
      <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-6 lg:px-8">
        <div className="flex items-center gap-8">
          <Link to="/" className="flex items-center gap-2 font-semibold tracking-tight">
            <span className="grid h-7 w-7 place-items-center rounded-md bg-zinc-900 text-zinc-50 text-sm">
              D
            </span>
            <span className="text-zinc-900">Drive</span>
          </Link>
          <nav className="hidden gap-1 md:flex">
            {[...PUBLIC_LINKS, ...(session ? AUTHED_LINKS : [])].map((l) => (
              <NavLink
                key={l.to}
                to={l.to}
                className={({ isActive }) =>
                  cn(
                    'rounded-md px-3 py-1.5 text-sm font-medium transition-colors',
                    isActive
                      ? 'bg-zinc-100 text-zinc-900'
                      : 'text-zinc-600 hover:bg-zinc-100 hover:text-zinc-900',
                  )
                }
              >
                {l.label}
              </NavLink>
            ))}
            {isAdmin && (
              <NavLink
                to="/admin"
                className={({ isActive }) =>
                  cn(
                    'flex items-center gap-1.5 rounded-md px-3 py-1.5 text-sm font-medium transition-colors',
                    isActive
                      ? 'bg-zinc-100 text-zinc-900'
                      : 'text-zinc-600 hover:bg-zinc-100 hover:text-zinc-900',
                  )
                }
              >
                <ShieldCheck className="h-3.5 w-3.5" /> Admin
              </NavLink>
            )}
          </nav>
        </div>

        <div className="flex items-center gap-2">
          <CurrencySelector />
          {session ? (
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="secondary" size="sm" className="gap-2">
                  <UserIcon className="h-4 w-4" />
                  <span className="hidden sm:inline">{session.email ?? 'Account'}</span>
                  <Menu className="h-4 w-4 sm:hidden" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuLabel>Signed in</DropdownMenuLabel>
                <DropdownMenuItem onSelect={() => navigate('/profile')}>
                  Profile
                </DropdownMenuItem>
                <DropdownMenuItem onSelect={() => navigate('/bookings')}>
                  My Bookings
                </DropdownMenuItem>
                {isAdmin && (
                  <DropdownMenuItem onSelect={() => navigate('/admin')}>
                    <ShieldCheck className="h-4 w-4" />
                    Admin
                  </DropdownMenuItem>
                )}
                <DropdownMenuSeparator />
                <DropdownMenuItem
                  onSelect={() => {
                    void logout();
                  }}
                  className="text-red-600 data-[highlighted]:bg-red-50 data-[highlighted]:text-red-700"
                >
                  <LogOut className="h-4 w-4" />
                  Log out
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          ) : (
            <>
              <Button asChild variant="ghost" size="sm">
                <Link to="/login">Log in</Link>
              </Button>
              <Button asChild size="sm">
                <Link to="/register">Sign up</Link>
              </Button>
            </>
          )}
        </div>
      </div>
    </header>
  );
}
