import type { Money } from './money';

export interface Car {
  id: string;
  brand: string;
  model: string;
  licensePlate: string;
  dailyRate: Money;
}
