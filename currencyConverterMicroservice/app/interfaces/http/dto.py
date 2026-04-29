"""Pydantic request/response shapes. Domain types stay inside the app/."""
from __future__ import annotations

from decimal import Decimal

from pydantic import BaseModel, Field


class ConvertRequest(BaseModel):
    amount: Decimal = Field(..., gt=0)
    from_: str = Field(..., alias="from", min_length=3, max_length=3)
    to: str = Field(..., min_length=3, max_length=3)


class ConvertResponse(BaseModel):
    amount: Decimal
    currency: str
    source_amount: Decimal
    source_currency: str
    rate_date: str


class RateEntry(BaseModel):
    currency: str
    rate: Decimal


class RatesResponse(BaseModel):
    base: str
    rate_date: str
    rates: list[RateEntry]


class ErrorResponse(BaseModel):
    code: str
    message: str
