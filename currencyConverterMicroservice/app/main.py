"""Composition root. Builds adapters, wires the use case, registers HTTP routes,
and starts the ECB refresher + gRPC server inside the FastAPI lifespan.
"""
from __future__ import annotations

import logging
from contextlib import asynccontextmanager
from typing import AsyncIterator

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from app.application.convert import ConvertUseCase
from app.infrastructure.config import Settings
from app.infrastructure.ecb_provider import ECBRatesProvider
from app.infrastructure.grpc.server import CurrencyGrpcServer
from app.infrastructure.logging import configure_logging
from app.interfaces.http.api import router

logger = logging.getLogger(__name__)


def _build_app() -> FastAPI:
    settings = Settings()
    configure_logging(settings.log_level)

    rates = ECBRatesProvider(
        feed_url=settings.ecb_feed_url,
        refresh_interval=settings.ecb_refresh_interval,
        fetch_timeout=settings.ecb_fetch_timeout,
    )
    usecase = ConvertUseCase(rates)
    grpc_server = CurrencyGrpcServer(port=settings.grpc_port, usecase=usecase)

    @asynccontextmanager
    async def lifespan(app_: FastAPI) -> AsyncIterator[None]:
        logger.info("startup_begin")
        await rates.start()
        await grpc_server.start()
        app_.state.rates_provider = rates
        app_.state.convert_usecase = usecase
        logger.info("startup_complete")
        try:
            yield
        finally:
            logger.info("shutdown_begin")
            await grpc_server.stop()
            await rates.stop()
            logger.info("shutdown_complete")

    application = FastAPI(title="currency-converter", version="0.1.0", lifespan=lifespan)
    application.add_middleware(
        CORSMiddleware,
        allow_origins=settings.cors_origins_list,
        allow_methods=["GET", "POST", "OPTIONS"],
        allow_headers=["Authorization", "Content-Type"],
        allow_credentials=False,
        max_age=3600,
    )
    application.include_router(router)
    return application


app = _build_app()
