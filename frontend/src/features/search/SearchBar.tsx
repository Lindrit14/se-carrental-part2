import { useEffect, useId, useRef } from 'react';
import { useForm } from 'react-hook-form';
import { useNavigate } from 'react-router-dom';
import { Search } from 'lucide-react';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { FormError, FormLabel } from '@/components/ui/Form';
import type { CarCategory } from '@/domain/car';
import { useGoogleMaps } from '@/features/maps/GoogleMapsProvider';
import { cn } from '@/lib/utils/cn';

interface SearchValues {
  location: string;
  from: string;
  to: string;
}

interface SearchBarProps {
  /** Initial values (from URL params on the results page). */
  defaultValues?: Partial<SearchValues>;
  /** Visual mode — hero (large, on homepage) or compact (re-edit on results page). */
  variant?: 'hero' | 'compact';
  /** Optional category from URL state, preserved when re-submitting. */
  category?: CarCategory | null;
  className?: string;
}

const todayIso = (offsetDays = 0): string => {
  const d = new Date();
  d.setDate(d.getDate() + offsetDays);
  return d.toISOString().slice(0, 10);
};

export function SearchBar({ defaultValues, variant = 'hero', category, className }: SearchBarProps) {
  const navigate = useNavigate();
  const id = useId();
  const { isLoaded: mapsLoaded } = useGoogleMaps();

  const { register, handleSubmit, setValue, formState: { errors } } = useForm<SearchValues>({
    defaultValues: {
      location: defaultValues?.location ?? '',
      from: defaultValues?.from ?? todayIso(),
      to: defaultValues?.to ?? todayIso(3),
    },
  });

  // RHF's register returns a ref-setter; we keep our own ref to the same DOM
  // node so we can hand it to Google Places Autocomplete. Both refs end up
  // pointing at the same input.
  const locationRef = useRef<HTMLInputElement | null>(null);
  const locationReg = register('location', { required: 'A pick-up location is required' });

  useEffect(() => {
    if (!mapsLoaded) return;
    const input = locationRef.current;
    if (!input || !window.google?.maps?.places) return;

    const autocomplete = new window.google.maps.places.Autocomplete(input, {
      fields: ['formatted_address', 'name', 'place_id'],
      types: ['geocode', 'establishment'],
    });

    const listener = autocomplete.addListener('place_changed', () => {
      const place = autocomplete.getPlace();
      const value = place.formatted_address || place.name || input.value;
      setValue('location', value, { shouldValidate: true, shouldDirty: true });
    });

    return () => {
      listener.remove();
      // Google injects a .pac-container into <body>; React doesn't manage it.
      // It's harmless to leave across re-mounts — Google reuses it.
    };
  }, [mapsLoaded, setValue]);

  const onSubmit = handleSubmit((values) => {
    const params = new URLSearchParams();
    if (values.location) params.set('location', values.location);
    if (values.from) params.set('from', values.from);
    if (values.to) params.set('to', values.to);
    if (category) params.set('category', category);
    navigate(`/cars?${params.toString()}`);
  });

  const isHero = variant === 'hero';

  return (
    <form
      onSubmit={onSubmit}
      noValidate
      className={cn(
        'rounded-xl bg-white shadow-lg ring-1 ring-amber-300',
        isHero ? 'p-2' : 'p-1.5',
        className,
      )}
    >
      <div className="grid grid-cols-1 items-end gap-2 md:grid-cols-[2fr_1fr_1fr_auto]">
        <div className="flex flex-col gap-1 px-3 py-1">
          <FormLabel htmlFor={`${id}-location`} className="text-[10px] uppercase tracking-wide text-zinc-500">
            Pick-up location
          </FormLabel>
          <input
            id={`${id}-location`}
            placeholder="City, airport, or station"
            autoComplete="off"
            className="border-0 bg-transparent p-0 text-sm font-medium text-zinc-900 placeholder:text-zinc-400 focus:outline-none"
            {...locationReg}
            ref={(el) => {
              locationReg.ref(el);
              locationRef.current = el;
            }}
          />
          <FormError>{errors.location?.message}</FormError>
        </div>

        <div className="flex flex-col gap-1 border-t border-zinc-100 px-3 py-1 md:border-l md:border-t-0">
          <FormLabel htmlFor={`${id}-from`} className="text-[10px] uppercase tracking-wide text-zinc-500">
            Pick-up date
          </FormLabel>
          <Input
            id={`${id}-from`}
            type="date"
            className="text-black border-0 bg-transparent px-0 text-sm font-medium focus-visible:ring-0"
            {...register('from', { required: true })}
          />
        </div>

        <div className="flex flex-col gap-1 border-t border-zinc-100 px-3 py-1 md:border-l md:border-t-0">
          <FormLabel htmlFor={`${id}-to`} className="text-[10px] uppercase tracking-wide text-zinc-500">
            Drop-off date
          </FormLabel>
          <Input
            id={`${id}-to`}
            type="date"
            className="text-black border-0 bg-transparent px-0 text-sm font-medium focus-visible:ring-0"
            {...register('to', { required: true })}
          />
        </div>

        <Button type="submit" size="lg" className="h-12 gap-2 bg-emerald-600 hover:bg-emerald-700">
          <Search className="h-4 w-4" />
          Search
        </Button>
      </div>
    </form>
  );
}
