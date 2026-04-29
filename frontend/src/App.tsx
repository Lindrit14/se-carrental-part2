import { BrowserRouter } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { AuthProvider } from '@/features/auth/AuthProvider';
import { CurrencyProvider } from '@/features/currency/CurrencyContext';
import GoogleMapsProvider from '@/features/maps/GoogleMapsProvider';
import { Toaster } from '@/components/ui/Toaster';
import { AppRouter } from '@/routes/AppRouter';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      refetchOnWindowFocus: false,
      staleTime: 30_000,
    },
  },
});

export function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <AuthProvider>
          <CurrencyProvider>
            <GoogleMapsProvider>
              <AppRouter />
              <Toaster />
            </GoogleMapsProvider>
          </CurrencyProvider>
        </AuthProvider>
      </BrowserRouter>
    </QueryClientProvider>
  );
}
