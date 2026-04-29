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

# Non-root user
RUN useradd --create-home --uid 1000 app
WORKDIR /app

COPY --from=builder /usr/local/lib/python3.12/site-packages /usr/local/lib/python3.12/site-packages
COPY --from=builder /usr/local/bin /usr/local/bin

COPY --from=builder --chown=app:app /build/app /app/app

USER app
EXPOSE 8000 9000

CMD ["python", "-m", "uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"]
