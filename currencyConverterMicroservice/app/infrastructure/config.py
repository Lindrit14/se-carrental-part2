"""Centralised settings. The only place that reads ENV vars."""
from __future__ import annotations

from pydantic import Field
from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    model_config = SettingsConfigDict(env_file=".env", extra="ignore")

    log_level: str = Field(default="info", alias="LOG_LEVEL")
    grpc_port: int = Field(default=9000, alias="GRPC_PORT")

    ecb_feed_url: str = Field(
        default="https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml",
        alias="ECB_FEED_URL",
    )
    ecb_refresh_interval: int = Field(default=3600, alias="ECB_REFRESH_INTERVAL")
    ecb_fetch_timeout: float = Field(default=10.0, alias="ECB_FETCH_TIMEOUT")
