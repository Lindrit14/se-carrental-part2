import type { Money } from './money';

export type BookingStatus = 'PENDING' | 'CONFIRMED' | 'CANCELLED';

export interface Booking {
  id: string;
  carId: string;
  startDate: string; // ISO date (YYYY-MM-DD)
  endDate: string;
  status: BookingStatus;
  totalSource: Money;
  totalTarget: Money;
}
