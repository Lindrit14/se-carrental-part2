import { createContext, useContext, useEffect, useState } from 'react';

interface GoogleMapsContextValue {
  isLoaded: boolean;
}

const GoogleMapsContext = createContext<GoogleMapsContextValue>({ isLoaded: false });

export function useGoogleMaps() {
  return useContext(GoogleMapsContext);
}

let loadPromise: Promise<void> | null = null;

function loadGoogleMapsScript(apiKey: string): Promise<void> {
  if (loadPromise) return loadPromise;

  loadPromise = new Promise((resolve, reject) => {
    if (document.querySelector('script[src*="maps.googleapis.com"]')) {
      resolve();
      return;
    }

    const script = document.createElement('script');
    script.src = `https://maps.googleapis.com/maps/api/js?key=${apiKey}&libraries=places&loading=async`;
    script.async = true;
    script.defer = true;
    script.onload = () => resolve();
    script.onerror = () => reject(new Error('Failed to load Google Maps script'));
    document.head.appendChild(script);
  });

  return loadPromise;
}

export default function GoogleMapsProvider({ children }: { children: React.ReactNode }) {
  const [isLoaded, setIsLoaded] = useState(false);
  const apiKey = import.meta.env.VITE_GOOGLE_MAPS_API_KEY ?? '';

  useEffect(() => {
    if (!apiKey) return;
    loadGoogleMapsScript(apiKey).then(() => setIsLoaded(true));
  }, [apiKey]);

  return (
    <GoogleMapsContext.Provider value={{ isLoaded }}>{children}</GoogleMapsContext.Provider>
  );
}
