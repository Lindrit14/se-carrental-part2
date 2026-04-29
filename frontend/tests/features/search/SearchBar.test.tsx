import { describe, expect, it, vi } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter, Routes, Route, useLocation } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { SearchBar } from '@/features/search/SearchBar';

// Avoid hitting the real Google Maps loader inside the jsdom env.
vi.mock('@/features/maps/GoogleMapsProvider', () => ({
  __esModule: true,
  default: ({ children }: { children: React.ReactNode }) => <>{children}</>,
  useGoogleMaps: () => ({ isLoaded: false }),
}));

function LocationProbe() {
  const loc = useLocation();
  return <div data-testid="location">{loc.pathname + loc.search}</div>;
}

function renderSearchBar() {
  const qc = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={qc}>
      <MemoryRouter initialEntries={['/']}>
        <Routes>
          <Route path="/" element={<SearchBar />} />
          <Route path="/cars" element={<LocationProbe />} />
        </Routes>
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

describe('SearchBar', () => {
  it('navigates to /cars with the form values as query params', async () => {
    renderSearchBar();
    const locInput = screen.getByLabelText(/pick-up location/i);
    await userEvent.type(locInput, 'Vienna');

    // The from/to inputs come pre-filled with today + today+3 — just submit.
    await userEvent.click(screen.getByRole('button', { name: /search/i }));

    await waitFor(() => {
      const probe = screen.getByTestId('location');
      expect(probe.textContent).toMatch(/^\/cars\?/);
      expect(probe.textContent).toContain('location=Vienna');
      expect(probe.textContent).toMatch(/from=\d{4}-\d{2}-\d{2}/);
      expect(probe.textContent).toMatch(/to=\d{4}-\d{2}-\d{2}/);
    });
  });
});
