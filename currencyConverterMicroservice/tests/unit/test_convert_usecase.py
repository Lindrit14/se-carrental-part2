from __future__ import annotations

from datetime import date
from decimal import Decimal

import pytest

from app.application.convert import ConvertUseCase
from app.domain.errors import UnknownCurrency
from app.domain.money import Money
from app.domain.rates import RateSet


class StaticRates:
    """Tiny test double that mimics RatesProvider."""

    def __init__(self) -> None:
        self._rs = RateSet(
            rate_date=date(2026, 4, 28),
            rates={"USD": Decimal("1.0823"), "GBP": Decimal("0.8634")},
        )

    def current(self) -> RateSet:
        return self._rs


def test_convert_usd_to_eur() -> None:
    uc = ConvertUseCase(StaticRates())
    res = uc.execute(Money(Decimal("100.00"), "USD"), "EUR")
    # 100 / 1.0823 = ~92.3958
    assert res.target.currency == "EUR"
    assert res.target.amount == Decimal("92.3958")


def test_convert_eur_to_usd() -> None:
    uc = ConvertUseCase(StaticRates())
    res = uc.execute(Money(Decimal("50"), "EUR"), "USD")
    # 50 * 1.0823 = 54.115
    assert res.target.amount == Decimal("54.1150")
    assert res.target.currency == "USD"


def test_convert_unknown_currency() -> None:
    uc = ConvertUseCase(StaticRates())
    with pytest.raises(UnknownCurrency):
        uc.execute(Money(Decimal("1"), "USD"), "XYZ")
