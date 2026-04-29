import { useSearchParams } from 'react-router-dom';
import { CAR_CATEGORIES, type CarCategory } from '@/domain/car';
import { SearchBar } from '@/features/search/SearchBar';
import { CategoryChips } from '@/features/cars/CategoryChips';
import { CarResultsList } from '@/features/cars/CarResultsList';
import { useCarSearch } from '@/features/cars/useCarSearch';

function readCategory(raw: string | null): CarCategory | null {
  if (!raw) return null;
  return (CAR_CATEGORIES as string[]).includes(raw) ? (raw as CarCategory) : null;
}

export function CarSearchResultsPage() {
  const [params, setParams] = useSearchParams();
  const location = params.get('location') ?? '';
  const from = params.get('from') ?? '';
  const to = params.get('to') ?? '';
  const category = readCategory(params.get('category'));

  const search = useCarSearch({
    location: location || undefined,
    category: category ?? undefined,
    from: from || undefined,
    to: to || undefined,
  });

  const setCategory = (c: CarCategory | null) => {
    const next = new URLSearchParams(params);
    if (c) next.set('category', c);
    else next.delete('category');
    setParams(next, { replace: true });
  };

  return (
    <div className="flex flex-col gap-8">
      <SearchBar
        variant="compact"
        defaultValues={{ location, from, to }}
        category={category}
      />
      <CategoryChips selected={category} onSelect={setCategory} />
      <CarResultsList
        cars={search.data}
        isLoading={search.isLoading}
        isError={search.isError}
        refetch={() => void search.refetch()}
        bookingDates={from && to ? { from, to } : undefined}
      />
    </div>
  );
}
