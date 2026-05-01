import { describe, expect, it } from 'vitest';
import { render, screen } from '@testing-library/react';
import { Price } from '@/components/ui/Price';

describe('Price', () => {
  it('renders the amount in its currency', () => {
    render(<Price amount={{ amount: '45.00', currency: 'EUR' }} suffix="/ day" />);
    expect(screen.getByText(/45\.00/)).toBeInTheDocument();
    expect(screen.getByText('/ day')).toBeInTheDocument();
  });

  it('renders the original below when both currencies differ', () => {
    render(
      <Price
        amount={{ amount: '110.00', currency: 'USD' }}
        original={{ amount: '100.00', currency: 'EUR' }}
      />,
    );
    expect(screen.getByText(/110\.00/)).toBeInTheDocument();
    expect(screen.getByText(/100\.00/)).toBeInTheDocument();
  });

  it('hides the original when currencies match', () => {
    render(
      <Price
        amount={{ amount: '100.00', currency: 'EUR' }}
        original={{ amount: '100.00', currency: 'EUR' }}
      />,
    );
    // Only one rendering of 100.00 — the duplicate "original" line is skipped
    // because the currencies are the same.
    expect(screen.getAllByText(/100\.00/)).toHaveLength(1);
  });
});
