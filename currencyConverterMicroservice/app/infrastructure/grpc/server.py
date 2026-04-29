"""gRPC server: implements ``currency.v1.CurrencyConverter`` over the existing
``ConvertUseCase``. Replaces the previous AMQP RPC server.

Domain errors are mapped to gRPC status codes so callers can react with the
standard tooling rather than parsing a JSON envelope:

  UnknownCurrency      -> INVALID_ARGUMENT
  RatesUnavailable     -> FAILED_PRECONDITION
  ValueError/TypeError -> INVALID_ARGUMENT (malformed request)
  unexpected           -> INTERNAL
"""
from __future__ import annotations

import logging
from decimal import Decimal, InvalidOperation

import grpc

from app.application.convert import ConvertUseCase
from app.domain.errors import RatesUnavailable, UnknownCurrency
from app.domain.money import Money
from app.infrastructure.grpc.generated.currency.v1 import currency_pb2, currency_pb2_grpc

logger = logging.getLogger(__name__)


class CurrencyConverterServicer(currency_pb2_grpc.CurrencyConverterServicer):
    def __init__(self, usecase: ConvertUseCase) -> None:
        self._usecase = usecase

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


class CurrencyGrpcServer:
    """Lifecycle wrapper for ``grpc.aio.server`` — start/stop hooks match the
    shape of the previous ``CurrencyRpcServer`` so ``main.py``'s lifespan stays
    structurally identical.
    """

    def __init__(self, port: int, usecase: ConvertUseCase) -> None:
        self._port = port
        self._usecase = usecase
        self._server: grpc.aio.Server | None = None

    async def start(self) -> None:
        self._server = grpc.aio.server()
        currency_pb2_grpc.add_CurrencyConverterServicer_to_server(
            CurrencyConverterServicer(self._usecase), self._server
        )
        bind_addr = f"[::]:{self._port}"
        self._server.add_insecure_port(bind_addr)
        await self._server.start()
        logger.info("grpc_server_started", extra={"address": bind_addr})

    async def stop(self) -> None:
        if self._server is not None:
            # 5s grace lets in-flight Convert calls finish; matches the
            # caller-side deadline so we never cut off a healthy request.
            await self._server.stop(grace=5.0)
            self._server = None
