import { useState } from 'react';
import { ArrowRight } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/Select';
import { FormField, FormLabel } from '@/components/ui/Form';
import { useConvert } from './useConvert';

const CURRENCIES = ['EUR', 'USD', 'GBP', 'CHF', 'JPY', 'NOK', 'SEK', 'PLN'];

export function CurrencyConverterCard() {
  const [amount, setAmount] = useState('100');
  const [from, setFrom] = useState('EUR');
  const [to, setTo] = useState('USD');
  const convert = useConvert();

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    convert.mutate({ amount, from, to });
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Quick currency check</CardTitle>
        <CardDescription>Powered by ECB daily reference rates.</CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={onSubmit} className="flex flex-col gap-4">
          <div className="grid grid-cols-1 items-end gap-3 sm:grid-cols-[1fr_auto_1fr]">
            <FormField>
              <FormLabel htmlFor="cc-amount">Amount</FormLabel>
              <div className="grid grid-cols-[1fr_7rem] gap-2">
                <Input
                  id="cc-amount"
                  type="number"
                  inputMode="decimal"
                  min="0"
                  step="0.01"
                  value={amount}
                  onChange={(e) => setAmount(e.target.value)}
                />
                <Select value={from} onValueChange={setFrom}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {CURRENCIES.map((c) => (
                      <SelectItem key={c} value={c}>
                        {c}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </FormField>

            <ArrowRight className="hidden h-5 w-5 self-end text-zinc-300 sm:block" />

            <FormField>
              <FormLabel>Into</FormLabel>
              <Select value={to} onValueChange={setTo}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {CURRENCIES.map((c) => (
                    <SelectItem key={c} value={c}>
                      {c}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </FormField>
          </div>

          <Button type="submit" disabled={convert.isPending} className="self-start">
            {convert.isPending ? 'Converting…' : 'Convert'}
          </Button>

          {convert.data && (
            <div className="rounded-[var(--radius-card)] border border-zinc-200 bg-zinc-50 p-4">
              <p className="text-xs uppercase tracking-wide text-zinc-500">Result</p>
              <p className="mt-1 text-2xl font-semibold tracking-tight text-zinc-900">
                {Number(convert.data.amount).toLocaleString(undefined, {
                  maximumFractionDigits: 2,
                })}{' '}
                {convert.data.currency}
              </p>
              <p className="mt-1 text-xs text-zinc-500">
                Rate date: {convert.data.rateDate}
              </p>
            </div>
          )}

          {convert.isError && (
            <p className="rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-700">
              Conversion failed. Please try again.
            </p>
          )}
        </form>
      </CardContent>
    </Card>
  );
}
