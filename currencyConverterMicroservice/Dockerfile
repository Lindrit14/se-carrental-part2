# syntax=docker/dockerfile:1.6
FROM python:3.12-slim AS builder

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    PIP_NO_CACHE_DIR=1 \
    PIP_DISABLE_PIP_VERSION_CHECK=1

WORKDIR /build

# Install uv for fast resolution + reproducible installs.
RUN pip install --no-cache-dir uv==0.5.4

COPY pyproject.toml ./
# Generate a lockfile inside the build (no uv.lock needed in repo).
RUN uv pip install --system --no-cache .
# grpcio-tools is dev-only — install it just for codegen.
RUN uv pip install --system --no-cache "grpcio-tools>=1.66.0"

# Codegen the proto stubs into the same layout the runtime imports.
COPY proto ./proto
COPY scripts ./scripts
COPY app ./app
RUN bash scripts/gen-proto.sh

# --- runtime ---------------------------------------------------------
FROM python:3.12-slim AS runtime

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1

# grpc_health_probe — used by docker / k8s liveness+readiness probes.
ARG GRPC_HEALTH_PROBE_VERSION=v0.4.25
RUN apt-get update \
 && apt-get install -y --no-install-recommends ca-certificates wget \
 && wget -qO /usr/local/bin/grpc_health_probe \
    https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-$(dpkg --print-architecture) \
 && chmod +x /usr/local/bin/grpc_health_probe \
 && apt-get purge -y wget && apt-get autoremove -y && rm -rf /var/lib/apt/lists/*

# Non-root user
RUN useradd --create-home --uid 1000 app
WORKDIR /app

COPY --from=builder /usr/local/lib/python3.12/site-packages /usr/local/lib/python3.12/site-packages
COPY --from=builder /usr/local/bin /usr/local/bin

COPY --from=builder --chown=app:app /build/app /app/app

USER app
EXPOSE 9000

HEALTHCHECK --interval=10s --timeout=5s --retries=10 --start-period=30s \
  CMD ["grpc_health_probe", "-addr=:9000"]

CMD ["python", "-m", "app.main"]
