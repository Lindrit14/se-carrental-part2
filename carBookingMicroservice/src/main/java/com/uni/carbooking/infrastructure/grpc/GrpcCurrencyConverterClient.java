package com.uni.carbooking.infrastructure.grpc;

import com.uni.carbooking.application.port.out.CurrencyConverterPort;
import com.uni.carbooking.domain.error.CurrencyConversionFailed;
import com.uni.carbooking.domain.money.Money;
import com.uni.currency.grpc.v1.ConvertRequest;
import com.uni.currency.grpc.v1.ConvertResponse;
import com.uni.currency.grpc.v1.CurrencyConverterGrpc;
import io.grpc.Status;
import io.grpc.StatusRuntimeException;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import java.math.BigDecimal;
import java.util.concurrent.TimeUnit;

/**
 * gRPC client for the currency-converter service. Replaces the previous
 * RabbitMQ-RPC adapter; the port boundary ({@link CurrencyConverterPort}) is
 * unchanged so {@code CreateBookingUseCase} doesn't move.
 *
 * <p>Any non-{@code OK} status — INVALID_ARGUMENT, FAILED_PRECONDITION,
 * UNAVAILABLE, DEADLINE_EXCEEDED, etc. — is surfaced as
 * {@link CurrencyConversionFailed}, the same exception the old adapter raised,
 * so upstream error mapping is unaffected.
 */
@Component
class GrpcCurrencyConverterClient implements CurrencyConverterPort {

    private final CurrencyConverterGrpc.CurrencyConverterBlockingStub stub;
    private final long deadlineMs;

    GrpcCurrencyConverterClient(
            CurrencyConverterGrpc.CurrencyConverterBlockingStub stub,
            @Value("${grpc.client.currency-converter.deadline-ms:5000}") long deadlineMs) {
        this.stub = stub;
        this.deadlineMs = deadlineMs;
    }

    @Override
    public Money convert(Money source, String targetCurrency) {
        var request = ConvertRequest.newBuilder()
                .setAmount(source.amount().toPlainString())
                .setFromCurrency(source.currency())
                .setToCurrency(targetCurrency)
                .build();

        ConvertResponse reply;
        try {
            reply = stub
                    .withDeadlineAfter(deadlineMs, TimeUnit.MILLISECONDS)
                    .convert(request);
        } catch (StatusRuntimeException e) {
            throw new CurrencyConversionFailed(messageFor(e.getStatus()));
        }

        return new Money(new BigDecimal(reply.getAmount()), reply.getCurrency());
    }

    private static String messageFor(Status status) {
        var code = status.getCode();
        var description = status.getDescription();
        return switch (code) {
            case DEADLINE_EXCEEDED -> "currency-converter timed out";
            case UNAVAILABLE -> "currency-converter unavailable: "
                    + (description == null ? code.name() : description);
            default -> "currency-converter error: " + code.name()
                    + (description == null ? "" : ": " + description);
        };
    }
}
