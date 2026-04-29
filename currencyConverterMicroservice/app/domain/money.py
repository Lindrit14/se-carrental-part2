"""Value object for an amount in a given currency. No framework imports."""
from __future__ import annotations

from dataclasses import dataclass
from decimal import Decimal


@dataclass(frozen=True, slots=True)
class Money:
    amount: Decimal
    currency: str  # ISO-4217, uppercase

    def __post_init__(self) -> None:
        if not isinstance(self.amount, Decimal):
            object.__setattr__(self, "amount", Decimal(str(self.amount)))
        if not self.currency or len(self.currency) != 3:
            raise ValueError(f"invalid currency code: {self.currency!r}")
        object.__setattr__(self, "currency", self.currency.upper())
