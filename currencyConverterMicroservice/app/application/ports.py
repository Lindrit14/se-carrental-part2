"""Outbound ports the application depends on. Adapters implement these."""
from __future__ import annotations

from typing import Protocol

from app.domain.rates import RateSet


class RatesProvider(Protocol):
    """Source of truth for current exchange rates."""

    def current(self) -> RateSet: ...
