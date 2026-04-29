import type { Car } from '@/domain/car';
import { Card, CardContent, CardFooter } from '@/components/ui/Card';
import { Badge } from '@/components/ui/Badge';
import { formatMoney } from '@/lib/format/money';

interface CarCardProps {
  car: Car;
  action?: React.ReactNode;
}

export function CarCard({ car, action }: CarCardProps) {
  return (
    <Card className="flex flex-col">
      <CardContent className="flex flex-1 flex-col gap-4 p-6">
        <div className="flex items-start justify-between gap-2">
          <div>
            <h3 className="text-base font-semibold tracking-tight text-zinc-900">
              {car.brand} {car.model}
            </h3>
            <p className="text-xs text-zinc-500">{car.licensePlate}</p>
          </div>
          <Badge variant="muted">Available</Badge>
        </div>
        <div className="mt-auto">
          <div className="flex items-baseline gap-1.5">
            <span className="text-xl font-semibold text-zinc-900">
              {formatMoney(car.dailyRate)}
            </span>
            <span className="text-xs text-zinc-500">/ day</span>
          </div>
        </div>
      </CardContent>
      {action && <CardFooter className="border-t border-zinc-100 px-6 py-4">{action}</CardFooter>}
    </Card>
  );
}
