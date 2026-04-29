import { Link } from 'react-router-dom';
import { ArrowRight, ShieldCheck, Zap, Globe2 } from 'lucide-react';
import { Button } from '@/components/ui/Button';
import { Card, CardContent } from '@/components/ui/Card';
import { CurrencyConverterCard } from '@/features/currency/CurrencyConverterCard';
import { useAuth } from '@/features/auth/useAuth';

const FEATURES = [
  {
    icon: Zap,
    title: 'Fast checkout',
    description: 'Pick a car, choose dates, pay in your currency. Done.',
  },
  {
    icon: Globe2,
    title: 'Live FX rates',
    description: 'Booking totals are converted with daily ECB reference rates.',
  },
  {
    icon: ShieldCheck,
    title: 'Privacy by default',
    description: 'Your account data is never shared between services.',
  },
];

export function HomePage() {
  const { session } = useAuth();

  return (
    <div className="flex flex-col gap-16">
      {/* Hero */}
      <section className="grid grid-cols-1 items-center gap-10 lg:grid-cols-[1.1fr_1fr]">
        <div className="flex flex-col gap-6">
          <span className="inline-flex w-fit items-center gap-2 rounded-full border border-zinc-200 bg-white px-3 py-1 text-xs font-medium text-zinc-700">
            <span className="h-1.5 w-1.5 rounded-full bg-emerald-500" />
            Now serving Munich
          </span>
          <h1 className="text-4xl font-semibold tracking-tight text-zinc-900 sm:text-5xl">
            Rent a car. <br className="hidden sm:inline" />
            Pay in your currency.
          </h1>
          <p className="max-w-prose text-base text-zinc-600">
            A no-nonsense booking experience. Browse the fleet, lock in a price, and we'll
            handle the rate conversion through the European Central Bank's daily feed.
          </p>
          <div className="flex flex-wrap items-center gap-2">
            {session ? (
              <Button asChild size="lg">
                <Link to="/cars">
                  Browse cars
                  <ArrowRight className="h-4 w-4" />
                </Link>
              </Button>
            ) : (
              <>
                <Button asChild size="lg">
                  <Link to="/register">
                    Create account
                    <ArrowRight className="h-4 w-4" />
                  </Link>
                </Button>
                <Button asChild variant="secondary" size="lg">
                  <Link to="/login">I already have an account</Link>
                </Button>
              </>
            )}
          </div>
        </div>

        <CurrencyConverterCard />
      </section>

      {/* Features */}
      <section className="grid grid-cols-1 gap-4 md:grid-cols-3">
        {FEATURES.map((f) => (
          <Card key={f.title}>
            <CardContent className="flex flex-col gap-3 p-6">
              <span className="grid h-9 w-9 place-items-center rounded-md border border-zinc-200 bg-zinc-50 text-zinc-700">
                <f.icon className="h-4 w-4" />
              </span>
              <h3 className="text-base font-semibold tracking-tight text-zinc-900">{f.title}</h3>
              <p className="text-sm text-zinc-500">{f.description}</p>
            </CardContent>
          </Card>
        ))}
      </section>
    </div>
  );
}
