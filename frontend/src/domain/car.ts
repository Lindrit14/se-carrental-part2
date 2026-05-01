import type { Money } from './money';

export type CarCategory =
  | 'SMALL'
  | 'MEDIUM'
  | 'LARGE'
  | 'SUV'
  | 'PEOPLE_CARRIER'
  | 'PREMIUM';

export const CAR_CATEGORIES: CarCategory[] = [
  'SMALL',
  'MEDIUM',
  'LARGE',
  'SUV',
  'PEOPLE_CARRIER',
  'PREMIUM',
];

export interface Car {
  id: string;
  brand: string;
  model: string;
  licensePlate: string;
  dailyRate: Money;
  /** Server-side conversion of dailyRate into the request's targetCurrency.
   *  Absent when no targetCurrency was sent, when source equals target, or
   *  when the currency-converter service was unavailable for the request. */
  dailyRateConverted?: Money;
  location: string;
  category: CarCategory;
}
