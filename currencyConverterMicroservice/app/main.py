"""Composition root. Builds adapters, wires the use case, starts the ECB
refresher and the gRPC server. Pure asyncio entrypoint — no HTTP."""
from __future__ import annotations

import asyncio
import logging
import signal

from app.application.convert import ConvertUseCase
from app.infrastructure.config import Settings
from app.infrastructure.ecb_provider import ECBRatesProvider
from app.infrastructure.grpc.server import CurrencyGrpcServer
from app.infrastructure.logging import configure_logging

logger = logging.getLogger(__name__)


async def _run() -> None:
    settings = Settings()
    configure_logging(settings.log_level)

    rates = ECBRatesProvider(
        feed_url=settings.ecb_feed_url,
        refresh_interval=settings.ecb_refresh_interval,
        fetch_timeout=settings.ecb_fetch_timeout,
    )
    usecase = ConvertUseCase(rates)
    server = CurrencyGrpcServer(port=settings.grpc_port, usecase=usecase, rates=rates)

    logger.info("startup_begin")
    # Bind & start the gRPC server first (health = NOT_SERVING) so probes can
    # connect immediately and report unready while ECB rates are still loading.
    await server.start()

    stop = asyncio.Event()
    loop = asyncio.get_running_loop()
    for sig in (signal.SIGTERM, signal.SIGINT):
        try:
            loop.add_signal_handler(sig, stop.set)
        except NotImplementedError:
            # Windows / non-mainthread fall-back: rely on KeyboardInterrupt.
            pass

    try:
        await rates.start()  # awaits the first ECB fetch
        await server.set_serving(True)
        logger.info("startup_complete")
        await stop.wait()
    finally:
        logger.info("shutdown_begin")
        await server.stop()  # flips health to NOT_SERVING + drains
        await rates.stop()
        logger.info("shutdown_complete")


def main() -> None:
    asyncio.run(_run())


if __name__ == "__main__":
    main()
