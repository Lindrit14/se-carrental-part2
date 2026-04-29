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
  location: string;
  category: CarCategory;
}
