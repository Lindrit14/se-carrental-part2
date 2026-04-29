import { SearchBar } from '@/features/search/SearchBar';

export function HomePage() {
  return (
    <div className="flex flex-col gap-16">
      <section className="-mx-6 -mt-8 bg-gradient-to-b from-blue-700 to-blue-500 px-6 pt-16 pb-24 text-white lg:-mx-8 lg:px-8">
        <div className="mx-auto max-w-5xl">
          <h1 className="text-4xl font-bold tracking-tight sm:text-5xl">
            Car Rental – Search, Compare &amp; Save
          </h1>
          <ul className="mt-4 flex flex-wrap items-center gap-x-6 gap-y-2 text-sm text-white/90">
          </ul>
          <div className="mt-8">
            <SearchBar />
          </div>
        </div>
      </section>
    </div>
  );
}
