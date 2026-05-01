from __future__ import annotations

import logging
from decimal import Decimal, InvalidOperation

import grpc
from grpc_health.v1 import health, health_pb2, health_pb2_grpc

from app.application.convert import ConvertUseCase
from app.application.ports import RatesProvider
from app.domain.errors import RatesUnavailable, UnknownCurrency
from app.domain.money import Money
from app.infrastructure.grpc.generated.currency.v1 import currency_pb2, currency_pb2_grpc

logger = logging.getLogger(__name__)

SERVICE_NAME = "currency.v1.CurrencyConverter"


class CurrencyConverterServicer(currency_pb2_grpc.CurrencyConverterServicer):
    def __init__(self, usecase: ConvertUseCase, rates: RatesProvider) -> None:
        self._usecase = usecase
        self._rates = rates

    async def Convert(  # noqa: N802 - gRPC method name is fixed by the proto
        self,
        request: currency_pb2.ConvertRequest,
        context: grpc.aio.ServicerContext,
    ) -> currency_pb2.ConvertResponse:
        try:
            amount = Decimal(request.amount)
        except (InvalidOperation, ValueError) as e:
            await context.abort(
                grpc.StatusCode.INVALID_ARGUMENT,
                f"invalid amount {request.amount!r}: {e}",
            )

        try:
            source = Money(amount=amount, currency=request.from_currency)
        except ValueError as e:
            await context.abort(grpc.StatusCode.INVALID_ARGUMENT, str(e))

        try:
            result = self._usecase.execute(source, request.to_currency)
        except UnknownCurrency as e:
            await context.abort(
                grpc.StatusCode.INVALID_ARGUMENT,
                f"unknown currency: {e.currency}",
            )
        except RatesUnavailable:
            await context.abort(
                grpc.StatusCode.FAILED_PRECONDITION,
                "ECB rates not yet loaded",
            )
        except ValueError as e:
            await context.abort(grpc.StatusCode.INVALID_ARGUMENT, str(e))

        return currency_pb2.ConvertResponse(
            amount=str(result.target.amount),
            currency=result.target.currency,
            source_amount=str(result.source.amount),
            source_currency=result.source.currency,
            rate_date=result.rate_date.isoformat(),
        )

    async def GetRates(  # noqa: N802 - gRPC method name is fixed by the proto
        self,
        request: currency_pb2.GetRatesRequest,
        context: grpc.aio.ServicerContext,
    ) -> currency_pb2.GetRatesResponse:
        try:
            snapshot = self._rates.current()
        except RatesUnavailable:
            await context.abort(
                grpc.StatusCode.FAILED_PRECONDITION,
                "ECB rates not yet loaded",
            )

        return currency_pb2.GetRatesResponse(
            base="EUR",
            rate_date=snapshot.rate_date.isoformat(),
            rates=[
                currency_pb2.Rate(currency=code, rate=str(rate))
                for code, rate in sorted(snapshot.rates.items())
            ],
        )


class CurrencyGrpcServer:

    def __init__(self, port: int, usecase: ConvertUseCase, rates: RatesProvider) -> None:
        self._port = port
        self._usecase = usecase
        self._rates = rates
        self._server: grpc.aio.Server | None = None
        self._health: health.aio.HealthServicer | None = None

    async def start(self) -> None:
        self._server = grpc.aio.server()
        currency_pb2_grpc.add_CurrencyConverterServicer_to_server(
            CurrencyConverterServicer(self._usecase, self._rates), self._server
        )
        # Standard gRPC health-checking protocol so probes (grpc_health_probe,
        # k8s grpc liveness/readiness) can hit us without HTTP.
        self._health = health.aio.HealthServicer()
        health_pb2_grpc.add_HealthServicer_to_server(self._health, self._server)
        await self._health.set("", health_pb2.HealthCheckResponse.NOT_SERVING)
        await self._health.set(SERVICE_NAME, health_pb2.HealthCheckResponse.NOT_SERVING)

        bind_addr = f"[::]:{self._port}"
        self._server.add_insecure_port(bind_addr)
        await self._server.start()
        logger.info("grpc_server_started", extra={"address": bind_addr})

    async def set_serving(self, ok: bool) -> None:
        if self._health is None:
            return
        status = (
            health_pb2.HealthCheckResponse.SERVING
            if ok
            else health_pb2.HealthCheckResponse.NOT_SERVING
        )
        await self._health.set("", status)
        await self._health.set(SERVICE_NAME, status)

    async def wait_for_termination(self) -> None:
        if self._server is not None:
            await self._server.wait_for_termination()

    async def stop(self) -> None:
        if self._server is not None:
            # Flip health to NOT_SERVING before tearing down so probes drain
            # ahead of the actual stop; 5s grace lets in-flight calls finish.
            await self.set_serving(False)
            await self._server.stop(grace=5.0)
            self._server = None
            self._health = None
