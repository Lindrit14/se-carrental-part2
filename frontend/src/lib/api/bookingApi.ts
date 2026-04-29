import type { Booking, BookingStatus } from '@/domain/booking';
import type { Car } from '@/domain/car';
import { request } from './client';

/* ---------- Wire DTOs (camelCase from the Java backend) ---------- */

interface CarDto {
  id: string;
  brand: string;
  model: string;
  licensePlate: string;
  dailyRateAmount: string | number;
  dailyRateCurrency: string;
}

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

const toCar = (dto: CarDto): Car => ({
  id: dto.id,
  brand: dto.brand,
  model: dto.model,
  licensePlate: dto.licensePlate,
  dailyRate: { amount: String(dto.dailyRateAmount), currency: dto.dailyRateCurrency },
});

const toBooking = (dto: BookingDto): Booking => ({
  id: dto.id,
  carId: dto.carId,
  startDate: dto.startDate,
  endDate: dto.endDate,
  status: dto.status,
  totalSource: { amount: String(dto.totalSourceAmount), currency: dto.totalSourceCurrency },
  totalTarget: { amount: String(dto.totalTargetAmount), currency: dto.totalTargetCurrency },
});

/* ---------- Endpoints ---------- */

export interface CreateBookingInput {
  carId: string;
  startDate: string;
  endDate: string;
  targetCurrency: string;
}

export interface CreateCarInput {
  brand: string;
  model: string;
  licensePlate: string;
  dailyRateAmount: string;
  dailyRateCurrency: string;
}

export const bookingApi = {
  async listCars(): Promise<Car[]> {
    const dtos = await request<CarDto[]>('booking', '/api/v1/cars');
    return dtos.map(toCar);
  },

  async getCar(id: string): Promise<Car> {
    const dto = await request<CarDto>('booking', `/api/v1/cars/${id}`);
    return toCar(dto);
  },

  async createBooking(input: CreateBookingInput): Promise<Booking> {
    const dto = await request<BookingDto>('booking', '/api/v1/bookings', {
      method: 'POST',
      body: input,
    });
    return toBooking(dto);
  },

  async listMyBookings(): Promise<Booking[]> {
    const dtos = await request<BookingDto[]>('booking', '/api/v1/bookings/me');
    return dtos.map(toBooking);
  },

  async cancelBooking(id: string): Promise<void> {
    await request<void>('booking', `/api/v1/bookings/${id}`, { method: 'DELETE' });
  },

  /* ---------- Admin ---------- */

  async createCar(input: CreateCarInput): Promise<Car> {
    const dto = await request<CarDto>('booking', '/api/v1/cars', {
      method: 'POST',
      body: input,
    });
    return toCar(dto);
  },

  async adminListBookings(): Promise<AdminBooking[]> {
    const dtos = await request<AdminBookingDto[]>('booking', '/api/v1/admin/bookings');
    return dtos.map(toAdminBooking);
  },
};
