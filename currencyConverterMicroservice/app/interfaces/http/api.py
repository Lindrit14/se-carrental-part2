"""FastAPI routes. Handlers translate DTOs ↔ use case I/O."""
from __future__ import annotations

from fastapi import APIRouter, HTTPException, Request

from app.application.convert import ConvertUseCase
from app.application.ports import RatesProvider
from app.domain.errors import RatesUnavailable, UnknownCurrency
from app.domain.money import Money
from app.interfaces.http.dto import (
    ConvertRequest,
    ConvertResponse,
    RateEntry,
    RatesResponse,
)

router = APIRouter()


def _usecase(request: Request) -> ConvertUseCase:
    return request.app.state.convert_usecase


def _rates(request: Request) -> RatesProvider:
    return request.app.state.rates_provider


@router.get("/healthz")
def liveness() -> dict[str, str]:
    return {"status": "ok"}


@router.get("/api/v1/rates", response_model=RatesResponse)
def get_rates(request: Request) -> RatesResponse:
    try:
        snapshot = _rates(request).current()
    except RatesUnavailable as exc:
        raise HTTPException(status_code=503, detail="rates_unavailable") from exc
    return RatesResponse(
        base="EUR",
        rate_date=snapshot.rate_date.isoformat(),
        rates=[RateEntry(currency=c, rate=r) for c, r in sorted(snapshot.rates.items())],
    )


@router.post("/api/v1/convert", response_model=ConvertResponse)
def convert(req: ConvertRequest, request: Request) -> ConvertResponse:
    try:
        result = _usecase(request).execute(
            Money(amount=req.amount, currency=req.from_),
            req.to,
        )
    except UnknownCurrency as exc:
        raise HTTPException(status_code=400, detail=f"unknown_currency:{exc.currency}") from exc
    except RatesUnavailable as exc:
        raise HTTPException(status_code=503, detail="rates_unavailable") from exc

    return ConvertResponse(
        amount=result.target.amount,
        currency=result.target.currency,
        source_amount=result.source.amount,
        source_currency=result.source.currency,
        rate_date=result.rate_date.isoformat(),
    )
