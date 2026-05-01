import type { Booking, BookingStatus } from '@/domain/booking';
import { request } from './client';

/* ---------- Wire DTOs (camelCase from the Java backend) ---------- */

interface BookingDto {
  id: string;
  carId: string;
  startDate: string;
  endDate: string;
  status: BookingStatus;
  totalSourceAmount: string | number;
  totalSourceCurrency: string;
  totalTargetAmount: string | number;
  totalTargetCurrency: string;
}

interface AdminBookingDto extends BookingDto {
  customerId: string;
  createdAt: string;
  updatedAt: string;
}

export interface AdminBooking extends Booking {
  customerId: string;
  createdAt: string;
  updatedAt: string;
}

const toBooking = (dto: BookingDto): Booking => ({
  id: dto.id,
  carId: dto.carId,
  startDate: dto.startDate,
  endDate: dto.endDate,
  status: dto.status,
  totalSource: { amount: String(dto.totalSourceAmount), currency: dto.totalSourceCurrency },
  totalTarget: { amount: String(dto.totalTargetAmount), currency: dto.totalTargetCurrency },
});

const toAdminBooking = (dto: AdminBookingDto): AdminBooking => ({
  id: dto.id,
  carId: dto.carId,
  startDate: dto.startDate,
  endDate: dto.endDate,
  status: dto.status,
  totalSource: { amount: String(dto.totalSourceAmount), currency: dto.totalSourceCurrency },
  totalTarget: { amount: String(dto.totalTargetAmount), currency: dto.totalTargetCurrency },
  customerId: dto.customerId,
  createdAt: dto.createdAt,
  updatedAt: dto.updatedAt,
});

/* ---------- Endpoints ---------- */

export interface CreateBookingInput {
  carId: string;
  startDate: string;
  endDate: string;
  targetCurrency: string;
}

export const bookingsApi = {
  async createBooking(input: CreateBookingInput): Promise<Booking> {
    const dto = await request<BookingDto>('bookings', '/api/v1/bookings', {
      method: 'POST',
      body: input,
    });
    return toBooking(dto);
  },

  async listMyBookings(): Promise<Booking[]> {
    const dtos = await request<BookingDto[]>('bookings', '/api/v1/bookings/me');
    return dtos.map(toBooking);
  },

  async cancelBooking(id: string): Promise<void> {
    await request<void>('bookings', `/api/v1/bookings/${id}`, { method: 'DELETE' });
  },

  /* ---------- Admin ---------- */

  async adminListBookings(): Promise<AdminBooking[]> {
    const dtos = await request<AdminBookingDto[]>('bookings', '/api/v1/admin/bookings');
    return dtos.map(toAdminBooking);
  },
};
