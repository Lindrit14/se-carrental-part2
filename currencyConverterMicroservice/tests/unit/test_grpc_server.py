"""Unit tests for the gRPC server: spin up grpc.aio.server() on an ephemeral
port with a real ConvertUseCase backed by a fake RatesProvider, then dial it
with a real client. Covers happy path + the two domain-error mappings.
"""
from __future__ import annotations

from datetime import date
from decimal import Decimal

import grpc
import pytest

from app.application.convert import ConvertUseCase
from app.domain.errors import RatesUnavailable
from app.domain.rates import RateSet
from app.infrastructure.grpc.generated.currency.v1 import currency_pb2, currency_pb2_grpc
from app.infrastructure.grpc.server import CurrencyConverterServicer


class StaticRates:
    def __init__(self) -> None:
        self._rs = RateSet(
            rate_date=date(2026, 4, 28),
            rates={"USD": Decimal("1.0823"), "GBP": Decimal("0.8634")},
        )

    def current(self) -> RateSet:
        return self._rs


class EmptyRates:
    def current(self) -> RateSet:
        raise RatesUnavailable()


async def _start_server(usecase: ConvertUseCase) -> tuple[grpc.aio.Server, int]:
    server = grpc.aio.server()
    currency_pb2_grpc.add_CurrencyConverterServicer_to_server(
        CurrencyConverterServicer(usecase), server
    )
    port = server.add_insecure_port("127.0.0.1:0")
    await server.start()
    return server, port


@pytest.mark.asyncio
async def test_convert_happy_path() -> None:
    server, port = await _start_server(ConvertUseCase(StaticRates()))
    try:
        async with grpc.aio.insecure_channel(f"127.0.0.1:{port}") as channel:
            stub = currency_pb2_grpc.CurrencyConverterStub(channel)
            reply = await stub.Convert(
                currency_pb2.ConvertRequest(
                    amount="100.00", from_currency="USD", to_currency="EUR"
                )
            )
        assert reply.currency == "EUR"
        # 100 / 1.0823 ≈ 92.3958
        assert reply.amount == "92.3958"
        assert reply.source_amount == "100.00"
        assert reply.source_currency == "USD"
        assert reply.rate_date == "2026-04-28"
    finally:
        await server.stop(grace=None)


@pytest.mark.asyncio
async def test_unknown_currency_returns_invalid_argument() -> None:
    server, port = await _start_server(ConvertUseCase(StaticRates()))
    try:
        async with grpc.aio.insecure_channel(f"127.0.0.1:{port}") as channel:
            stub = currency_pb2_grpc.CurrencyConverterStub(channel)
            with pytest.raises(grpc.aio.AioRpcError) as exc:
                await stub.Convert(
                    currency_pb2.ConvertRequest(
                        amount="1", from_currency="USD", to_currency="XYZ"
                    )
                )
            assert exc.value.code() == grpc.StatusCode.INVALID_ARGUMENT
            assert "XYZ" in (exc.value.details() or "")
    finally:
        await server.stop(grace=None)


@pytest.mark.asyncio
async def test_rates_unavailable_returns_failed_precondition() -> None:
    server, port = await _start_server(ConvertUseCase(EmptyRates()))
    try:
        async with grpc.aio.insecure_channel(f"127.0.0.1:{port}") as channel:
            stub = currency_pb2_grpc.CurrencyConverterStub(channel)
            with pytest.raises(grpc.aio.AioRpcError) as exc:
                await stub.Convert(
                    currency_pb2.ConvertRequest(
                        amount="1", from_currency="USD", to_currency="EUR"
                    )
                )
            assert exc.value.code() == grpc.StatusCode.FAILED_PRECONDITION
    finally:
        await server.stop(grace=None)


@pytest.mark.asyncio
async def test_malformed_amount_returns_invalid_argument() -> None:
    server, port = await _start_server(ConvertUseCase(StaticRates()))
    try:
        async with grpc.aio.insecure_channel(f"127.0.0.1:{port}") as channel:
            stub = currency_pb2_grpc.CurrencyConverterStub(channel)
            with pytest.raises(grpc.aio.AioRpcError) as exc:
                await stub.Convert(
                    currency_pb2.ConvertRequest(
                        amount="not-a-number", from_currency="USD", to_currency="EUR"
                    )
                )
            assert exc.value.code() == grpc.StatusCode.INVALID_ARGUMENT
    finally:
        await server.stop(grace=None)
