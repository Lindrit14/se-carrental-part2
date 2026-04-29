/**
 * Minimal global typings for the Google Maps JS API surface we actually use.
 * Avoids pulling the full `@types/google.maps` package as a dependency just
 * for one Autocomplete widget.
 */
export {};

declare global {
  // eslint-disable-next-line @typescript-eslint/no-namespace
  namespace google.maps.places {
    interface AutocompleteOptions {
      fields?: string[];
      types?: string[];
      componentRestrictions?: { country: string | string[] };
    }

    interface PlaceResult {
      formatted_address?: string;
      name?: string;
      place_id?: string;
    }

    class Autocomplete {
      constructor(input: HTMLInputElement, opts?: AutocompleteOptions);
      addListener(event: 'place_changed', handler: () => void): { remove(): void };
      getPlace(): PlaceResult;
    }
  }

  interface Window {
    google?: typeof google;
  }
}
