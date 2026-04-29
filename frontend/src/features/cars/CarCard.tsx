import { MapPin } from 'lucide-react';
import type { Car, CarCategory } from '@/domain/car';
import { Card, CardContent, CardFooter } from '@/components/ui/Card';
import { Badge } from '@/components/ui/Badge';
import { Price } from '@/components/ui/Price';

const CATEGORY_LABEL: Record<CarCategory, string> = {
  SMALL: 'Small car',
  MEDIUM: 'Medium car',
  LARGE: 'Large car',
  SUV: 'SUV',
  PEOPLE_CARRIER: 'People carrier',
  PREMIUM: 'Premium car',
};

interface CarCardProps {
  car: Car;
  action?: React.ReactNode;
  /** When true, also render the license plate (admin-only). */
  showLicensePlate?: boolean;
}

export function CarCard({ car, action, showLicensePlate }: CarCardProps) {
  return (
    <Card className="flex flex-col">
      <CardContent className="flex flex-1 flex-col gap-4 p-6">
        <div className="flex items-start justify-between gap-2">
          <div>
            <h3 className="text-base font-semibold tracking-tight text-zinc-900">
              {car.brand} {car.model}
            </h3>
            {showLicensePlate && (
              <p className="font-mono text-xs text-zinc-500">{car.licensePlate}</p>
            )}
            {car.location && (
              <p className="mt-1 inline-flex items-center gap-1 text-xs text-zinc-500">
                <MapPin className="h-3 w-3" />
                {car.location}
              </p>
            )}
          </div>
          <Badge variant="muted">{CATEGORY_LABEL[car.category]}</Badge>
        </div>
        <div className="mt-auto">
          <Price source={car.dailyRate} suffix="/ day" />
        </div>
      </CardContent>
      {action && <CardFooter className="border-t border-zinc-100 px-6 py-4">{action}</CardFooter>}
    </Card>
  );
}
