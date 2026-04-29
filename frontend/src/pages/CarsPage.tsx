import { PageHeader } from '@/components/layout/PageHeader';
import { CarList } from '@/features/cars/CarList';

export function CarsPage() {
  return (
    <div className="flex flex-col gap-8">
      <PageHeader
        title="Available cars"
        description="Daily rates shown in the car's listed currency. We convert your total at booking."
      />
      <CarList />
    </div>
  );
}
