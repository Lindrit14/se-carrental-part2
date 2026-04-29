# syntax=docker/dockerfile:1.6
FROM golang:1.26-alpine AS builder
WORKDIR /src

# Cache mod download
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

# Build
COPY . .
ARG VERSION=dev
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux \
    go build -trimpath -ldflags="-s -w -X main.version=${VERSION}" \
      -o /out/server ./cmd/server

# --- Final image ----------------------------------------------------
FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=builder /out/server /app/server
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/app/server"]
