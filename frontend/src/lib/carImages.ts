import type { Car } from '@/domain/car';

const DEFAULT_CAR_IMAGE = '/cars/Default-Car.jpg';

const carImages: Record<string, string> = {
  'audi-a6': '/cars/Audi-A6.jpg',
  'bmw-m4': '/cars/BMW-M4.jpg',
  'bmw-x3': '/cars/BMW-X3.jpg',
  'fiat-500': '/cars/Fiat-500.jpg',
  'ford-galaxy': '/cars/Ford-Galaxy.jpg',
  'ford-toureno': '/cars/Ford-Toureno.jpg',
  'mercedes-e': '/cars/Mercedes-E.jpg',
  'opel-corsa': '/cars/Opel-corsa.jpg',
  'opel-crossland': '/cars/Opel-Crossland.jpg',
  'renault-megane': '/cars/Renault-Megane.jpg',
  'suzuki-swift': '/cars/Suzuki-Swift.jpg',
  'toyota-corolla': '/cars/Toyota-Corolla.jpg',
  'volkswagen-touran': '/cars/Volkswagen-Touran.jpg',
  'vw-polo': '/cars/VW-Polo.jpg',
  'vw-up': '/cars/VW-up.jpg',

  'volkswagen-polo': '/cars/VW-Polo.jpg',
  'volkswagen-up': '/cars/VW-up.jpg',
  'vw-touran': '/cars/Volkswagen-Touran.jpg',
};

function normalize(value: string): string {
  return value.trim().toLowerCase().replace(/\s+/g, '-');
}

export function getCarImage(car: Pick<Car, 'brand' | 'model'>): string {
  const key = `${normalize(car.brand)}-${normalize(car.model)}`;
  return carImages[key] ?? DEFAULT_CAR_IMAGE;
}
