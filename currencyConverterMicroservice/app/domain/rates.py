"""ECB rates are EUR-based: ``rates[X]`` = how many X you get for 1 EUR."""
from __future__ import annotations

from dataclasses import dataclass
from datetime import date
from decimal import Decimal

from .errors import UnknownCurrency

EUR = "EUR"


@dataclass(frozen=True, slots=True)
class RateSet:
    """A snapshot of EUR-based exchange rates for a given date."""

    rate_date: date
    rates: dict[str, Decimal]  # currency code → rate vs EUR (e.g. {"USD": 1.0823})

    def rate_for(self, currency: str) -> Decimal:
        currency = currency.upper()
        if currency == EUR:
            return Decimal("1")
        rate = self.rates.get(currency)
        if rate is None:
            raise UnknownCurrency(currency)
        return rate

    def supported(self) -> list[str]:
        return [EUR, *sorted(self.rates.keys())]
