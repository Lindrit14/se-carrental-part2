class DomainError(Exception):
    """Base for domain-level errors."""


class UnknownCurrency(DomainError):
    def __init__(self, currency: str) -> None:
        super().__init__(f"unknown currency: {currency}")
        self.currency = currency


class RatesUnavailable(DomainError):
    """Raised when no rates have ever been loaded."""
