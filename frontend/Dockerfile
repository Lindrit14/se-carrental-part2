# syntax=docker/dockerfile:1.6
FROM node:20-alpine AS builder
WORKDIR /build

# Build args become Vite env vars at build time. Override per-environment.
ARG VITE_AUTH_URL=http://localhost:8080
ARG VITE_BOOKING_URL=http://localhost:8082
ARG VITE_CURRENCY_URL=http://localhost:8000
ARG VITE_GOOGLE_MAPS_API_KEY=
ENV VITE_AUTH_URL=$VITE_AUTH_URL \
    VITE_BOOKING_URL=$VITE_BOOKING_URL \
    VITE_CURRENCY_URL=$VITE_CURRENCY_URL \
    VITE_GOOGLE_MAPS_API_KEY=$VITE_GOOGLE_MAPS_API_KEY

COPY package.json package-lock.json* ./
RUN npm ci --no-audit --no-fund

COPY . .
RUN npm run build

# --- runtime ---------------------------------------------------------
FROM nginx:1.27-alpine AS runtime

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /build/dist /usr/share/nginx/html

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
