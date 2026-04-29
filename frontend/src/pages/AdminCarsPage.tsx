import { useState } from 'react';
import { Plus } from 'lucide-react';
import { Button } from '@/components/ui/Button';
import { CarList } from '@/features/cars/CarList';
import { CreateCarDialog } from '@/features/admin/CreateCarDialog';

export function AdminCarsPage() {
  const [open, setOpen] = useState(false);
  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-center justify-between">
        <p className="text-sm text-zinc-500">
          The full fleet. Newly added cars become bookable immediately.
        </p>
        <Button onClick={() => setOpen(true)}>
          <Plus className="h-4 w-4" /> Add car
        </Button>
      </div>
      <CarList />
      <CreateCarDialog open={open} onOpenChange={setOpen} />
    </div>
  );
}
