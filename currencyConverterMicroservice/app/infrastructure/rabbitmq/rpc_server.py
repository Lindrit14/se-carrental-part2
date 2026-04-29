"""RabbitMQ-RPC server: consumes from ``currency.requests`` and replies on ``reply_to``.

This is what booking calls via Spring AMQP's ``convertSendAndReceive`` — the
``correlation_id`` lets the caller match the reply to its waiting Future.
"""
from __future__ import annotations

import json
import logging
from decimal import Decimal

import aio_pika

from app.application.convert import ConvertUseCase
from app.domain.errors import RatesUnavailable, UnknownCurrency
from app.domain.money import Money

logger = logging.getLogger(__name__)


class CurrencyRpcServer:
    def __init__(
        self,
        amqp_url: str,
        queue_name: str,
        usecase: ConvertUseCase,
    ) -> None:
        self._url = amqp_url
        self._queue_name = queue_name
        self._usecase = usecase
        self._connection: aio_pika.RobustConnection | None = None
        self._channel: aio_pika.abc.AbstractChannel | None = None

    async def start(self) -> None:
        self._connection = await aio_pika.connect_robust(self._url)
        self._channel = await self._connection.channel()
        await self._channel.set_qos(prefetch_count=10)
        queue = await self._channel.declare_queue(self._queue_name, durable=True)
        await queue.consume(self._handle)
        logger.info("rpc_server_started", extra={"queue": self._queue_name})

    async def stop(self) -> None:
        if self._connection is not None:
            await self._connection.close()
            self._connection = None
            self._channel = None

    async def _handle(self, message: aio_pika.abc.AbstractIncomingMessage) -> None:
        async with message.process(requeue=False):
            response_body = self._dispatch(message.body)
            if message.reply_to is None:
                logger.warning("rpc_request_without_reply_to")
                return
            assert self._channel is not None
            await self._channel.default_exchange.publish(
                aio_pika.Message(
                    body=response_body,
                    content_type="application/json",
                    correlation_id=message.correlation_id,
                ),
                routing_key=message.reply_to,
            )

    def _dispatch(self, raw: bytes) -> bytes:
        try:
            payload = json.loads(raw)
            source = Money(amount=Decimal(str(payload["amount"])), currency=payload["from"])
            target_currency = str(payload["to"])
            result = self._usecase.execute(source, target_currency)
            return json.dumps(
                {
                    "amount": str(result.target.amount),
                    "currency": result.target.currency,
                    "source_amount": str(result.source.amount),
                    "source_currency": result.source.currency,
                    "rate_date": result.rate_date.isoformat(),
                }
            ).encode()
        except UnknownCurrency as e:
            return json.dumps({"error": "unknown_currency", "currency": e.currency}).encode()
        except RatesUnavailable:
            return json.dumps({"error": "rates_unavailable"}).encode()
        except (KeyError, ValueError, TypeError) as e:
            logger.warning("rpc_invalid_request", extra={"error": str(e)})
            return json.dumps({"error": "invalid_request", "message": str(e)}).encode()
