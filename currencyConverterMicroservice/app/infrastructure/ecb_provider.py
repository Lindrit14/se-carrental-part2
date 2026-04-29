"""ECBRatesProvider — fetches and caches the daily EUR-FX feed."""
from __future__ import annotations

import asyncio
import logging
from datetime import date, datetime
from decimal import Decimal
from xml.etree import ElementTree as ET

import httpx

from app.application.ports import RatesProvider
from app.domain.errors import RatesUnavailable
from app.domain.rates import RateSet

logger = logging.getLogger(__name__)

ECB_NS = {
    "gesmes": "http://www.gesmes.org/xml/2002-08-01",
    "eurofx": "http://www.ecb.int/vocabulary/2002-08-01/eurofxref",
}


def parse_ecb_xml(xml_text: str) -> RateSet:
    """Parse the ECB daily XML into a RateSet (EUR-based).

    The feed structure is:
        <gesmes:Envelope>
          <Cube>
            <Cube time="2026-04-28">
              <Cube currency="USD" rate="1.0823"/>
              ...
            </Cube>
          </Cube>
        </gesmes:Envelope>
    """
    root = ET.fromstring(xml_text)
    day = root.find(".//eurofx:Cube/eurofx:Cube[@time]", ECB_NS)
    if day is None:
        raise ValueError("ECB feed: no dated <Cube> element")

    time_attr = day.attrib.get("time", "")
    rate_date = datetime.strptime(time_attr, "%Y-%m-%d").date()

    rates: dict[str, Decimal] = {}
    for entry in day.findall("eurofx:Cube", ECB_NS):
        currency = entry.attrib.get("currency")
        rate = entry.attrib.get("rate")
        if currency and rate:
            rates[currency.upper()] = Decimal(rate)
    return RateSet(rate_date=rate_date, rates=rates)


class ECBRatesProvider(RatesProvider):
    """Background-refreshed in-memory cache of ECB rates."""

    def __init__(self, feed_url: str, refresh_interval: int, fetch_timeout: float) -> None:
        self._url = feed_url
        self._interval = refresh_interval
        self._timeout = fetch_timeout
        self._cache: RateSet | None = None
        self._task: asyncio.Task[None] | None = None

    def current(self) -> RateSet:
        if self._cache is None:
            raise RatesUnavailable("rates not yet loaded")
        return self._cache

    async def refresh_once(self) -> None:
        """Fetch the feed once. On failure, log and keep the previous cache."""
        try:
            async with httpx.AsyncClient(timeout=self._timeout) as client:
                resp = await client.get(self._url)
                resp.raise_for_status()
            self._cache = parse_ecb_xml(resp.text)
            logger.info(
                "ecb_rates_refreshed",
                extra={"date": str(self._cache.rate_date), "count": len(self._cache.rates)},
            )
        except Exception:
            if self._cache is None:
                logger.exception("ecb_initial_fetch_failed")
                raise
            logger.warning("ecb_refresh_failed_keeping_stale_cache", exc_info=True)

    async def start(self) -> None:
        await self.refresh_once()
        self._task = asyncio.create_task(self._loop(), name="ecb-refresher")

    async def stop(self) -> None:
        if self._task is not None:
            self._task.cancel()
            try:
                await self._task
            except (asyncio.CancelledError, Exception):
                pass
            self._task = None

    async def _loop(self) -> None:
        while True:
            try:
                await asyncio.sleep(self._interval)
                await self.refresh_once()
            except asyncio.CancelledError:
                return
            except Exception:
                logger.exception("ecb_refresh_loop_error")
