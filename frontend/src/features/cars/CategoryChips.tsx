import { Car as CarIcon, Truck, Bus, Crown, Caravan, CarFront } from 'lucide-react';
import { CAR_CATEGORIES, type CarCategory } from '@/domain/car';
import { cn } from '@/lib/utils/cn';

const META: Record<CarCategory, { label: string; Icon: typeof CarIcon }> = {
  SMALL: { label: 'Small car', Icon: CarFront },
  MEDIUM: { label: 'Medium car', Icon: CarIcon },
  LARGE: { label: 'Large car', Icon: Truck },
  SUV: { label: 'SUVs', Icon: Caravan },
  PEOPLE_CARRIER: { label: 'People carrier', Icon: Bus },
  PREMIUM: { label: 'Premium car', Icon: Crown },
};

interface CategoryChipsProps {
  selected: CarCategory | null;
  onSelect: (c: CarCategory | null) => void;
}

export function CategoryChips({ selected, onSelect }: CategoryChipsProps) {
  return (
    <div className="flex flex-wrap items-center gap-2">
      {CAR_CATEGORIES.map((c) => {
        const { label, Icon } = META[c];
        const active = selected === c;
        return (
          <button
            key={c}
            type="button"
            onClick={() => onSelect(active ? null : c)}
            className={cn(
              'flex flex-col items-center gap-1 rounded-lg border px-4 py-3 text-xs transition-colors',
              active
                ? 'border-zinc-900 bg-zinc-900 text-white'
                : 'border-zinc-200 bg-white text-zinc-700 hover:bg-zinc-50',
            )}
          >
            <Icon className="h-5 w-5" />
            <span>{label}</span>
          </button>
        );
      })}
      {selected && (
        <button
          type="button"
          onClick={() => onSelect(null)}
          className="flex flex-col justify-center items-center gap-1 rounded-lg border border-zinc-200 px-4 py-3 text-xs text-zinc-600 transition-colors hover:bg-red-50"
        >
          
          <span>Clear</span>
        </button>
      )}
    </div>
  );
}
