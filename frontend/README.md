# Frontend

Vite + React 18 + TypeScript + Tailwind v4. Talks to the platform's three
backend services via REST + JWT (Bearer).

## Architecture

Clean Architecture — imports flow inward:

```
pages → features → components/ui + lib → domain
```

```
src/
  domain/           Pure types (User, Car, Booking, Money) — no framework imports
  lib/              Framework-free: API client, token storage, formatters
  components/
    ui/             shadcn-style primitives (Button, Input, Dialog, …)
    layout/         AppShell, Navbar, PageHeader
  features/         Auth, cars, bookings, currency, profile (UI + React Query hooks)
  pages/            Thin route-level components composing features
  routes/           Router + ProtectedRoute
```

API calls live **only** in `lib/api/`. UI consumes data through React-Query
hooks under `features/<x>/use*.ts`.

## Run locally

```bash
cp .env.example .env       # default URLs match the platform compose
npm install
npm run dev                # http://localhost:5173

# Backend (separate terminal)
cd .. && make up
```

## Scripts

```bash
npm run dev          # Vite dev server
npm run build        # tsc -b && vite build → dist/
npm run preview      # serve the production build locally
npm run typecheck    # tsc -b --noEmit
npm run lint         # eslint
npm run test         # vitest run (smoke tests)
npm run test:watch
npm run format       # prettier
```

## Production image

```bash
docker build -t frontend:dev .
docker run --rm -p 3000:80 frontend:dev
# → http://localhost:3000
```

The platform compose builds this image automatically as the `frontend`
service.

## Security note

Access + refresh tokens are stored in `localStorage` for simplicity. That
makes them readable by any script running on this origin (XSS exposure).
For production use HttpOnly cookies + CSRF protection issued by user-auth.
