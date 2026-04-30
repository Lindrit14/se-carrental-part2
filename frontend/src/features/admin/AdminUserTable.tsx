import { useState } from 'react';
import { toast } from 'sonner';
import { Trash2, ShieldCheck, Shield, AlertCircle } from 'lucide-react';
import { Badge } from '@/components/ui/Badge';
import { Button } from '@/components/ui/Button';
import { Skeleton } from '@/components/ui/Skeleton';
import { EmptyState } from '@/components/ui/EmptyState';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/Dialog';
import type { AdminUser } from '@/lib/api/authApi';
import type { Role } from '@/domain/user';
import { ApiError } from '@/lib/api/errors';
import { useAuth } from '@/features/auth/useAuth';
import { useAdminDeleteUser, useAdminSetRoles, useAdminUsers } from './useAdmin';

function fmtDate(iso: string): string {
  const d = new Date(iso);
  return Number.isNaN(d.getTime()) ? iso : d.toLocaleString();
}

export function AdminUserTable() {
  const { data, isLoading, isError, refetch } = useAdminUsers({ limit: 100 });
  const setRoles = useAdminSetRoles();
  const deleteUser = useAdminDeleteUser();
  const { session } = useAuth();
  const [pendingDelete, setPendingDelete] = useState<AdminUser | null>(null);

  if (isLoading) {
    return (
      <div className="flex flex-col gap-2">
        {Array.from({ length: 6 }).map((_, i) => (
          <Skeleton key={i} className="h-14 w-full" />
        ))}
      </div>
    );
  }

  if (isError) {
    return (
      <EmptyState
        icon={<AlertCircle className="h-8 w-8" />}
        title="Could not load users"
        description="Something went wrong reaching the auth service."
        action={<Button onClick={() => refetch()}>Try again</Button>}
      />
    );
  }

  if (!data || data.items.length === 0) {
    return (
      <EmptyState
        icon={<Shield className="h-8 w-8" />}
        title="No users yet"
        description="No registered accounts in the system."
      />
    );
  }

  const onToggleAdmin = async (user: AdminUser) => {
    const isAdmin = user.roles.includes('admin');
    const next: Role[] = isAdmin
      ? (user.roles.filter((r) => r !== 'admin') as Role[])
      : ([...new Set([...user.roles, 'admin' as Role])] as Role[]);
    if (next.length === 0) next.push('user');
    try {
      await setRoles.mutateAsync({ userId: user.id, roles: next });
      toast.success(isAdmin ? 'Admin role removed' : 'Admin role granted');
    } catch (err) {
      const msg =
        err instanceof ApiError && err.code === 'last_admin'
          ? 'Cannot remove the last admin.'
          : 'Could not update roles.';
      toast.error(msg);
    }
  };

  const onConfirmDelete = async () => {
    if (!pendingDelete) return;
    try {
      await deleteUser.mutateAsync(pendingDelete.id);
      toast.success('User deleted');
    } catch (err) {
      const msg =
        err instanceof ApiError && err.code === 'last_admin'
          ? 'Cannot delete the last admin.'
          : 'Could not delete the user.';
      toast.error(msg);
    } finally {
      setPendingDelete(null);
    }
  };

  return (
    <>
      <div className="overflow-hidden rounded-lg border border-zinc-200 bg-white">
        <table className="w-full text-sm">
          <thead className="border-b border-zinc-200 bg-zinc-50 text-left text-zinc-500">
            <tr>
              <th className="px-4 py-3 font-medium">Email</th>
              <th className="px-4 py-3 font-medium">Roles</th>
              <th className="px-4 py-3 font-medium">Verified</th>
              <th className="px-4 py-3 font-medium">Created</th>
              <th className="px-4 py-3 font-medium text-right">Actions</th>
            </tr>
          </thead>
          <tbody>
            {data.items.map((u) => {
              const isAdmin = u.roles.includes('admin');
              const isSelf = session?.userId === u.id;
              return (
                <tr key={u.id} className="border-b border-zinc-100 last:border-0">
                  <td className="px-4 py-3 font-medium text-zinc-900">{u.email}</td>
                  <td className="px-4 py-3">
                    <div className="flex flex-wrap gap-1">
                      {u.roles.map((r: string) => (
                        <Badge key={r} variant={r === 'admin' ? 'success' : 'neutral'}>
                          {r}
                        </Badge>
                      ))}
                    </div>
                  </td>
                  <td className="px-4 py-3 text-zinc-600">
                    {u.verified ? (
                      <Badge variant="success">verified</Badge>
                    ) : (
                      <Badge variant="muted">pending</Badge>
                    )}
                  </td>
                  <td className="px-4 py-3 text-zinc-500">{fmtDate(u.createdAt)}</td>
                  <td className="px-4 py-3">
                    <div className="flex items-center justify-end gap-2">
                      <Button
                        variant="secondary"
                        size="sm"
                        disabled={setRoles.isPending || isSelf}
                        onClick={() => onToggleAdmin(u)}
                        title={isSelf ? "Can't change your own role" : undefined}
                      >
                        {isAdmin ? (
                          <>
                            <Shield className="h-4 w-4" /> Revoke admin
                          </>
                        ) : (
                          <>
                            <ShieldCheck className="h-4 w-4" /> Make admin
                          </>
                        )}
                      </Button>
                      <Button
                        variant="ghost"
                        size="sm"
                        disabled={isSelf}
                        onClick={() => setPendingDelete(u)}
                        className="text-red-600 hover:bg-red-50 hover:text-red-700"
                        title={isSelf ? "Can't delete yourself" : undefined}
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </div>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>

      <Dialog
        open={Boolean(pendingDelete)}
        onOpenChange={(open) => !open && setPendingDelete(null)}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete user?</DialogTitle>
            <DialogDescription>
              This permanently deletes <strong>{pendingDelete?.email}</strong>, revokes their
              tokens, and emits a <code>user.deleted</code> event. This cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="secondary" onClick={() => setPendingDelete(null)}>
              Cancel
            </Button>
            <Button
              variant="destructive"
              disabled={deleteUser.isPending}
              onClick={onConfirmDelete}
            >
              {deleteUser.isPending ? 'Deleting…' : 'Delete user'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}
