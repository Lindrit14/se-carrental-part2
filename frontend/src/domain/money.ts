/** Value object: an amount in a given currency. */
export interface Money {
  amount: string; // BigDecimal string from the backend, parse only when needed
  currency: string; // ISO-4217
}
