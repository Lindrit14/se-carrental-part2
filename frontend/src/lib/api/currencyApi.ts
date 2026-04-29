import { request } from './client';

export interface ConvertResult {
  amount: string;
  currency: string;
  sourceAmount: string;
  sourceCurrency: string;
  rateDate: string;
}

interface ConvertResponseDto {
  amount: string;
  currency: string;
  source_amount: string;
  source_currency: string;
  rate_date: string;
}

interface RateEntryDto {
  currency: string;
  rate: string | number;
}

interface RatesResponseDto {
  base: string;
  rate_date: string;
  rates: RateEntryDto[];
}

export interface RateSet {
  base: string;
  rateDate: string;
  rates: { currency: string; rate: string }[];
}

export const currencyApi = {
  async convert(amount: string, from: string, to: string): Promise<ConvertResult> {
    const dto = await request<ConvertResponseDto>(
      'currency',
      '/api/v1/convert',
      { method: 'POST', body: { amount, from, to }, auth: false },
    );
    return {
      amount: dto.amount,
      currency: dto.currency,
      sourceAmount: dto.source_amount,
      sourceCurrency: dto.source_currency,
      rateDate: dto.rate_date,
    };
  },

  async rates(): Promise<RateSet> {
    const dto = await request<RatesResponseDto>('currency', '/api/v1/rates', { auth: false });
    return {
      base: dto.base,
      rateDate: dto.rate_date,
      rates: dto.rates.map((r) => ({ currency: r.currency, rate: String(r.rate) })),
    };
  },
};
