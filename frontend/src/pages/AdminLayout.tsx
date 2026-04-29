import { NavLink, Outlet } from 'react-router-dom';
import { PageHeader } from '@/components/layout/PageHeader';
import { cn } from '@/lib/utils/cn';

const TABS = [
  { to: '/admin/users', label: 'Users' },
  { to: '/admin/cars', label: 'Cars' },
  { to: '/admin/bookings', label: 'Bookings' },
];

export function AdminLayout() {
  return (
    <div className="flex flex-col gap-6">
      <PageHeader
        title="Admin"
        description="Manage users, cars and bookings across the platform."
      />
      <nav className="flex gap-1 border-b border-zinc-200">
        {TABS.map((t) => (
          <NavLink
            key={t.to}
            to={t.to}
            className={({ isActive }) =>
              cn(
                '-mb-px rounded-t-md border-b-2 px-3 py-2 text-sm font-medium transition-colors',
                isActive
                  ? 'border-zinc-900 text-zinc-900'
                  : 'border-transparent text-zinc-500 hover:text-zinc-900',
              )
            }
          >
            {t.label}
          </NavLink>
        ))}
      </nav>
      <Outlet />
    </div>
  );
}
