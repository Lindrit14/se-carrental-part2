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
from grpc_health.v1 import health_pb2, health_pb2_grpc

from app.application.ports import RatesProvider
from app.infrastructure.grpc.generated.currency.v1 import currency_pb2, currency_pb2_grpc
from app.infrastructure.grpc.server import (
    SERVICE_NAME,
    CurrencyConverterServicer,
    CurrencyGrpcServer,
)


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


async def _start_server(usecase: ConvertUseCase, rates) -> tuple[grpc.aio.Server, int]:
    server = grpc.aio.server()
    currency_pb2_grpc.add_CurrencyConverterServicer_to_server(
        CurrencyConverterServicer(usecase, rates), server
    )
    port = server.add_insecure_port("127.0.0.1:0")
    await server.start()
    return server, port


@pytest.mark.asyncio
async def test_convert_happy_path() -> None:
    rates = StaticRates()
    server, port = await _start_server(ConvertUseCase(rates), rates)
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
    rates = StaticRates()
    server, port = await _start_server(ConvertUseCase(rates), rates)
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
    rates = EmptyRates()
    server, port = await _start_server(ConvertUseCase(rates), rates)
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
    rates = StaticRates()
    server, port = await _start_server(ConvertUseCase(rates), rates)
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


@pytest.mark.asyncio
async def test_get_rates_happy_path() -> None:
    rates = StaticRates()
    server, port = await _start_server(ConvertUseCase(rates), rates)
    try:
        async with grpc.aio.insecure_channel(f"127.0.0.1:{port}") as channel:
            stub = currency_pb2_grpc.CurrencyConverterStub(channel)
            reply = await stub.GetRates(currency_pb2.GetRatesRequest())
        assert reply.base == "EUR"
        assert reply.rate_date == "2026-04-28"
        # Sorted by ISO code, EUR base not echoed back as a row.
        codes = [r.currency for r in reply.rates]
        assert codes == ["GBP", "USD"]
        rates_by_code = {r.currency: r.rate for r in reply.rates}
        assert rates_by_code == {"USD": "1.0823", "GBP": "0.8634"}
    finally:
        await server.stop(grace=None)


@pytest.mark.asyncio
async def test_get_rates_unavailable_returns_failed_precondition() -> None:
    rates = EmptyRates()
    server, port = await _start_server(ConvertUseCase(rates), rates)
    try:
        async with grpc.aio.insecure_channel(f"127.0.0.1:{port}") as channel:
            stub = currency_pb2_grpc.CurrencyConverterStub(channel)
            with pytest.raises(grpc.aio.AioRpcError) as exc:
                await stub.GetRates(currency_pb2.GetRatesRequest())
            assert exc.value.code() == grpc.StatusCode.FAILED_PRECONDITION
    finally:
        await server.stop(grace=None)


async def _start_full_server(rates: RatesProvider) -> tuple[CurrencyGrpcServer, int]:
    """Start the full CurrencyGrpcServer (incl. health) on an ephemeral port."""
    server = CurrencyGrpcServer(port=0, usecase=ConvertUseCase(rates), rates=rates)
    # Bypass the `[::]:{port}` formatting so we can grab a real port from grpc.
    server._server = grpc.aio.server()
    currency_pb2_grpc.add_CurrencyConverterServicer_to_server(
        CurrencyConverterServicer(rates=rates, usecase=ConvertUseCase(rates)),
        server._server,
    )
    from grpc_health.v1 import health as health_mod  # local to avoid top-level cycle
    server._health = health_mod.aio.HealthServicer()
    health_pb2_grpc.add_HealthServicer_to_server(server._health, server._server)
    await server._health.set("", health_pb2.HealthCheckResponse.NOT_SERVING)
    await server._health.set(SERVICE_NAME, health_pb2.HealthCheckResponse.NOT_SERVING)
    port = server._server.add_insecure_port("127.0.0.1:0")
    await server._server.start()
    return server, port


@pytest.mark.asyncio
async def test_health_starts_not_serving_and_flips_to_serving() -> None:
    """Health gating mirrors readiness: NOT_SERVING until set_serving(True)."""
    server, port = await _start_full_server(StaticRates())
    try:
        async with grpc.aio.insecure_channel(f"127.0.0.1:{port}") as channel:
            health_stub = health_pb2_grpc.HealthStub(channel)

            reply = await health_stub.Check(health_pb2.HealthCheckRequest(service=""))
            assert reply.status == health_pb2.HealthCheckResponse.NOT_SERVING

            await server.set_serving(True)

            reply = await health_stub.Check(health_pb2.HealthCheckRequest(service=""))
            assert reply.status == health_pb2.HealthCheckResponse.SERVING
            reply = await health_stub.Check(
                health_pb2.HealthCheckRequest(service=SERVICE_NAME)
            )
            assert reply.status == health_pb2.HealthCheckResponse.SERVING
    finally:
        await server.stop()
