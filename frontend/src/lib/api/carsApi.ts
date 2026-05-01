import type { Car, CarCategory } from '@/domain/car';
import { request } from './client';

/* ---------- Wire DTOs (camelCase from the Java backend) ---------- */

interface CarDto {
  id: string;
  brand: string;
  model: string;
  licensePlate: string;
  dailyRateAmount: string | number;
  dailyRateCurrency: string;
  location: string;
  category: CarCategory;
}

const toCar = (dto: CarDto): Car => ({
  id: dto.id,
  brand: dto.brand,
  model: dto.model,
  licensePlate: dto.licensePlate,
  dailyRate: { amount: String(dto.dailyRateAmount), currency: dto.dailyRateCurrency },
  location: dto.location,
  category: dto.category,
});

/* ---------- Endpoints ---------- */

export interface CreateCarInput {
  brand: string;
  model: string;
  licensePlate: string;
  dailyRateAmount: string;
  dailyRateCurrency: string;
  location: string;
  category: CarCategory;
}

export interface SearchCarsInput {
  location?: string;
  category?: CarCategory;
  from?: string;
  to?: string;
}

export const carsApi = {
  async listCars(): Promise<Car[]> {
    const dtos = await request<CarDto[]>('cars', '/api/v1/cars', { auth: false });
    return dtos.map(toCar);
  },

  async searchCars(input: SearchCarsInput): Promise<Car[]> {
    const params = new URLSearchParams();
    if (input.location) params.set('location', input.location);
    if (input.category) params.set('category', input.category);
    if (input.from) params.set('from', input.from);
    if (input.to) params.set('to', input.to);
    const query = params.toString();
    const path = `/api/v1/cars/search${query ? `?${query}` : ''}`;
    const dtos = await request<CarDto[]>('cars', path);
    return dtos.map(toCar);
  },

  async getCar(id: string): Promise<Car> {
    const dto = await request<CarDto>('cars', `/api/v1/cars/${id}`, { auth: false });
    return toCar(dto);
  },

  /* ---------- Admin ---------- */

  async createCar(input: CreateCarInput): Promise<Car> {
    const dto = await request<CarDto>('cars', '/api/v1/cars', {
      method: 'POST',
      body: input,
    });
    return toCar(dto);
  },

  async updateCar(id: string, input: CreateCarInput): Promise<Car> {
    const dto = await request<CarDto>('cars', `/api/v1/cars/${id}`, {
      method: 'PUT',
      body: input,
    });
    return toCar(dto);
  },

  async deleteCar(id: string): Promise<void> {
    await request<void>('cars', `/api/v1/cars/${id}`, { method: 'DELETE' });
  },
};
