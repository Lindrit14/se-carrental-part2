import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { SearchBar } from '@/features/search/SearchBar';
import { CarSearchResultsPage } from '@/pages/CarSearchResultsPage';
import { AuthContext } from '@/features/auth/AuthProvider';
import { CurrencyContext } from '@/features/currency/CurrencyContext';
import { tokenStorage } from '@/lib/auth/tokenStorage';

vi.mock('@/features/maps/GoogleMapsProvider', () => ({
  __esModule: true,
  default: ({ children }: { children: React.ReactNode }) => <>{children}</>,
  useGoogleMaps: () => ({ isLoaded: false }),
}));

const CARS_RESPONSE = [
  {
    id: 'car-1',
    brand: 'Toyota',
    model: 'Yaris',
    licensePlate: 'W-1234',
    dailyRateAmount: '42.00',
    dailyRateCurrency: 'EUR',
    location: 'Vienna',
    category: 'SMALL',
  },
];

function jsonResponse(body: unknown, status = 200): Response {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'content-type': 'application/json' },
  });
}

function renderUnauthenticatedHomeWithCarsRoute(initialPath: string = '/') {
  const qc = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={qc}>
      <AuthContext.Provider
        value={{
          session: null,
          initialized: true,
          login: vi.fn(),
          register: vi.fn(),
          logout: vi.fn().mockResolvedValue(undefined),
        }}
      >
        <CurrencyContext.Provider
          value={{
            displayCurrency: 'EUR',
            setDisplayCurrency: vi.fn(),
            rates: new Map([['EUR', 1]]),
            available: ['EUR'],
            ratesLoading: false,
            ratesError: false,
            rateDate: null,
          }}
        >
          <MemoryRouter initialEntries={[initialPath]}>
            <Routes>
              <Route path="/" element={<SearchBar />} />
              <Route path="/cars" element={<CarSearchResultsPage />} />
            </Routes>
          </MemoryRouter>
        </CurrencyContext.Provider>
      </AuthContext.Provider>
    </QueryClientProvider>,
  );
}

describe('Unauthenticated car search', () => {
  beforeEach(() => {
    tokenStorage.clear();
  });
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('hits /api/v1/cars/search without an Authorization header and renders the results', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse(CARS_RESPONSE));
    vi.stubGlobal('fetch', fetchMock);

    renderUnauthenticatedHomeWithCarsRoute();

    await userEvent.type(screen.getByLabelText(/pick-up location/i), 'Vienna');
    await userEvent.click(screen.getByRole('button', { name: /search/i }));

    // The results page should render the car returned by the search endpoint.
    expect(await screen.findByText(/toyota\s+yaris/i)).toBeInTheDocument();

    // Verify the fetch contract: GET /api/v1/cars/search with our query
    // params, and crucially no Authorization header since there's no session.
    const searchCall = fetchMock.mock.calls.find(([url]) =>
      String(url).includes('/api/v1/cars/search'),
    );
    expect(searchCall, 'expected a call to /api/v1/cars/search').toBeDefined();

    const [url, init] = searchCall!;
    const requestUrl = String(url);
    expect(requestUrl).toContain('http://localhost:8080/api/v1/cars/search');
    expect(requestUrl).toContain('location=Vienna');
    expect(requestUrl).toMatch(/from=\d{4}-\d{2}-\d{2}/);
    expect(requestUrl).toMatch(/to=\d{4}-\d{2}-\d{2}/);
    expect(init?.method ?? 'GET').toBe('GET');
    expect((init?.headers as Record<string, string>).Authorization).toBeUndefined();
  });

  it('loads all cars when /cars is opened without a pickup location', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse(CARS_RESPONSE));
    vi.stubGlobal('fetch', fetchMock);

    renderUnauthenticatedHomeWithCarsRoute('/cars');

    // The page should fetch and render the unfiltered car list.
    expect(await screen.findByText(/toyota\s+yaris/i)).toBeInTheDocument();

    const searchCall = fetchMock.mock.calls.find(([url]) =>
      String(url).includes('/api/v1/cars/search'),
    );
    expect(searchCall, 'expected a call to /api/v1/cars/search').toBeDefined();

    const [url, init] = searchCall!;
    const requestUrl = String(url);
    // No filters → carsApi builds the path with no query string.
    expect(requestUrl).toBe('http://localhost:8080/api/v1/cars/search');
    expect(init?.method ?? 'GET').toBe('GET');
    expect((init?.headers as Record<string, string>).Authorization).toBeUndefined();
  });
});
