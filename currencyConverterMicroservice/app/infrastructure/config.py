"""Centralised settings. The only place that reads ENV vars."""
from __future__ import annotations

from pydantic import Field
from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    model_config = SettingsConfigDict(env_file=".env", extra="ignore")

    http_port: int = Field(default=8000, alias="HTTP_PORT")
    log_level: str = Field(default="info", alias="LOG_LEVEL")

    rabbitmq_url: str = Field(
        default="amqp://guest:guest@localhost:5672/", alias="RABBITMQ_URL"
    )
    rabbitmq_rpc_queue: str = Field(default="currency.requests", alias="RABBITMQ_RPC_QUEUE")

    ecb_feed_url: str = Field(
        default="https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml",
        alias="ECB_FEED_URL",
    )
    ecb_refresh_interval: int = Field(default=3600, alias="ECB_REFRESH_INTERVAL")
    ecb_fetch_timeout: float = Field(default=10.0, alias="ECB_FETCH_TIMEOUT")

    # Comma-separated list of origins allowed to call REST endpoints.
    cors_allowed_origins: str = Field(
        default="http://localhost:3000,http://localhost:5173",
        alias="CORS_ALLOWED_ORIGINS",
    )

    @property
    def cors_origins_list(self) -> list[str]:
        return [o.strip() for o in self.cors_allowed_origins.split(",") if o.strip()]
