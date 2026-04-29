"""ConvertUseCase: converts an amount between two currencies via the EUR pivot."""
from __future__ import annotations

from dataclasses import dataclass
from datetime import date
from decimal import ROUND_HALF_UP, Decimal

from app.application.ports import RatesProvider
from app.domain.money import Money


@dataclass(frozen=True, slots=True)
class ConvertResult:
    source: Money
    target: Money
    rate_date: date


class ConvertUseCase:
    """Convert ``source`` Money into ``target_currency`` using the latest RateSet.

    ECB rates are quoted as ``rate[X] = X per 1 EUR``. So:
        amount_EUR = amount_X / rate[X]
        amount_Y   = amount_EUR * rate[Y]
    """

    def __init__(self, rates: RatesProvider) -> None:
        self._rates = rates

    def execute(self, source: Money, target_currency: str) -> ConvertResult:
        rate_set = self._rates.current()
        rate_from = rate_set.rate_for(source.currency)
        rate_to = rate_set.rate_for(target_currency)

        amount_in_target = (source.amount / rate_from) * rate_to
        amount_in_target = amount_in_target.quantize(Decimal("0.0001"), rounding=ROUND_HALF_UP)

        return ConvertResult(
            source=source,
            target=Money(amount=amount_in_target, currency=target_currency),
            rate_date=rate_set.rate_date,
        )
